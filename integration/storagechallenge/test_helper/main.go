package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pastelnetwork/gonode/pastel"
)

var blockCount int64 = 0
var ticker = time.NewTicker(time.Second * 3)
var stopChan chan struct{}

func startGenerateblock() {
	for {
		select {
		case <-ticker.C:
			atomic.AddInt64(&blockCount, 1)
		case <-stopChan:
			ticker.Stop()
		}
	}
}

func stopGenerateBlock() {
	stopChan <- struct{}{}
}

func getBlockCount(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte(fmt.Sprint(blockCount)))
}

func getMasternodeList(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(mnList)
}

func storeImageByID(k int, storageName, storageIPAddr string, keyList *[]string) {
	b, err := ioutil.ReadFile(fmt.Sprintf("/root/pastel/%d.jpg", k))
	if err != nil {
		log.Fatal(k, "read image", err)
	}
	b, err = json.Marshal(map[string]interface{}{"value": b})
	if err != nil {
		log.Fatal(k, "marshal", err)
	}
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://%s/p2p", storageIPAddr), bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	req.Close = true
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Fatal(k, "post p2p store", err)
	}
	defer resp.Body.Close()

	var mapKey map[string]string
	err = json.NewDecoder(resp.Body).Decode(&mapKey)
	if err != nil {
		log.Fatal(k, "unmarshal", err)
	}

	*keyList = append(*keyList, mapKey["key"])

	fmt.Println("Finish pushing file", k, "to p2p storage")
}

func incrementalKeyStoring(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("OK"))
	var keyList = make([]string, 0)
	b, _ := ioutil.ReadFile("p2pkeys.json")
	json.Unmarshal(b, &keyList)
	storageName := "mn1key"
	storageIPAddr := "192.168.100.14:9090"
	for i := 1; i < 5; i++ {
		storeImageByID(i, storageName, storageIPAddr, &keyList)
	}

	b, _ = json.Marshal(keyList)

	if err := ioutil.WriteFile("p2pkeys.json", b, 0644); err != nil {
		log.Fatal("write json file", err)
	}

	w.Write(b)
}

func resetChallengeStatus(w http.ResponseWriter, _ *http.Request) {
	mapNodeStatictis = make(map[string]*statictis)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func listKeys(w http.ResponseWriter, _ *http.Request) {
	b, err := ioutil.ReadFile("p2pkeys.json")
	if err != nil {
		b = []byte("[]")
	}

	w.Write(b)
}

func main() {
	runtime.GOMAXPROCS(5)
	go startGenerateblock()
	defer stopGenerateBlock()

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("OK")) })
	http.HandleFunc("/getblockcount", getBlockCount)
	http.HandleFunc("/mnlist", getMasternodeList)

	// store more incremental files
	http.HandleFunc("/store/incrementals", incrementalKeyStoring)
	http.HandleFunc("/store/keys", listKeys)

	http.HandleFunc("/sts/sent", challengeSent)
	http.HandleFunc("/sts/respond", challengeResponded)
	http.HandleFunc("/sts/succeeded", challengeVerified)
	http.HandleFunc("/sts/failed", challengeFailed)
	http.HandleFunc("/sts/timeout", challengeTimeout)
	http.HandleFunc("/sts/show", statictisShow)
	http.HandleFunc("/sts/reset", resetChallengeStatus)

	go verifyTimeout()
	defer close(stopCh)

	fmt.Println("server starting")
	err := http.ListenAndServe("0.0.0.0:8088", nil)
	if err != nil {
		log.Fatal("Listen and serve", err)
	}

	fmt.Println("server shutdown")
}

var mtx sync.Mutex
var sentMap = make(map[string]map[string]map[string]string)
var mapNodeStatictis = make(map[string]*statictis)
var mapFileStatictis = make(map[string]*statictis)
var stopCh = make(chan struct{})

type statictis struct {
	Sent      int32 `json:"sent"`
	Responded int32 `json:"respond"`
	Succeeded int32 `json:"success"`
	Failed    int32 `json:"failed"`
	Timeout   int32 `json:"timeout"`
}

func challengeSent(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.PostForm.Get("id")
	key := r.PostForm.Get("key")
	nodeID := r.PostForm.Get("node_id")
	sentBlock := r.PostForm.Get("sent_block")
	mtx.Lock()
	if sentMap[nodeID] == nil {
		sentMap[nodeID] = make(map[string]map[string]string)
	}
	sentMap[nodeID][id] = map[string]string{"key": key, "sent_block": sentBlock}
	mtx.Unlock()

	if _, ok := mapNodeStatictis[nodeID]; !ok {
		mapNodeStatictis[nodeID] = &statictis{}
	}
	if st, ok := mapNodeStatictis[nodeID]; ok {
		atomic.AddInt32(&st.Sent, 1)
	}
	log.Println("handled challenge sent statictis", mapNodeStatictis[nodeID])
	w.WriteHeader(http.StatusOK)
}

func challengeVerified(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.PostForm.Get("id")
	nodeID := r.PostForm.Get("node_id")
	mtx.Lock()
	delete(sentMap[nodeID], id)
	mtx.Unlock()

	if st, ok := mapNodeStatictis[nodeID]; ok {
		atomic.AddInt32(&st.Succeeded, 1)
	}
	log.Println("handled challenge succeeded statictis", mapNodeStatictis[nodeID])
	w.WriteHeader(http.StatusOK)
}

func challengeFailed(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.PostForm.Get("id")
	nodeID := r.PostForm.Get("node_id")
	mtx.Lock()
	delete(sentMap[nodeID], id)
	mtx.Unlock()

	if st, ok := mapNodeStatictis[nodeID]; ok {
		atomic.AddInt32(&st.Failed, 1)
	}
	log.Println("handled challenge failed statictis", mapNodeStatictis[nodeID])
	w.WriteHeader(http.StatusOK)
}

func challengeResponded(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	nodeID := r.PostForm.Get("node_id")

	if st, ok := mapNodeStatictis[nodeID]; ok {
		atomic.AddInt32(&st.Responded, 1)
	}
	log.Println("handled challenge respond statictis", mapNodeStatictis[nodeID])
	w.WriteHeader(http.StatusOK)
}

func challengeTimeout(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.PostForm.Get("id")
	nodeID := r.PostForm.Get("node_id")
	mtx.Lock()
	delete(sentMap[nodeID], id)
	mtx.Unlock()

	if st, ok := mapNodeStatictis[nodeID]; ok {
		atomic.AddInt32(&st.Timeout, 1)
	}
	log.Println("handled challenge timeout statictis", mapNodeStatictis[nodeID])
	w.WriteHeader(http.StatusOK)
}

func verifyTimeout() {
	tc := time.NewTicker(time.Second * 6)
	defer tc.Stop()
	for {
		select {
		case <-stopCh:
			return
		case <-tc.C:
			log.Println("current checking block count:", blockCount)
			for nodeID, mapID := range sentMap {
				for id, mapKey := range mapID {
					sentBlkStr := mapKey["sent_block"]
					if sentBlkStr == "" {
						mtx.Lock()
						delete(sentMap[nodeID], id)
						mtx.Unlock()
					}

					sentBlk, err := strconv.Atoi(sentBlkStr)
					if err != nil {
						continue
					}
					if blockCount > int64(sentBlk)+1 {
						mtx.Lock()
						delete(sentMap[nodeID], id)
						mtx.Unlock()
						if st, ok := mapNodeStatictis[nodeID]; ok {
							atomic.AddInt32(&st.Timeout, 1)
						}
						log.Println("handled challenge timeout statictis", mapNodeStatictis[nodeID])
					}
				}
			}
		}
	}
}

func statictisShow(w http.ResponseWriter, _ *http.Request) {
	b, _ := json.Marshal(mapNodeStatictis)
	w.Write(b)
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var mapWsConn = make(map[string]*websocket.Conn)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}

var mnList = pastel.MasterNodes{}

type message struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

func deregistrationWs(id string) {
	var idx = 0
	var mn pastel.MasterNode
	for idx, mn = range mnList {
		if mn.ExtKey == id {
			break
		}
	}
	if len(mnList) > idx {
		mnList = append(mnList[:idx], mnList[idx+1:]...)
	}
	removeConn := mapWsConn[id]
	defer removeConn.Close()
	delete(mapWsConn, id)
	for _, conn := range mapWsConn {
		conn.WriteJSON(&message{Type: "deregistration", Data: []byte(id)})
	}
}

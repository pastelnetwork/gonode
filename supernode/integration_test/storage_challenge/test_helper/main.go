package main

import (
	"bytes"
	"context"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pastelnetwork/gonode/pastel"
	"gopkg.in/yaml.v2"
)

var mockFileStored, imageDownloaded bool
var blockCount int64 = 0
var ticker = time.NewTicker(time.Second * 10)
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
func prepareMockFiles() {
	urlMask := "https://picsum.photos/id/%d/60/90"
	wg := sync.WaitGroup{}
	for i := 1; i < 24; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()
			req, _ := http.NewRequest("GET", fmt.Sprintf(urlMask, 200+i), nil)
			req.Close = true
			resp, err := (&http.Client{}).Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			image, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			err = ioutil.WriteFile(fmt.Sprintf("/root/pastel/%d.jpg", i), image, 0644)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("succeeded download test image %d.png\n", i)
		}(i, &wg)
		if math.Mod(float64(i), 5) == 0 {
			wg.Wait()
		}
	}
	imageDownloaded = true
}

func storeMockFiles() {
	if !imageDownloaded {
		prepareMockFiles()
	}

	var k int
	var storageIPAddr, storageName string
	var mapKeys = make(map[string]map[string]string)
	var keyList = make([]string, 0)
	for j := 1; j <= len(mapWsConn); j++ {
		storageIPAddr = fmt.Sprintf("192.168.100.1%d", j)
		storageName = fmt.Sprintf("mn%dkey", j)
		mapKeys[storageName] = make(map[string]string)
		for i := 1; i <= 4; i++ {
			k = i + j*4

			if k > 24 {
				k -= 24
			}
			storeImageByID(k, storageName, storageIPAddr, mapKeys, &keyList)
		}
	}

	b, err := yaml.Marshal(mapKeys)
	if err != nil {
		log.Fatal("yaml marshal", err)
	}

	if err = ioutil.WriteFile("p2pkeys.yml", b, 0644); err != nil {
		log.Fatal("write yaml file", err)
	}

	b, err = json.Marshal(keyList)
	if err != nil {
		log.Fatal("yaml marshal key list", err)
	}

	if err = ioutil.WriteFile("p2pkeys.json", b, 0644); err != nil {
		log.Fatal("write json file", err)
	}
	mockFileStored = true
}

func storeImageByID(k int, storageName, storageIPAddr string, mapKeys map[string]map[string]string, keyList *[]string) {
	b, err := ioutil.ReadFile(fmt.Sprintf("/root/pastel/%d.jpg", k))
	if err != nil {
		log.Fatal(k, "read image", err)
	}
	b, err = json.Marshal(map[string]interface{}{"value": b})
	if err != nil {
		log.Fatal(k, "marshal", err)
	}
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://%s:9090/p2p", storageIPAddr), bytes.NewReader(b))
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

	mapKeys[storageName][fmt.Sprintf("%d.png", k)] = mapKey["key"]
	*keyList = append(*keyList, mapKey["key"])

	fmt.Println("Finish pushing file", k, "to p2p storage")
}

func incrementalKeyStoring(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("OK"))
	b, _ := ioutil.ReadFile("p2pkeys.yml")
	var mapKeys = make(map[string]map[string]string)
	yaml.Unmarshal(b, &mapKeys)
	var keyList = make([]string, 0)
	b, _ = ioutil.ReadFile("p2pkeys.json")
	json.Unmarshal(b, &keyList)
	storageName := "mn1key"
	storageIPAddr := "192.168.100.11"
	for i := 20; i < 25; i++ {
		storeImageByID(i, storageName, storageIPAddr, mapKeys, &keyList)
	}

	b, err := yaml.Marshal(mapKeys)
	if err != nil {
		log.Fatal("yaml marshal", err)
	}

	if err = ioutil.WriteFile("p2pkeys.yml", b, 0644); err != nil {
		log.Fatal("write yaml file", err)
	}

	b, err = json.Marshal(keyList)
	if err != nil {
		log.Fatal("yaml marshal key list", err)
	}

	if err = ioutil.WriteFile("p2pkeys.json", b, 0644); err != nil {
		log.Fatal("write json file", err)
	}

	w.Write(b)
}

func resetChallengeStatus(w http.ResponseWriter, _ *http.Request) {
	mapStatictis = make(map[string]*statictis)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func keyStoring(w http.ResponseWriter, _ *http.Request) {
	if !mockFileStored {
		storeMockFiles()
	}
	w.Write([]byte("OK"))
}

func listKeys(w http.ResponseWriter, _ *http.Request) {
	b, err := ioutil.ReadFile("p2pkeys.json")
	if err != nil {
		b = []byte("[]")
	}

	w.Write(b)
	w.WriteHeader(http.StatusOK)
}

func main() {
	runtime.GOMAXPROCS(5)
	go startGenerateblock()
	defer stopGenerateBlock()

	http.HandleFunc("/ws", wsHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("OK")) })
	http.HandleFunc("/getblockcount", getBlockCount)
	http.HandleFunc("/mnlist", getMasternodeList)

	http.HandleFunc("/store/mocks", keyStoring)
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
var mapStatictis = make(map[string]*statictis)
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

	if _, ok := mapStatictis[nodeID]; !ok {
		mapStatictis[nodeID] = &statictis{}
	}
	if st, ok := mapStatictis[nodeID]; ok {
		atomic.AddInt32(&st.Sent, 1)
	}
	log.Println("handled challenge sent statictis", mapStatictis[nodeID])
	w.WriteHeader(http.StatusOK)
}

func challengeVerified(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.PostForm.Get("id")
	nodeID := r.PostForm.Get("node_id")
	mtx.Lock()
	delete(sentMap[nodeID], id)
	mtx.Unlock()

	if st, ok := mapStatictis[nodeID]; ok {
		atomic.AddInt32(&st.Succeeded, 1)
	}
	log.Println("handled challenge succeeded statictis", mapStatictis[nodeID])
	w.WriteHeader(http.StatusOK)
}

func challengeFailed(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.PostForm.Get("id")
	nodeID := r.PostForm.Get("node_id")
	mtx.Lock()
	delete(sentMap[nodeID], id)
	mtx.Unlock()

	if st, ok := mapStatictis[nodeID]; ok {
		atomic.AddInt32(&st.Failed, 1)
	}
	log.Println("handled challenge failed statictis", mapStatictis[nodeID])
	w.WriteHeader(http.StatusOK)
}

func challengeResponded(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	nodeID := r.PostForm.Get("node_id")

	if st, ok := mapStatictis[nodeID]; ok {
		atomic.AddInt32(&st.Responded, 1)
	}
	log.Println("handled challenge respond statictis", mapStatictis[nodeID])
	w.WriteHeader(http.StatusOK)
}

func challengeTimeout(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.PostForm.Get("id")
	nodeID := r.PostForm.Get("node_id")
	mtx.Lock()
	delete(sentMap[nodeID], id)
	mtx.Unlock()

	if st, ok := mapStatictis[nodeID]; ok {
		atomic.AddInt32(&st.Timeout, 1)
	}
	log.Println("handled challenge timeout statictis", mapStatictis[nodeID])
	w.WriteHeader(http.StatusOK)
}

func verifyTimeout() {
	tc := time.NewTicker(time.Second * 10)
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
						if st, ok := mapStatictis[nodeID]; ok {
							atomic.AddInt32(&st.Timeout, 1)
						}
						log.Println("handled challenge timeout statictis", mapStatictis[nodeID])
					}
				}
			}
		}
	}
}

func statictisShow(w http.ResponseWriter, _ *http.Request) {
	b, _ := json.Marshal(mapStatictis)
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

var mnList = pastel.MasterNodes{
	{
		Rank:       "1",
		IPAddress:  "192.168.100.11:18232",
		ExtAddress: "192.168.100.11:14444",
		ExtKey:     base32.StdEncoding.EncodeToString([]byte("mn1key")),
	},
	{
		Rank:       "2",
		IPAddress:  "192.168.100.12:18232",
		ExtAddress: "192.168.100.12:14444",
		ExtKey:     base32.StdEncoding.EncodeToString([]byte("mn2key")),
	},
	{
		Rank:       "3",
		IPAddress:  "192.168.100.13:18232",
		ExtAddress: "192.168.100.13:14444",
		ExtKey:     base32.StdEncoding.EncodeToString([]byte("mn3key")),
	},
	{
		Rank:       "4",
		IPAddress:  "192.168.100.14:18232",
		ExtAddress: "192.168.100.14:14444",
		ExtKey:     base32.StdEncoding.EncodeToString([]byte("mn4key")),
	},
}

type message struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	defer conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	ctx, cncl := context.WithCancel(r.Context())
	defer cncl()
	go func(ctx context.Context, conn *websocket.Conn) {
		ticker := time.NewTicker(pingPeriod)
		ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}(ctx, conn)
	id := r.Header.Get("id")
	if id == "" {
		log.Println("required node id in header, but got empty")
		return
	}
	if id == base32.StdEncoding.EncodeToString([]byte("mn5key")) {
		mnList = append(mnList, pastel.MasterNode{
			Rank:       "5",
			IPAddress:  "192.168.100.15:18232",
			ExtAddress: "192.168.100.15:14444",
			ExtKey:     id,
		})
	} else if id == base32.StdEncoding.EncodeToString([]byte("mn6key")) {
		mnList = append(mnList, pastel.MasterNode{
			Rank:       "5",
			IPAddress:  "192.168.100.16:18232",
			ExtAddress: "192.168.100.16:14444",
			ExtKey:     id,
		})
	}
	_, msg, err := conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("read registration IsUnexpectedCloseError: %v", err)
		}
		return
	}
	if string(msg) != "registration" {
		log.Printf("node registration message expected, got %s", string(msg))
		return
	}

	mapWsConn[id] = conn
	defer deregistrationWs(id)
	var nodeList = make([]string, 0)
	for nodeID := range mapWsConn {
		nodeList = append(nodeList, nodeID)
	}
	b, _ := json.Marshal(nodeList)
	err = conn.WriteJSON(&message{Type: "registration", Data: b})
	if err != nil {
		log.Printf("could not reply registration message: %v", err)
		return
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read message IsUnexpectedCloseError: %v", err)
			}
			break
		}
		var messageData *message
		if err = json.Unmarshal(msg, &messageData); err != nil {
			log.Printf("read message IsUnexpectedCloseError: %v", err)
			break
		}
		switch messageData.Type {
		case "store":
			for nodeID, nodeConn := range mapWsConn {
				err = nodeConn.WriteJSON(messageData)
				if err != nil {
					if nodeConn != nil {
						nodeConn.Close()
					}
					deregistrationWs(nodeID)
				}
			}
		}
	}
}

func deregistrationWs(id string) {
	var idx = 0
	var mn pastel.MasterNode
	for idx, mn = range mnList {
		if mn.ExtKey == id {
			break
		}
	}
	mnList = append(mnList[:idx], mnList[idx+1:]...)
	delete(mapWsConn, id)
	for _, conn := range mapWsConn {
		conn.WriteJSON(&message{Type: "deregistration", Data: []byte(id)})
	}
}

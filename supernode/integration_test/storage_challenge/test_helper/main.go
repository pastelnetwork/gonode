package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
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
	w.WriteHeader(http.StatusOK)
}

func prepareMockFiles() {
	urlMask := "https://picsum.photos/%d/60/90"
	for i := 0; i < 500; i++ {
		resp, err := http.Get(fmt.Sprint(urlMask, i))
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
		fmt.Printf("succeeded download test image %d.png", i)
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
	for j := 1; j < 6; j++ {
		storageIPAddr = fmt.Sprintf("192.168.100.1%d", j)
		storageName = fmt.Sprintf("mn%dkey", j)
		mapKeys[storageName] = make(map[string]string)
		for i := 1; i <= 20; i++ {
			k = i + j*20
			if k > 100 {
				k -= 100
			}
			b, err := ioutil.ReadFile(fmt.Sprintf("/root/pastel/%d.jpg", k))
			if err != nil {
				log.Fatal(j, i, "read image", err)
			}
			b, err = json.Marshal(map[string]interface{}{"value": b})
			if err != nil {
				log.Fatal(j, i, "marshal", err)
			}
			resp, err := (&http.Client{}).Post(fmt.Sprintf("http://%s:9090/p2p", storageIPAddr), "application/json", bytes.NewReader(b))
			if err != nil {
				log.Fatal(j, i, "post p2p store", err)
			}
			defer resp.Body.Close()

			b, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(j, i, "read response", err)
			}

			var mapKey map[string]string
			err = json.Unmarshal(b, &mapKey)
			if err != nil {
				log.Fatal(j, i, "unmarshal", err)
			}

			mapKeys[storageName][fmt.Sprintf("%d.png", k)] = mapKey["key"]
			keyList = append(keyList, mapKey["key"])

			fmt.Println("Finish pushing file", k, "to p2p storage", j)
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

func keyStoring(w http.ResponseWriter, _ *http.Request) {
	if !mockFileStored {
		storeMockFiles()
	}
	w.WriteHeader(http.StatusOK)
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
	go startGenerateblock()
	defer stopGenerateBlock()

	http.HandleFunc("/ws", wsHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {})
	http.HandleFunc("/getblockcount", getBlockCount)

	http.HandleFunc("/store/mocks", keyStoring)
	http.HandleFunc("/store/keys", listKeys)

	http.HandleFunc("/sts/sent", challengeSent)
	http.HandleFunc("/sts/respond", challengeResponded)
	http.HandleFunc("/sts/succeeded", challengeVerified)
	http.HandleFunc("/sts/failed", challengeFailed)
	http.HandleFunc("/sts/timeout", challengeTimeout)
	http.HandleFunc("/sts/show", statictisShow)

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
			resp, err := (&http.Client{}).Get(fmt.Sprintf("http://%s/getblockcount", "192.168.100.10"))
			if err != nil {
				log.Println("could not current getblockcount", err)
				continue
			}
			defer resp.Body.Close()
			b, _ := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("could not read current getblockcount response", err)
				continue
			}
			b = bytes.Trim(b, "\n")
			blkCount, err := strconv.Atoi(string(b))
			if err != nil {
				log.Println("could not convert current getblockcount to integer", err)
				continue
			}
			log.Println("current checking block count:", blkCount)
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
					if blkCount > sentBlk+1 {
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
	delete(mapWsConn, id)
	for _, conn := range mapWsConn {
		conn.WriteJSON(&message{Type: "deregistration", Data: []byte(id)})
	}
}

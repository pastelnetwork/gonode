package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/btcsuite/btcutil/base58"
	"github.com/gorilla/websocket"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/p2p"
	_ "github.com/rqlite/go-sqlite3" //go-sqlite3
	"golang.org/x/crypto/sha3"
)

const (
	// Maximum message size allowed from peer.
	maxMessageSize = 1024
	dbName         = "data001.sqlite3"
)

type mockP2P struct {
	Nodes []string
	ID    string
	conn  *websocket.Conn
	db    *sql.DB
}

type message struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

func (s *mockP2P) registrationNode(ctx context.Context) error {
	header := http.Header{}
	header.Add("id", s.ID)
	log.P2P().WithField("id", s.ID).Info("Dial ws connection")
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, "ws://helper:8088/ws", header)
	if err != nil {
		log.P2P().WithContext(ctx).WithError(err).Error("could not dial websocket join peers")
		conn.Close()
		return err
	}
	if err = conn.WriteMessage(websocket.TextMessage, []byte("registration")); err != nil {
		log.P2P().WithContext(ctx).WithError(err).Error("could not send node registration to peers")
		conn.Close()
		return err
	}
	if msgType, msg, err := conn.ReadMessage(); err != nil {
		log.P2P().WithContext(ctx).WithError(err).Error("could not send receive node registration accepted from peers")
		conn.Close()
		return err
	} else {
		if msgType == websocket.TextMessage {
			var msgData *message
			if err = json.Unmarshal(msg, &msgData); err != nil {
				log.P2P().WithContext(ctx).WithError(err).Error("could not parse node registration accepted from peers")
				conn.Close()
				return err
			}
			if msgData.Type != "registration" {
				log.P2P().WithContext(ctx).Error("node registration accepted data is not correct")
				conn.Close()
				return err
			}

			if err = json.Unmarshal(msgData.Data, &s.Nodes); err != nil {
				log.P2P().WithContext(ctx).WithError(err).Error("could not parse registration node list")
				conn.Close()
				return err
			}
		} else {
			log.P2P().WithContext(ctx).Error("websocket message type is not correct")
			conn.Close()
			return err
		}
	}
	s.conn = conn
	return nil
}

func (s *mockP2P) deregistration(_ context.Context) {
	s.conn.Close()
	s.conn = nil
}

func (s *mockP2P) Run(ctx context.Context) error {
	if err := s.newStore(ctx); err != nil {
		return err
	}

	if err := s.registrationNode(ctx); err != nil {
		return err
	}
	defer s.deregistration(ctx)
	var msgData *message

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, msg, err := s.conn.ReadMessage()
			if err != nil {
				log.P2P().WithContext(ctx).WithError(err).Error("could not read ws message")
				return err
			}
			log.P2P().WithContext(ctx).Infof("recv: %s", string(msg))
			if err = json.Unmarshal(msg, &msgData); err != nil {
				log.P2P().WithContext(ctx).WithError(err).Error("could not parse received message from peers")
				return err
			}

			switch msgData.Type {
			case "deretistration":
				id := string(msgData.Data)
				for idx := 0; idx < len(s.Nodes); idx++ {
					if s.Nodes[idx] == id {
						s.Nodes = append(s.Nodes[:idx], s.Nodes[idx+1:]...)
						break
					}
				}
				log.P2P().WithContext(ctx).Infof("node deregistration: %s", id)
			case "store":
				key, err := s.storeLocal(ctx, msgData.Data)
				if err != nil {
					log.P2P().WithContext(ctx).WithError(err).Errorf("could store record recived from peers: %s", key)
					return err
				}
			}
		}
	}
}

func (s *mockP2P) newStore(ctx context.Context) error {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.P2P().WithContext(ctx).WithError(err).Error("cannot open sqlite database")
		return err
	}
	s.db = db
	if !s.checkStore(ctx) {
		if err = s.migrate(ctx); err != nil {
			log.P2P().WithContext(ctx).WithError(err).Error("cannot create table(s) in sqlite database")
			return err
		}
	}

	return nil
}

func (s *mockP2P) checkStore(ctx context.Context) bool {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='data'`
	var name string
	return s.db.QueryRowContext(ctx, query).Scan(&name) == nil
}

func (s *mockP2P) migrate(ctx context.Context) error {
	query := `
    CREATE TABLE IF NOT EXISTS data(
        key TEXT PRIMARY KEY,
        value BLOB NOT NULL
    );
    `
	if _, err := s.db.ExecContext(ctx, query); err != nil {
		log.P2P().WithContext(ctx).WithError(err).Error("failed to create table 'data'")
		return err
	}
	return nil
}

func (s *mockP2P) Store(ctx context.Context, value []byte) (string, error) {
	key, err := s.storeLocal(ctx, value)
	if err != nil {
		return key, err
	}

	if err = s.conn.WriteJSON(&message{Type: "store", Data: value}); err != nil {
		log.P2P().WithContext(ctx).WithError(err).Errorf("replicas record with key %s", key)
		return key, err
	}

	return key, nil
}

func (s *mockP2P) storeLocal(ctx context.Context, value []byte) (string, error) {
	sha := sha3.Sum256(value)
	keyBytes := sha[:]
	key := base58.Encode(keyBytes)
	query := `
	INSERT INTO data(key, value) values(?, ?) ON CONFLICT(key) DO UPDATE SET value=?
    `
	res, err := s.db.ExecContext(ctx, query, key, value, value)
	if err != nil {
		log.P2P().WithContext(ctx).WithError(err).Errorf("cannot insert or update record with key %s", key)
		return "", err
	}

	if rowsAffected, err := res.RowsAffected(); err != nil || rowsAffected == 0 {
		log.P2P().WithContext(ctx).WithError(err).Errorf("failed to insert/update record with key %s", key)
		return "", fmt.Errorf("failed to insert/update record with key %s: %w", key, err)
	}

	return key, nil
}

func (s *mockP2P) Retrieve(ctx context.Context, key string, _ ...bool) ([]byte, error) {
	var value = []byte{}
	err := s.db.QueryRowContext(ctx, `SELECT value FROM data WHERE key = ?`, key).Scan(&value)
	if err != nil {
		return nil, fmt.Errorf("failed to get record by key %s: %w", key, err)
	}

	return value, nil
}

func (s *mockP2P) Delete(ctx context.Context, key string) error {
	res, err := s.db.Exec("DELETE FROM data WHERE key = ?", key)
	if err != nil {
		log.P2P().WithContext(ctx).Errorf("cannot delete record by key %s: %v", key, err)
		return err
	}

	if rowsAffected, err := res.RowsAffected(); err != nil || rowsAffected == 0 {
		log.P2P().WithContext(ctx).WithError(err).Errorf("failed to delete record by key %s: %v", key, err)
		return err
	}
	return nil
}

// Stats returns stats of store
func (s *mockP2P) Stats(_ context.Context) (map[string]interface{}, error) {
	stats := map[string]interface{}{}
	return stats, nil
}

func (s *mockP2P) NClosestNodes(ctx context.Context, n int, key string, ignores ...string) []string {
	return utils.GetNClosestXORDistanceStringToAGivenComparisonString(n, key, s.Nodes, ignores...)
}

func newMockP2P(id string) p2p.P2P {
	return &mockP2P{ID: id}
}

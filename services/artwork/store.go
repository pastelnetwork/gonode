package artwork

import (
	"sync"

	"github.com/pastelnetwork/go-commons/log"
	"github.com/pastelnetwork/go-commons/random"
)

type Store struct {
	sync.Mutex
	images map[string]*Image
}

func (store *Store) Get(id string) *Image {
	store.Lock()
	defer store.Unlock()

	if image, ok := store.images[id]; ok {
		return image
	}
	return nil
}

func (store *Store) Remove(id string) {
	store.Lock()
	defer store.Unlock()

	if _, ok := store.images[id]; ok {
		log.WithField("id", id).Debugln(logPrefix, "Remove image")
		delete(store.images, id)
	}
}

func (store *Store) Add(image *Image) string {
	store.Lock()
	defer store.Unlock()

	id, _ := random.String(8, random.Base62Chars)

	log.WithField("id", id).WithField("bytes", image.Size).Debugln(logPrefix, "Add image")
	store.images[id] = image

	return id
}

// func (store *Store) Add(image *Image) (string, time.Time) {
// 	store.Lock()
// 	defer store.Unlock()

// 	id, _ := random.String(8, random.Base62Chars)
// 	expiresIn := time.Now().Add(imageTTL)

// 	log.Debugln(logPrefix, "Add image", id)
// 	store.images[id] = image

// 	go func() {
// 		time.AfterFunc(time.Until(expiresIn), func() {
// 			store.Remove(id)
// 		})
// 	}()

// 	return id, expiresIn
// }

func (store *Store) LoadImage(image *Image) {

}

func NewStore() *Store {
	return &Store{
		images: make(map[string]*Image),
	}
}

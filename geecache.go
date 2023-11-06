package geecache

import (
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		return nil
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}
func (g *Group) Get(key string) (ByteView, error) {
	// key is not provided
	if key == "" {
		return ByteView{}, fmt.Errorf("keys is required")
	}
	//key-value is stored in mainCache and hit
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[geeCache] hit")
		return v, nil
	}
	//k-v not stores in mainCache and loads from disk
	return g.load(key)
}

//func (g *Group) load(key string) (ByteView, error) {
//	return g.getLocally(key)
//}

func (g *Group) getLocally(key string) (ByteView, error) {
	//get date from source
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: byteClones(bytes)}
	//take the data into the cache
	g.populateCache(key, value)
	return value, nil
}
func (g *Group) populateCache(key string, value ByteView) {
	//mu.Lock()
	//defer mu.Unlock()
	g.mainCache.Add(key, value)
}

func (g *Group) RegisterPeer(peers PeerPicker) {
	if g.peers != nil {
		panic("registerPeerPickPeer called more than once")
	}
	g.peers = peers
}

func (g *Group) load(key string) (value ByteView, err error) {
	if g.peers != nil {
		if peer, ok := g.peers.PickPeer(key); ok {
			if value, err := g.getFromPeer(peer, key); err == nil {
				return value, nil
			}
			log.Println("[GeeCache] Failed to get from peer", err)
		}
	}
	return g.getLocally(key)
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}

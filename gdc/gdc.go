package gdc

import (
	pb "github.com/GallifreyGoTutoural/ggt-dist-cache/gdccachepb"
	"github.com/GallifreyGoTutoural/ggt-dist-cache/singleflight"
	"sync"
)

// Getter is the interface that must be implemented to fetch data.
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc implements Getter with a function.
type GetterFunc func(key string) ([]byte, error)

// Get implements Getter interface.
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group is a cache namespace and associated data loaded spread over
type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker
	loader    *singleflight.Group
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup creates a new instance of Group.
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &singleflight.Group{},
	}
	groups[name] = g
	return g
}

// GetGroup returns the named group previously created with NewGroup, or nil if there's no such group.
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}

// Get value for a key from cache
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, nil
	}
	//if mainCache has the key, return it
	if v, ok := g.mainCache.get(key); ok {
		return v, nil
	}
	//else load it
	return g.load(key)
}

// load value for a key from cache
func (g *Group) load(key string) (ByteView, error) {
	//use singleflight to ensure that each key is only fetched once
	view, err := g.loader.Do(key, func() (interface{}, error) {
		//if peers is not nil, use PeerPicker to get value from remote peer
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err := g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
			}
		}

		//else use local getter to get value
		return g.getLocally(key)
	})
	if err == nil {
		return view.(ByteView), nil
	}
	return ByteView{}, err

}

// getLocally value for a key from cache
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

// populateCache value for a key from cache
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

// RegisterPeerPicker registers a PeerPicker for choosing remote peer
func (g *Group) RegisterPeerPicker(picker PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = picker
}

// PickPeer picks a peer according to key
func (g *Group) PickPeer(key string) (PeerGetter, bool) {
	if g.peers == nil {
		return nil, false
	}
	return g.peers.PickPeer(key)
}

// getFromPeer gets value from remote peer
func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &pb.Request{
		Group: g.name,
		Key:   key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
}

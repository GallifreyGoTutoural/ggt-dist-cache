package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

/*
Do executes and returns the results of the given function, making sure that
only one execution is in-flight for a given key at a time. If a duplicate
comes in, the duplicate caller waits for the original to complete and
receives the same results. The return value shared indicates whether v was
given to multiple callers.
*/
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	//lazy initialization
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	//if the call is in-flight, wait for it and return its results
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	//if not exist, create a new call
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()
	//execute the function
	c.val, c.err = fn()
	c.wg.Done()
	//delete the call
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
	//return the result
	return c.val, c.err
}

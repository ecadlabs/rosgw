package pool

import (
	"container/list"
	"sync"
)

type Finalizer interface {
	Finalize()
}

type Item interface {
	Key() interface{}
}

type poolItem struct {
	queueItem *list.Element
	idxItem   *list.Element
	value     Item
}

type Pool struct {
	MaxLength int
	New       func(arg interface{}) (Item, error)

	len   int
	queue list.List
	idx   map[interface{}]*list.List
	mtx   sync.Mutex
}

func (p *Pool) Len() int { return p.len }

func (p *Pool) get(key interface{}) Item {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	group, ok := p.idx[key]
	if !ok || group.Len() == 0 {
		return nil
	}

	i := group.Remove(group.Front()).(*poolItem)
	if group.Len() == 0 {
		delete(p.idx, key)
	}

	p.queue.Remove(i.queueItem)
	p.len--

	return i.value
}

func (p *Pool) Get(key, arg interface{}) (Item, error) {
	if itm := p.get(key); itm != nil {
		return itm, nil
	}

	if p.New != nil {
		// Create new
		return p.New(arg)
	}

	return nil, nil
}

func (p *Pool) Put(item Item) {
	key := item.Key()

	p.mtx.Lock()
	defer p.mtx.Unlock()

	// Remove first item
	if p.MaxLength > 0 && p.len == p.MaxLength {
		i := p.queue.Remove(p.queue.Back()).(*poolItem)
		k := i.value.Key()
		group := p.idx[k]
		group.Remove(i.idxItem)
		if group.Len() == 0 {
			delete(p.idx, k)
		}
		p.len--

		if f, ok := i.value.(Finalizer); ok {
			f.Finalize()
		}
	}

	// Push given item
	i := &poolItem{
		value: item,
	}

	i.queueItem = p.queue.PushFront(i)

	var (
		group *list.List
		ok    bool
	)

	if group, ok = p.idx[key]; !ok {
		if p.idx == nil {
			p.idx = make(map[interface{}]*list.List)
		}

		group = list.New()
		p.idx[key] = group
	}

	i.idxItem = group.PushFront(i)
	p.len++
}

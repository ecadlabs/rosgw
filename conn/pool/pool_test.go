package pool

import (
	"testing"
)

type testItem struct {
	key int
	val int
}

func (i *testItem) Key() interface{} {
	return i.key
}

func TestPool(t *testing.T) {
	pool := Pool{
		MaxLength: 10,
	}

	for i := 0; i < 20; i++ {
		pool.Put(&testItem{
			key: i,
			val: i,
		})
	}

	for i := 0; i < 10; i++ {
		if v, _ := pool.Get(i, nil); v != nil {
			t.Fatalf("got %v; want nil\n", v)
		}
	}

	for i := 10; i < 20; i++ {
		tmp, _ := pool.Get(i, nil)
		itm := tmp.(*testItem)
		if itm.key != i {
			t.Fatalf("got %v; want %v\n", itm.key, i)
		}
	}

	if pool.Len() != 0 {
		t.Fatalf("assumed empty\n")
	}

	for i := 0; i < 5; i++ {
		for j := 0; j < 2; j++ {
			pool.Put(&testItem{
				key: i,
				val: j,
			})
		}
	}

	for i := 0; i < 5; i++ {
		for j := 0; j < 2; j++ {
			tmp, _ := pool.Get(i, nil)
			itm := tmp.(*testItem)

			if itm.key != i {
				t.Fatalf("got %v; want %v\n", itm.key, i)
			}

			if itm.val != 1-j {
				t.Fatalf("got %v; want %v\n", itm.val, 1-j)
			}
		}
	}

	if pool.Len() != 0 {
		t.Fatalf("assumed empty\n")
	}
}

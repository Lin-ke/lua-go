package state

import (
	"luago54/number"
	"math"
)

// lobject.h #724 : hashtable + array
type luaTable struct {
	arr  []luaValue
	_map map[luaValue]luaValue
}

func newLuaTable(nArr, nRec int) *luaTable {
	t := &luaTable{}
	if nArr > 0 {
		t.arr = make([]luaValue, 0, nArr)
	}
	if nRec > 0 {
		t._map = make(map[luaValue]luaValue, nRec)
	}
	return t
}

func (Table *luaTable) len() int {
	return len(Table.arr)
}

// t[1.0] => t[1]
func (Table *luaTable) get(key luaValue) luaValue {
	key = _floatToInteger(key)

	if idx, ok := key.(int64); ok {

		if idx >= 1 && idx <= int64(len(Table.arr)) {
			return Table.arr[idx-1]
		}
	}
	return Table._map[key]
}
func _floatToInteger(key luaValue) luaValue {
	if f, ok := key.(float64); ok {
		if i, ok := number.FloatToInteger(f); ok {
			return i
		}
		// failed
	}
	// if failed, then return key.
	return key
}

// enable inf, not enable nan
func (Table *luaTable) put(key, val luaValue) {
	if key == nil {
		panic("table index is nil!")
	}
	if f, ok := key.(float64); ok && math.IsNaN(f) {
		panic("table index is NaN!")
	}

	key = _floatToInteger(key)
	// x.0 will be x , (x < MaxInt64)
	if idx, ok := key.(int64); ok && idx >= 1 {
		arrLen := int64(len(Table.arr))
		if idx <= arrLen {
			Table.arr[idx-1] = val
			if idx == arrLen && val == nil {
				Table._shrinkArray()

				return
			}

			return
		}

		// expand, unless val is nil.
		if idx == arrLen+1 {
			delete(Table._map, key)
			if val != nil {
				Table.arr = append(Table.arr, val)
				Table._expandArray()
			}

			return
		}
		// idx > arrLen +1 will expand in hashmap

	}
	// else : k-v store in hashmap
	if val != nil {
		if Table._map == nil {
			Table._map = make(map[luaValue]luaValue, 8) // init
		}
		Table._map[key] = val
	} else {
		delete(Table._map, key)
	}

}

func (Table *luaTable) _shrinkArray() {
	i := len(Table.arr) - 1
	for ; i >= 0; i-- {
		if Table.arr[i] != nil {
			break
		}
	}
	Table.arr = Table.arr[0:i]

}

func (Table *luaTable) _expandArray() {
	// from the next.
	for idx := int64(len(Table.arr)) + 1; true; idx++ {
		if val, found := Table._map[idx]; found {
			delete(Table._map, idx)
			Table.arr = append(Table.arr, val)
		} else {
			break
		}
	}
}

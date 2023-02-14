package sortedmap

import (
	"fmt"
	"strings"
)

// template type OrderedKMap(KType,VType,sort)
type KType int
type VType interface{}
type SortFunc func(i, j KType) int

func sort(i, j KType) int {
	return 0
}

type OrderedKMap struct {
	kv map[KType]VType
	ll *mlist
}

func NewOrderedKMap() *OrderedKMap {
	return &OrderedKMap{
		kv: make(map[KType]VType),
		ll: newList(sort),
	}
}

func (m *OrderedKMap) Contains(key KType) bool {
	_, ok := m.kv[key]
	return ok
}

// Get returns the Value for a Key. If the Key does not exist, the second return
// parameter will be false and the Value will be nil.
func (m *OrderedKMap) Get(key KType) (VType, bool) {
	value, ok := m.kv[key]
	if ok {
		return value, true
	}
	return nil, false
}

// Set will set (or replace) a Value for a Key. If the Key was new, then true
// will be returned. The returned Value will be false if the Value was replaced
// (even if the Value was the same).
func (m *OrderedKMap) Set(key KType, value VType) VType {
	_, didExist := m.kv[key]
	m.kv[key] = value
	if !didExist {
		m.ll.Add(key)
	}
	return m.kv[key]
}

// Delete will remove a Key from the map. It will return true if the Key was
// removed (the Key did exist).
func (m *OrderedKMap) Delete(key KType) (didDelete bool) {
	_, ok := m.kv[key]
	if ok {
		m.ll.Remove(key)
		delete(m.kv, key)
	}
	return ok
}

// Len returns the number of elements in the map.
func (m *OrderedKMap) Len() int {
	return len(m.kv)
}

// Keys returns all of the keys in the order they were inserted. If a Key was
// replaced it will retain the same position. To ensure most recently set keys
// are always at the end you must always Delete before Set.
func (m *OrderedKMap) Keys() (keys []KType) {
	keys = make([]KType, 0, m.Len())
	for i := 0; i < m.ll.size; i++ {
		keys = append(keys, m.ll.elements[i])
	}
	return keys
}

func (m *OrderedKMap) Range(f func(key KType, value VType) bool) {
	for i := 0; i < m.ll.size; i++ {
		key := m.ll.elements[i]
		value, ok := m.kv[key]
		if !ok {
			continue
		}
		if !f(key, value) {
			return
		}
	}
}

func (m *OrderedKMap) Clear() {
	m.kv = make(map[KType]VType)
	m.ll.Clear()
}

//------------------------------------------arraylist------------------------------------------------------

type mlist struct {
	elements []KType
	size     int
	sort     SortFunc
}

const (
	growthFactor = float32(2.0)  // growth by 100%
	shrinkFactor = float32(0.25) // shrink when size is 25% of capacity (0 means never shrink)
)

// New instantiates a new list and adds the passed values, if any, to the list
func newList(sort SortFunc) *mlist {
	list := &mlist{}
	list.sort = sort
	return list
}

func (list *mlist) Add(key KType) {
	index := list.insertIndex(key)
	if !list.withinRange(index) {
		// Append
		if index == list.size {
			list.add(key)
		}
		return
	}
	list.size++
	list.growBy(1)
	copy(list.elements[index+1:], list.elements[index:list.size-1])
	list.elements[index] = key
}

// Remove removes the element at the given index from the list.
func (list *mlist) Remove(key KType) {
	index := list.findIndex(key)
	if !list.withinRange(index) {
		return
	}
	//list.elements[index] = 0                                      // cleanup reference
	copy(list.elements[index:], list.elements[index+1:list.size]) // shift to the left by one (slow operation, need ways to optimize this)
	list.size--
	list.shrink()
}

// Values returns all elements in the list.
func (list *mlist) Values() []KType {
	newElements := make([]KType, list.size, list.size)
	copy(newElements, list.elements[:list.size])
	return newElements
}

// Size returns number of elements within the list.
func (list *mlist) Size() int {
	return list.size
}

// Clear removes all elements from the list.
func (list *mlist) Clear() {
	list.size = 0
	list.elements = []KType{}
}

func (list *mlist) add(value KType) {
	list.growBy(1)
	list.elements[list.size] = value
	list.size++
}

// String returns a string representation of container
func (list *mlist) String() string {
	builder := strings.Builder{}
	builder.WriteString("ArrayList\n")
	for _, value := range list.elements[:list.size] {
		builder.WriteString(fmt.Sprintf("%v ", value))
	}
	return builder.String()
}

// Check that the index is within bounds of the list
func (list *mlist) withinRange(index int) bool {
	return index >= 0 && index < list.size
}

func (list *mlist) resize(cap int) {
	newElements := make([]KType, cap, cap)
	copy(newElements, list.elements)
	list.elements = newElements
}

// Expand the array if necessary, i.e. capacity will be reached if we add n elements
func (list *mlist) growBy(n int) {
	// When capacity is reached, grow by a factor of growthFactor and add number of elements
	currentCapacity := cap(list.elements)
	if list.size+n >= currentCapacity {
		newCapacity := int(growthFactor * float32(currentCapacity+n))
		list.resize(newCapacity)
	}
}

// Shrink the array if necessary, i.e. when size is shrinkFactor percent of current capacity
func (list *mlist) shrink() {
	if shrinkFactor == 0.0 {
		return
	}
	// Shrink when size is at shrinkFactor * capacity
	currentCapacity := cap(list.elements)
	if list.size <= int(float32(currentCapacity)*shrinkFactor) {
		list.resize(list.size)
	}
}

func (list *mlist) findIndex(element KType) int {
	bi, mi := 0, 0
	ei := list.size - 1
	for bi <= ei {
		mi = bi + (ei-bi)>>1
		mid := list.elements[mi]
		v := list.sort(element, mid)
		if v < 0 {
			bi = mi + 1
		} else if v > 0 {
			ei = mi - 1
		} else {
			return mi
		}
	}
	return -1
}

func (list *mlist) insertIndex(element KType) int {
	if list.size == 0 {
		return 0
	}
	bi, mi := 0, 0
	ei := list.size - 1
	for bi <= ei {
		mi = bi + (ei-bi)>>1
		mid := list.elements[mi]
		v := list.sort(element, mid)
		if v < 0 {
			bi = mi + 1
		} else {
			ei = mi - 1
		}
	}
	return bi
}

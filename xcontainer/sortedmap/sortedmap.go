package sortedmap

import (
	"fmt"
	"strings"
)

type Map[K comparable, V any] struct {
	kv map[K]V
	ll *mlist[K]
}

// New returns an empty OrderedKMap with the specified sort function.
// Note: sort function must be stable. not safe for concurrent use.
func New[K comparable, V any](sort func(K, K) int) *Map[K, V] {
	return &Map[K, V]{
		kv: make(map[K]V),
		ll: newList(sort),
	}
}

func (m *Map[K, V]) Contains(key K) bool {
	_, ok := m.kv[key]
	return ok
}

// Get returns the Value for a Key. If the Key does not exist, the second return
// parameter will be false and the Value will be nil.
func (m *Map[K, V]) Get(key K) (V, bool) {
	value, ok := m.kv[key]
	if ok {
		return value, true
	}
	var v V
	return v, false
}

// Set will set (or replace) a Value for a Key. If the Key was new, then true
// will be returned. The returned Value will be false if the Value was replaced
// (even if the Value was the same).
func (m *Map[K, V]) Set(key K, value V) V {
	_, didExist := m.kv[key]
	m.kv[key] = value
	if !didExist {
		m.ll.Add(key)
	}
	return m.kv[key]
}

// Delete will remove a Key from the map. It will return true if the Key was
// removed (the Key did exist).
func (m *Map[K, V]) Delete(key K) (didDelete bool) {
	_, ok := m.kv[key]
	if ok {
		m.ll.Remove(key)
		delete(m.kv, key)
	}
	return ok
}

// Len returns the number of elements in the map.
func (m *Map[K, V]) Len() int {
	return len(m.kv)
}

// Keys returns all of the keys in the order they were inserted. If a Key was
// replaced it will retain the same position. To ensure most recently set keys
// are always at the end you must always Delete before Set.
func (m *Map[K, V]) Keys() (keys []K) {
	keys = make([]K, 0, m.Len())
	for i := 0; i < m.ll.size; i++ {
		keys = append(keys, m.ll.elements[i])
	}
	return keys
}

// Range will call the passed function for each Key/Value pair in the map.
// If the function returns false, iteration will stop.
// Note: Avoid using Set or Delete inside the loop, as it may lead to errors in the loop execution.
func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	size := m.ll.size
	for i := 0; i < m.ll.size; i++ {
		key := m.ll.elements[i]
		value, ok := m.kv[key]
		if !ok {
			continue
		}
		if !f(key, value) {
			return
		}
		if size != m.ll.size {
			panic(fmt.Sprintf("map was mutated during iteration: %d != %d", size, m.ll.size))
		}
	}
}

func (m *Map[K, V]) Clear() {
	m.kv = make(map[K]V)
	m.ll.Clear()
}

//------------------------------------------arraylist------------------------------------------------------

type mlist[K comparable] struct {
	elements []K
	size     int
	sort     func(K, K) int
}

const (
	growthFactor = float32(2.0)  // growth by 100%
	shrinkFactor = float32(0.25) // shrink when size is 25% of capacity (0 means never shrink)
)

// New instantiates a new list and adds the passed values, if any, to the list
func newList[K comparable](sort func(K, K) int) *mlist[K] {
	list := &mlist[K]{}
	list.sort = sort
	return list
}

func (list *mlist[K]) Add(key K) {
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
func (list *mlist[K]) Remove(key K) {
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
func (list *mlist[K]) Values() []K {
	newElements := make([]K, list.size, list.size)
	copy(newElements, list.elements[:list.size])
	return newElements
}

// Size returns number of elements within the list.
func (list *mlist[K]) Size() int {
	return list.size
}

// Clear removes all elements from the list.
func (list *mlist[K]) Clear() {
	list.size = 0
	list.elements = []K{}
}

func (list *mlist[K]) add(value K) {
	list.growBy(1)
	list.elements[list.size] = value
	list.size++
}

// String returns a string representation of container
func (list *mlist[K]) String() string {
	builder := strings.Builder{}
	builder.WriteString("ArrayList\n")
	for _, value := range list.elements[:list.size] {
		builder.WriteString(fmt.Sprintf("%v ", value))
	}
	return builder.String()
}

// Check that the index is within bounds of the list
func (list *mlist[K]) withinRange(index int) bool {
	return index >= 0 && index < list.size
}

func (list *mlist[K]) resize(cap int) {
	newElements := make([]K, cap, cap)
	copy(newElements, list.elements)
	list.elements = newElements
}

// Expand the array if necessary, i.e. capacity will be reached if we add n elements
func (list *mlist[K]) growBy(n int) {
	// When capacity is reached, grow by a factor of growthFactor and add number of elements
	currentCapacity := cap(list.elements)
	if list.size+n >= currentCapacity {
		newCapacity := int(growthFactor * float32(currentCapacity+n))
		list.resize(newCapacity)
	}
}

// Shrink the array if necessary, i.e. when size is shrinkFactor percent of current capacity
func (list *mlist[K]) shrink() {
	if shrinkFactor == 0.0 {
		return
	}
	// Shrink when size is at shrinkFactor * capacity
	currentCapacity := cap(list.elements)
	if list.size <= int(float32(currentCapacity)*shrinkFactor) {
		list.resize(list.size)
	}
}

func (list *mlist[K]) findIndex(element K) int {
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

func (list *mlist[K]) insertIndex(element K) int {
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

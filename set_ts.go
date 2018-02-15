package set

import (
	"sync/atomic"
	"fmt"
	"strings"
)

// Set defines a thread safe set data structure.
type Set SetNonTS

// New creates and initialize a new Set. It's accept a variable number of
// arguments to populate the initial set. If nothing passed a Set with zero
// Length is created.
func New(items ...string) *Set {
	s := (*Set)(NewNonTS())

	s.Add(items...)
	return s
}


func (s Set) GetBucketID(item string) int {
	return hash(item, s.Buckets)
}

func (s *Set) GetSet(item string) *set {
	return s.Sets[s.GetBucketID(item)]
}

// Add includes the specified items (one or more) to the set. The underlying
// Set s is modified. If passed nothing it silently returns.
func (s *Set) Add(items ...string) {
	var count int
	var t *set

	for _, item := range items {
		t = s.GetSet(item)

		t.Lock()
		count += t.Add(item)
		t.Unlock()
	}
	atomic.AddUint64(&s.Length, uint64(count))
}

// Remove deletes the specified items from the set.  The underlying Set s is
// modified. If passed nothing it silently returns.
func (s *Set) Remove(items ...string) {
	var count int

	for _, item := range items {
		set := s.GetSet(item)

		set.Lock()
		count += set.Remove(item)
		set.Unlock()
	}
	atomic.AddUint64(&s.Length, ^uint64(count - 1))
}

// Pop  deletes and return an item from the set. The underlying Set s is
// modified. If set is empty, nil is returned.
func (s *Set) Pop() string {
	if s.Size() == 0 {
		return ""
	}

	i := 0
	set := s.Sets[i]
	for set.Size() == 0 {
		i++
		set = s.Sets[i]
	}

		set.Lock()
		res := set.Pop()
		set.Unlock()

		if res != "" {
			atomic.AddUint64(&s.Length, ^uint64(0))
			return res
		}


	return ""
}

// Has looks for the existence of items passed. It returns false if nothing is
// passed. For multiple items it returns true only if all of  the items exist.
func (s *Set) Has(items ...string) bool {
	if s.Size() == 0 {
		return false
	}

	has := true
	for _, item := range items {
		set := s.GetSet(item)

		set.RLock()
		if !set.Has(item) {
			has = false
			set.RUnlock()
			break
		}
		set.RUnlock()
	}
	return has
}

// Size returns the number of items in a set.
func (s *Set) Size() int {
	return int(atomic.LoadUint64(&s.Length))
}

// Clear removes all items from the set.
func (s *Set) Clear() {
	for i := range s.Sets {
		set := s.Sets[i]
		setLength := set.Size()

		set.Lock()
		s.Sets[i] = newSet()
		atomic.AddUint64(&s.Length, ^uint64(setLength-1))
		set.Unlock()
	}
}

// IsEmpty reports whether the Set is empty.
func (s *Set) IsEmpty() bool {
	return s.Size() == 0
}

// IsEqual test whether s and t are the same in Length and have the same items.
func (s *Set) IsEqual(t *Set) bool {
	if s.Size() != t.Size() {
		return false
	}

	for i := 0; i < s.Buckets; i++ {
		if s.Sets[i].Size() != t.Sets[i].Size() {
			return false
		}
	}

	equal := true

	var sSet, tSet *set

Loop:
	for i := 0; i < s.Buckets; i++ {
		sSet = s.Sets[i]
		tSet = t.Sets[i]

		sSet.RLock()
		tSet.RLock()
		for item := range sSet.Storage {
			if _, ok := tSet.Storage[item]; !ok {
				equal = false

				sSet.RUnlock()
				tSet.RUnlock()
				break Loop
			}
		}

		sSet.RUnlock()
		tSet.RUnlock()
	}

	return equal
}

// IsSubset tests whether t is a subset of s.
func (s *Set) IsSubset(t *Set) (subset bool) {
	if t.Size() > s.Size() {
		return false
	}

	subset = true
	for i:=0; i < t.Buckets; i++ {
		for item := range t.Sets[i].Storage {
			if !s.Has(item) {
				subset = false
				break
			}
		}
	}

	return
}

// Each traverses the items in the Set, calling the provided function for each
// set member. Traversal will continue until all items in the Set have been
// visited, or if the closure returns false.
func (s *Set) Each(f func(item string) bool) {
	for i:=0; i < s.Buckets; i++ {
		for item := range s.Sets[i].Storage {
			if !f(item) {
				break
			}
		}
	}
}

// String returns a string representation of s
func (s *Set) String() string {
	return fmt.Sprintf("[%s]", strings.Join(s.List(), ", "))
}

// List returns a slice of all items. There is also StringSlice() and
// IntSlice() methods for returning slices of type string or int.
func (s *Set) List() []string {
	list := make([]string, 0, s.Size())

	for i:=0; i < s.Buckets; i++ {
		for item := range s.Sets[i].Storage {
			list = append(list, item)
		}
	}

	return list
}

// Copy returns a new Set with a copy of s.
func (s *Set) Copy() *Set {
	return New(s.List()...)
}

// Merge is like Union, however it modifies the current set it's applied on
// with the given t set.
func (s *Set) Merge(t *Set) {
	for i:=0; i < s.Buckets; i++ {
		s.Add(t.Sets[i].List()...)
	}
}

// it's not the opposite of Merge.
// Separate removes the set items containing in t from set s. Please aware that
func (s *Set) Separate(t *Set) {
	for i:=0; i < s.Buckets; i++ {
		s.Remove(t.Sets[i].List()...)
	}
}

// IsSuperset tests whether t is a superset of s.
func (s *Set) IsSuperset(t *Set) bool {
	return t.IsSubset(s)
}

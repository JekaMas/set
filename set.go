package set

import (
	"sync"
	"fmt"
	"sync/atomic"
	"strings"
)

// Provides a common set baseline for both threadsafe and non-ts Sets.
type set struct {
	Storage map[string]struct{}
	size    uint64
	sync.RWMutex
}

func newSet() *set {
	s := &set{}
	s.Storage = make(map[string]struct{})

	return s
}

// Add includes the specified items (one or more) to the set. The underlying
// set s is modified. If passed nothing it silently returns.
func (s *set) Add(items ...string) int {
	var count int
	for _, item := range items {
		if _, ok := s.Storage[item]; ok {
			continue
		}
		s.Storage[item] = struct{}{}
		count++
	}
	atomic.AddUint64(&s.size, uint64(count))

	return count
}

func (s *set) add(item string) {
	if _, ok := s.Storage[item]; ok {
		return
	}

	s.Storage[item] = struct{}{}
	atomic.AddUint64(&s.size, 1)
}

// Remove deletes the specified items from the set.  The underlying set s is
// modified. If passed nothing it silently returns.
func (s *set) Remove(items ...string) int {
	var diff int

	for _, item := range items {
		if _, ok := s.Storage[item]; !ok {
			diff++
			continue
		}

		delete(s.Storage, item)
	}
	atomic.AddUint64(&s.size, ^uint64(len(items)-diff-1))

	return len(items)-diff
}

func (s *set) remove(item string) {
	if _, ok := s.Storage[item]; !ok {
		return
	}

	delete(s.Storage, item)
	atomic.AddUint64(&s.size, ^uint64(0))
}

// Pop  deletes and return an item from the set. The underlying set s is
// modified. If set is empty, nil is returned.
func (s *set) Pop() string {
	for item := range s.Storage {
		s.remove(item)
		return item
	}
	return ""
}

// Has looks for the existence of items passed. It returns false if nothing is
// passed. For multiple items it returns true only if all of  the items exist.
func (s *set) Has(items ...string) bool {
	has := true
	for _, item := range items {
		if _, has = s.Storage[item]; !has {
			break
		}
	}
	return has
}

// Size returns the number of items in a set.
func (s *set) Size() int {
	return int(atomic.LoadUint64(&s.size))
}

// Clear removes all items from the set.
func (s *set) Clear() {
	s.Storage = make(map[string]struct{})
	atomic.StoreUint64(&s.size, 0)
}

// IsEmpty reports whether the set is empty.
func (s *set) IsEmpty() bool {
	return s.Size() == 0
}

// IsEqual test whether s and t are the same in Length and have the same items.
func (s *set) IsEqual(t *set) bool {
	if s.Size() != t.Size() {
		return false
	}

	equal := true
	t.Each(func(item string) bool {
		_, equal = s.Storage[item]
		return equal // if false, Each() will end
	})

	return equal
}

// IsSubset tests whether t is a subset of s.
func (s *set) IsSubset(t *set) (subset bool) {
	subset = true

	t.Each(func(item string) bool {
		_, subset = s.Storage[item]
		return subset
	})

	return
}

// IsSuperset tests whether t is a superset of s.
func (s *set) IsSuperset(t *set) bool {
	return t.IsSubset(s)
}

// Each traverses the items in the set, calling the provided function for each
// set member. Traversal will continue until all items in the set have been
// visited, or if the closure returns false.
func (s *set) Each(f func(item string) bool) {
	for item := range s.Storage {
		if !f(item) {
			break
		}
	}
}

// String returns a string representation of s
func (s *set) String() string {
	return fmt.Sprintf("[%s]", strings.Join(s.List(), ", "))
}

// List returns a slice of all items. There is also StringSlice() and
// IntSlice() methods for returning slices of type string or int.
func (s *set) List() []string {
	list := make([]string, 0, len(s.Storage))

	for item := range s.Storage {
		list = append(list, item)
	}

	return list
}

// Copy returns a new set with a copy of s.
func (s *set) Copy() *set {
	clone := newSet()
	clone.Add(s.List()...)
	return clone
}

// Merge is like Union, however it modifies the current set it's applied on
// with the given t set.
func (s *set) Merge(t *set) {
	t.Each(func(item string) bool {
		s.add(item)
		return true
	})
}

// it's not the opposite of Merge.
// Separate removes the set items containing in t from set s. Please aware that
func (s *set) Separate(t *set) {
	s.Remove(t.List()...)
}
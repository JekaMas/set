package set

import (
	"fmt"
	"strings"
	"sync/atomic"
)

// Provides a common set baseline for both threadsafe and non-ts Sets.
type set struct {
	m map[string]struct{} // struct{} doesn't take up space
	size uint64
}

// SetNonTS defines a non-thread safe set data structure.
type SetNonTS struct {
	set
}

// NewNonTS creates and initialize a new non-threadsafe Set.
// It accepts a variable number of arguments to populate the initial set.
// If nothing is passed a SetNonTS with zero size is created.
func NewNonTS(items ...string) *SetNonTS {
	s := &SetNonTS{}
	s.m = make(map[string]struct{})

	s.Add(items...)
	return s
}

// New creates and initalizes a new Set interface. It accepts a variable
// number of arguments to populate the initial set. If nothing is passed a
// zero size Set based on the struct is created.
func (s *set) New(items ...string) Interface {
	return NewNonTS(items...)
}

// Add includes the specified items (one or more) to the set. The underlying
// Set s is modified. If passed nothing it silently returns.
func (s *set) Add(items ...string) {
	var count int
	for _, item := range items {
		if _, ok := s.m[item]; ok {
			continue
		}
		s.m[item] = struct{}{}
		count++
	}
	atomic.AddUint64(&s.size, uint64(count))
}

func (s *set) add(item string) {
	if _, ok := s.m[item]; ok {
		return
	}

	s.m[item] = struct{}{}
	atomic.AddUint64(&s.size, 1)
}

// Remove deletes the specified items from the set.  The underlying Set s is
// modified. If passed nothing it silently returns.
func (s *set) Remove(items ...string) {
	var diff int

	for _, item := range items {
		if _, ok := s.m[item]; !ok {
			diff++
			continue
		}

		delete(s.m, item)
	}
	atomic.AddUint64(&s.size, ^uint64(len(items)-diff-1))
}

func (s *set) remove(item string) {
	if _, ok := s.m[item]; !ok {
		return
	}

	delete(s.m, item)
	atomic.AddUint64(&s.size, ^uint64(0))
}

// Pop  deletes and return an item from the set. The underlying Set s is
// modified. If set is empty, nil is returned.
func (s *set) Pop() string {
	for item := range s.m {
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
		if _, has = s.m[item]; !has {
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
	s.m = make(map[string]struct{})
	atomic.StoreUint64(&s.size, 0)
}

// IsEmpty reports whether the Set is empty.
func (s *set) IsEmpty() bool {
	return s.Size() == 0
}

// IsEqual test whether s and t are the same in size and have the same items.
func (s *set) IsEqual(t Interface) bool {
	if s.Size() != t.Size() {
		return false
	}

	equal := true
	t.Each(func(item string) bool {
		_, equal = s.m[item]
		return equal // if false, Each() will end
	})

	return equal
}

// IsSubset tests whether t is a subset of s.
func (s *set) IsSubset(t Interface) (subset bool) {
	subset = true

	t.Each(func(item string) bool {
		_, subset = s.m[item]
		return subset
	})

	return
}

// IsSuperset tests whether t is a superset of s.
func (s *set) IsSuperset(t Interface) bool {
	return t.IsSubset(s)
}

// Each traverses the items in the Set, calling the provided function for each
// set member. Traversal will continue until all items in the Set have been
// visited, or if the closure returns false.
func (s *set) Each(f func(item string) bool) {
	for item := range s.m {
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
	list := make([]string, 0, len(s.m))

	for item := range s.m {
		list = append(list, item)
	}

	return list
}

// Copy returns a new Set with a copy of s.
func (s *set) Copy() Interface {
	return NewNonTS(s.List()...)
}

// Merge is like Union, however it modifies the current set it's applied on
// with the given t set.
func (s *set) Merge(t Interface) {
	t.Each(func(item string) bool {
		s.add(item)
		return true
	})
}

// it's not the opposite of Merge.
// Separate removes the set items containing in t from set s. Please aware that
func (s *set) Separate(t Interface) {
	s.Remove(t.List()...)
}

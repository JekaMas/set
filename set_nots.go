package set

import (
	"fmt"
	"sync/atomic"
	"strings"
	"math/rand"
)

// SetNonTS defines a non-thread safe set data structure.
type SetNonTS struct {
	Sets    []*set
	Buckets int
	Length  uint64
}

// NewNonTS creates and initialize a new non-threadsafe Set.
// It accepts a variable number of arguments to populate the initial set.
// If nothing is passed a SetNonTS with zero Length is created.
func NewNonTS(items ...string) *SetNonTS {
	buckets := 64
	s := &SetNonTS{}
	s.Buckets = buckets
	s.Sets = make([]*set, buckets)

	for i := range s.Sets {
		s.Sets[i] = newSet()
	}

	s.Add(items...)

	return s
}

func (s SetNonTS) GetBucketID(item string) int {
	return hash(item, s.Buckets)
}

func (s *SetNonTS) GetSet(item string) *set {
	return s.Sets[s.GetBucketID(item)]
}

// Add includes the specified items (one or more) to the set. The underlying
// Set s is modified. If passed nothing it silently returns.
func (s *SetNonTS) Add(items ...string) {
	var count int

	for _, item := range items {
		set := s.GetSet(item)
		count += set.Add(item)
	}
	atomic.AddUint64(&s.Length, uint64(count))
}

// Remove deletes the specified items from the set.  The underlying Set s is
// modified. If passed nothing it silently returns.
func (s *SetNonTS) Remove(items ...string) {
	var count int

	for _, item := range items {
		set := s.GetSet(item)
		count += set.Remove(item)
	}
	atomic.AddUint64(&s.Length, ^uint64(count - 1))
}

// Pop  deletes and return an item from the set. The underlying Set s is
// modified. If set is empty, nil is returned.
func (s *SetNonTS) Pop() string {
	i := rand.Intn(s.Buckets)

	for true {
		if i >= s.Buckets {
			return ""
		}
		set := s.Sets[i]
		if res := set.Pop(); res != "" {
			atomic.AddUint64(&s.Length, ^uint64(0))
			return res
		}
		i++
	}

	return ""
}

// Has looks for the existence of items passed. It returns false if nothing is
// passed. For multiple items it returns true only if all of  the items exist.
func (s *SetNonTS) Has(items ...string) bool {
	if s.Size() == 0 {
		return false
	}

	has := true
	for _, item := range items {
		set := s.GetSet(item)
		if !set.Has(item) {
			has = false
			break
		}
	}
	return has
}

// Size returns the number of items in a set.
func (s *SetNonTS) Size() int {
	return int(atomic.LoadUint64(&s.Length))
}

// Clear removes all items from the set.
func (s *SetNonTS) Clear() {
	s.Sets = make([]*set, s.Buckets)

	for i := range s.Sets {
		s.Sets[i] = newSet()
	}
	atomic.StoreUint64(&s.Length, 0)
}

// IsEmpty reports whether the Set is empty.
func (s *SetNonTS) IsEmpty() bool {
	return s.Size() == 0
}

// IsEqual test whether s and t are the same in Length and have the same items.
func (s *SetNonTS) IsEqual(t *SetNonTS) bool {
	if s.Size() != t.Size() {
		return false
	}

	equal := true
	for i := 0; i < s.Buckets; i++ {
		for item := range s.Sets[i].Storage {
			if !t.Has(item) {
				equal = false
				break
			}
		}
	}

	return equal
}

// IsSubset tests whether t is a subset of s.
func (s *SetNonTS) IsSubset(t *SetNonTS) (subset bool) {
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

// IsSuperset tests whether t is a superset of s.
func (s *SetNonTS) IsSuperset(t *SetNonTS) bool {
	return t.IsSubset(s)
}

// Each traverses the items in the Set, calling the provided function for each
// set member. Traversal will continue until all items in the Set have been
// visited, or if the closure returns false.
func (s *SetNonTS) Each(f func(item string) bool) {
	for i:=0; i < s.Buckets; i++ {
		for item := range s.Sets[i].Storage {
			if !f(item) {
				break
			}
		}
	}
}

// String returns a string representation of s
func (s *SetNonTS) String() string {
	return fmt.Sprintf("[%s]", strings.Join(s.List(), ", "))
}

// List returns a slice of all items. There is also StringSlice() and
// IntSlice() methods for returning slices of type string or int.
func (s *SetNonTS) List() []string {
	list := make([]string, 0, s.Size())

	for i:=0; i < s.Buckets; i++ {
		for item := range s.Sets[i].Storage {
			list = append(list, item)
		}
	}

	return list
}

// Copy returns a new Set with a copy of s.
func (s *SetNonTS) Copy() *SetNonTS {
	return NewNonTS(s.List()...)
}

// Merge is like Union, however it modifies the current set it's applied on
// with the given t set.
func (s *SetNonTS) Merge(t *SetNonTS) {
	for i:=0; i < s.Buckets; i++ {
		s.Add(t.Sets[i].List()...)
	}
}

// it's not the opposite of Merge.
// Separate removes the set items containing in t from set s. Please aware that
func (s *SetNonTS) Separate(t *SetNonTS) {
	for i:=0; i < s.Buckets; i++ {
		s.Remove(t.Sets[i].List()...)
	}
}

// Package set provides both threadsafe and non-threadsafe implementations of
// a generic set data structure. In the threadsafe set, safety encompasses all
// operations on one set. Operations on multiple sets are consistent in that
// the elements of each set used was valid at exactly one point in time
// between the start and the end of the operation.
package set

// Interface is describing a Set. Sets are an unordered, unique list of values.
type Interface interface {
	New(items ...string) Interface
	Add(items ...string)
	Remove(items ...string)
	Pop() string
	Has(items ...string) bool
	Size() int
	Clear()
	IsEmpty() bool
	IsEqual(s Interface) bool
	IsSubset(s Interface) bool
	IsSuperset(s Interface) bool
	Each(func(string) bool)
	String() string
	List() []string
	Copy() Interface
	Merge(s Interface)
	Separate(s Interface)
}

// Union is the merger of multiple sets. It returns a new set with all the
// elements present in all the sets that are passed.
//
// The dynamic type of the returned set is determined by the first passed set's
// implementation of the New() method.
func Union(sets ...Interface) Interface {
	result := sets[0].New()

	for _, set := range sets {
		set.Each(func(item string) bool {
			if !result.Has(item) {
				result.Add(item)
			}

			return true
		})
	}

	return result
}

// Difference returns a new set which contains items which are in in the first
// set but not in the others. Unlike the Difference() method you can use this
// function separately with multiple sets.
func Difference(sets ...Interface) Interface {
	result := sets[0].New()

	sets[0].Each(func(item string) bool {
		inDifference := true
		for i, set := range sets {
			if i == 0 {
				continue
			}

			if set.Has(item) {
				inDifference = false
				break
			}
		}
		if inDifference {
			result.Add(item)
		}
		return true
	})
	return result
}

// Intersection returns a new set which contains items that only exist in all given sets.
func Intersection(sets ...Interface) Interface {
	result := sets[0].New()
	smallestIndex := getSmallestSet(sets...)

	sets[smallestIndex].Each(func(item string) bool {
		inIntersection := true
		for i, set := range sets {
			if i == smallestIndex {
				continue
			}

			if !set.Has(item) {
				inIntersection = false
				break
			}
		}
		if inIntersection {
			result.Add(item)
		}
		return true
	})
	return result
}

func getSmallestSet(sets ...Interface) int {
	var smallestLen int
	var smallestIndex int
	var setSize int
	for i, set := range sets {
		setSize = set.Size()
		if set.Size() < smallestLen || smallestLen == 0 {
			smallestLen = setSize
			smallestIndex = i
		}
	}

	return smallestIndex
}

// SymmetricDifference returns a new set which s is the difference of items which are in
// one of either, but not in both.
func SymmetricDifference(s Interface, t Interface) Interface {
	u := Difference(s, t)
	v := Difference(t, s)
	return Union(u, v)
}

// StringSlice is a helper function that returns a slice of strings of s. If
// the set contains mixed types of items only items of type string are returned.
func StringSlice(s Interface) []string {
	var slice []string
	for _, item := range s.List() {
		slice = append(slice, item)
	}
	return slice
}

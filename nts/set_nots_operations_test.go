package nts

import (
	"reflect"
	"testing"
	"strconv"
	"github.com/JekaMas/set"
)

func Test_NonTS_Union(t *testing.T) {
	s := set.NewNonTS("1", "2", "3")
	r := set.NewNonTS("3", "4", "5")
	x := set.NewNonTS("5", "6", "7")

	u := Union(s, r, x)
	if u.Size() != 7 {
		t.Error("Union: the merged set doesn't have all items in it.")
	}

	if !u.Has("1", "2", "3", "4", "5", "6", "7") {
		t.Error("Union: merged items are not availabile in the set.")
	}

	z := Union(x, r)
	if z.Size() != 5 {
		t.Error("Union: Union of 2 Sets doesn't have the proper number of items.")
	}
}

func Test_NonTS_Difference(t *testing.T) {
	s := set.NewNonTS("1", "2", "3")
	r := set.NewNonTS("3", "4", "5")
	x := set.NewNonTS("5", "6", "7")
	u := Difference(s, r, x)

	if u.Size() != 2 {
		t.Error("Difference: the set doesn't have all items in it.")
	}

	if !u.Has("1", "2") {
		t.Error("Difference: items are not availabile in the set.")
	}

	y := Difference(r, r)
	if y.Size() != 0 {
		t.Error("Difference: size should be zero")
	}

}

func Test_NonTS_Intersection(t *testing.T) {
	s1 := set.NewNonTS("1", "3", "4", "5")
	s2 := set.NewNonTS("2", "3", "5", "6")
	s3 := set.NewNonTS("4", "5", "6", "7")
	u := Intersection(s1, s2, s3)

	if u.Size() != 1 {
		t.Error("Intersection: the set doesn't have all items in it.", u.List())
	}

	if !u.Has("5") {
		t.Error("Intersection: items after intersection are not availabile in the set.")
	}
}

func Test_NonTS_SymmetricDifference(t *testing.T) {
	s := set.NewNonTS("1", "2", "3")
	r := set.NewNonTS("3", "4", "5")
	u := SymmetricDifference(s, r)

	if u.Size() != 4 {
		t.Error("SymmetricDifference: the set doesn't have all items in it.")
	}

	if !u.Has("1", "2", "4", "5") {
		t.Error("SymmetricDifference: items are not availabile in the set.")
	}
}

func Test_NonTS_StringSlice(t *testing.T) {
	s := set.NewNonTS("san francisco", "istanbul", "3.14", "1321", "ankara")
	u := StringSlice(s)

	if len(u) != 5 {
		t.Error("StringSlice: slice should only have three items")
	}

	for _, item := range u {
		r := reflect.TypeOf(item)
		if r.Kind().String() != "string" {
			t.Error("StringSlice: slice item should be a string")
		}
	}
}

func Benchmark_NonTS_SetEquality(b *testing.B) {
	s := set.NewNonTS()
	u := set.NewNonTS()

	for i := 0; i < b.N; i++ {
		v := strconv.Itoa(i)
		s.Add(v)
		u.Add(v)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.IsEqual(u)
	}
}

func Benchmark_NonTS_Subset(b *testing.B) {
	s := set.NewNonTS()
	u := set.NewNonTS()

	for i := 0; i < b.N; i++ {
		v := strconv.Itoa(i)
		s.Add(v)
		u.Add(v)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.IsSubset(u)
	}
}

func benchmarkNonTSIntersection(b *testing.B, numberOfItems int) {
	s1 := set.NewNonTS()
	s2 := set.NewNonTS()

	for i := 0; i < numberOfItems/2; i++ {
		v := strconv.Itoa(i)
		s1.Add(v)
	}
	for i := 0; i < numberOfItems; i++ {
		v := strconv.Itoa(i)
		s2.Add(v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Intersection(s1, s2)
	}
}

func Benchmark_NonTS_Intersection10(b *testing.B) {
	benchmarkNonTSIntersection(b, 10)
}

func Benchmark_NonTS_Intersection100(b *testing.B) {
	benchmarkNonTSIntersection(b, 100)
}

func Benchmark_NonTS_Intersection1000(b *testing.B) {
	benchmarkNonTSIntersection(b, 1000)
}

func Benchmark_NonTS_Intersection10000(b *testing.B) {
	benchmarkNonTSIntersection(b, 10000)
}

func Benchmark_NonTS_Intersection100000(b *testing.B) {
	benchmarkNonTSIntersection(b, 100000)
}

func Benchmark_NonTS_Intersection1000000(b *testing.B) {
	benchmarkNonTSIntersection(b, 1000000)
}
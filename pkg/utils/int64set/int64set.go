// int64 set
// not concurrent safety
package int64set

var exists = struct{}{}

type Set struct {
	m map[int64]struct{}
}

type EnumerateFunc func(item int64)

func NewSet(items...int64) *Set {
	s := &Set{
		m: make(map[int64]struct{}),
	}
	s.Add(items...)
	return s
}

func (s *Set) Add(values ...int64) {
	for _, v := range values {
		s.m[v] = exists
	}
}

func (s *Set) Remove(values ...int64) {
	for _, v := range values {
		delete(s.m, v)
	}
}

func (s *Set) Contains(value int64) bool {
	_, c := s.m[value]
	return c
}

func (s *Set) Enumerate(f EnumerateFunc) {
	for key := range s.m {
		f(key)
	}
}

func (s *Set) Merge(another *Set) {
	another.Enumerate(func(items int64) {
		s.Add(items)
	})
}

func (s *Set) GetSlice() []int64 {
	slice := make([]int64, 0 ,len(s.m))
	s.Enumerate(func(items int64) {
		slice = append(slice, items)
	})
	return slice
}

func (s *Set) Length() int {
	return len(s.m)
}

// string set
// not concurrent safety
package stringset

var exists = struct{}{}

type Set struct {
	m map[string]struct{}
}

type EnumerateFunc func(str string)

func NewSet(strs...string) *Set {
	s := &Set{
		m: make(map[string]struct{}),
	}
	s.Add(strs...)
	return s
}

func (s *Set) Add(values ...string) {
	for _, v := range values {
		s.m[v] = exists
	}
}

func (s *Set) Remove(values ...string) {
	for _, v := range values {
		delete(s.m, v)
	}
}

func (s *Set) Contains(value string) bool {
	_, c := s.m[value]
	return c
}

func (s *Set) Enumerate(f EnumerateFunc) {
	for key := range s.m {
		f(key)
	}
}

func (s *Set) Merge(another *Set) {
	another.Enumerate(func(str string) {
		s.Add(str)
	})
}

func (s *Set) GetSlice() []string {
	slice := make([]string, 0 ,len(s.m))
	s.Enumerate(func(str string) {
		slice = append(slice, str)
	})
	return slice
}

func (s *Set) Length() int {
	return len(s.m)
}

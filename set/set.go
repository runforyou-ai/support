package set

import (
	"cmp"
	"slices"
)

// Set is an unordered collection of distinct values. The zero value is an
// empty set ready to use.
type Set[E comparable] map[E]struct{}

// Of returns a new set holding values.
func Of[E comparable](values ...E) Set[E] {
	return Collect(values)
}

// Collect returns a new set holding the elements of s.
func Collect[E comparable](s []E) Set[E] {
	out := make(Set[E], len(s))
	for _, v := range s {
		out[v] = struct{}{}
	}
	return out
}

// CollectBy returns a new set holding key(element) for each element of s.
func CollectBy[E any, K comparable](s []E, key func(E) K) Set[K] {
	out := make(Set[K], len(s))
	for _, v := range s {
		out[key(v)] = struct{}{}
	}
	return out
}

// Sorted returns the elements of s in ascending order, or nil when s is
// empty.
func Sorted[E cmp.Ordered](s Set[E]) []E {
	out := s.Values()
	slices.Sort(out)
	return out
}

// Add adds v to the set and reports whether it was not already present. It
// allocates the map when the set is nil, so the zero value of a Set variable
// or field can be used directly.
func (s *Set[E]) Add(v E) bool {
	if *s == nil {
		*s = make(Set[E])
	}
	if _, ok := (*s)[v]; ok {
		return false
	}
	(*s)[v] = struct{}{}
	return true
}

// Has reports whether v is in the set.
func (s Set[E]) Has(v E) bool {
	_, ok := s[v]
	return ok
}

// Delete removes v from the set. It does nothing when v is absent or the set
// is nil.
func (s Set[E]) Delete(v E) {
	delete(s, v)
}

// Len returns the number of elements in the set.
func (s Set[E]) Len() int {
	return len(s)
}

// Clone returns a new set holding the elements of s. It never returns nil.
func (s Set[E]) Clone() Set[E] {
	out := make(Set[E], len(s))
	for v := range s {
		out[v] = struct{}{}
	}
	return out
}

// Union returns a new set holding the elements that are in s, o or both.
func (s Set[E]) Union(o Set[E]) Set[E] {
	out := make(Set[E], max(len(s), len(o)))
	for v := range s {
		out[v] = struct{}{}
	}
	for v := range o {
		out[v] = struct{}{}
	}
	return out
}

// Intersect returns a new set holding the elements that are in both s and o.
func (s Set[E]) Intersect(o Set[E]) Set[E] {
	small, large := s, o
	if len(small) > len(large) {
		small, large = large, small
	}
	out := make(Set[E])
	for v := range small {
		if _, ok := large[v]; ok {
			out[v] = struct{}{}
		}
	}
	return out
}

// Diff returns a new set holding the elements of s that are not in o.
func (s Set[E]) Diff(o Set[E]) Set[E] {
	out := make(Set[E])
	for v := range s {
		if _, ok := o[v]; !ok {
			out[v] = struct{}{}
		}
	}
	return out
}

// Values returns the elements of the set in unspecified order, or nil when
// the set is empty.
func (s Set[E]) Values() []E {
	if len(s) == 0 {
		return nil
	}
	out := make([]E, 0, len(s))
	for v := range s {
		out = append(out, v)
	}
	return out
}

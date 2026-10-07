package arr

// Unique returns the elements of s without duplicates, keeping the first
// occurrence of each and the original order.
func Unique[S ~[]E, E comparable](s S) S {
	return UniqueBy(s, func(v E) E { return v })
}

// UniqueBy returns the elements of s whose fn(element) has not appeared
// before, keeping the first occurrence of each key and the original order.
func UniqueBy[S ~[]E, E any, K comparable](s S, fn func(E) K) S {
	var out S
	seen := make(map[K]struct{}, len(s))
	for _, v := range s {
		k := fn(v)
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		out = append(out, v)
	}
	return out
}

// Diff returns the elements of s that appear in none of others, in order.
// Duplicates in s are kept.
func Diff[S ~[]E, E comparable](s S, others ...[]E) S {
	exclude := make(map[E]struct{})
	for _, o := range others {
		for _, v := range o {
			exclude[v] = struct{}{}
		}
	}
	return Filter(s, func(v E) bool {
		_, ok := exclude[v]
		return !ok
	})
}

// Intersect returns the elements of s that appear in every one of others, in
// order. Duplicates in s are kept. With no others it returns a copy of s.
func Intersect[S ~[]E, E comparable](s S, others ...[]E) S {
	sets := make([]map[E]struct{}, len(others))
	for i, o := range others {
		set := make(map[E]struct{}, len(o))
		for _, v := range o {
			set[v] = struct{}{}
		}
		sets[i] = set
	}
	return Filter(s, func(v E) bool {
		for _, set := range sets {
			if _, ok := set[v]; !ok {
				return false
			}
		}
		return true
	})
}

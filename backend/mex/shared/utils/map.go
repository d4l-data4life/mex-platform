package utils

func KeysOfMap[K comparable, T any](m map[K]T) []K {
	stringKeys := make([]K, len(m))
	i := 0
	for k := range m {
		stringKeys[i] = k
		i++
	}
	return stringKeys
}

func Map[A any, B any](as []A, f func(x A) B) []B {
	bs := make([]B, len(as))
	for i, a := range as {
		bs[i] = f(a)
	}
	return bs
}

func Filter[T any](ts []T, f func(t T) bool) []T {
	r := []T{}
	for _, t := range ts {
		if f(t) {
			r = append(r, t)
		}
	}
	return r
}

func GetOrZero[T any](ts []T, idx int) T {
	var zero T
	if idx < 0 {
		panic("index is negative")
	}

	if len(ts) < idx+1 {
		return zero
	}

	return ts[idx]
}

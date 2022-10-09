package common

func Map[TSource any, TTarget any](source []TSource, f func(TSource) TTarget) []TTarget {
	res := make([]TTarget, len(source))
	for i, item := range source {
		res[i] = f(item)
	}

	return res
}

func MapWithIndex[TSource any, TTarget any](source []TSource, f func(int, TSource) TTarget) []TTarget {
	res := make([]TTarget, len(source))
	for i, item := range source {
		res[i] = f(i, item)
	}

	return res
}

func Sum[TSource any](source []TSource, f func(TSource) float64) float64 {
	var res float64
	res = 0

	for _, item := range source {
		res += f(item)
	}

	return res
}

func Filter[TSource any](source []TSource, filter func(TSource) bool) []TSource {
	res := make([]TSource, 0)

	for _, item := range source {
		if filter(item) {
			res = append(res, item)
		}
	}

	return res
}

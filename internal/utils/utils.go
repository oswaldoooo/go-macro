package utils

func SliceConvert[T1, T2 any](src []T1, dst []T2, convfunc func(src T1, dst *T2)) {
	for i := range src {
		convfunc(src[i], &dst[i])
	}
}

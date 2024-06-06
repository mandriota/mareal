package utils

func MapCopyNoOverwrite[K, V comparable](dst, src map[K]V) {
	for k, v := range src {
		if _, ok := dst[k]; !ok {
			dst[k] = v
		}
	}
}

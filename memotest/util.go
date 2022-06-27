package memotest

func removeNils[T any](els ...any) []T {
	ret := make([]T,0,len(els))
	for _,el := range els {
		if (el != nil) {
			ret = append(ret, el.(T))
		}
	}
	return ret
}

func limitate[T int | uint | uint8 | ~string ](ptr *T, min T, max T) {
	if( (*ptr) < min) {
		(*ptr) = min
	} else
	if( (*ptr) > max) {
		(*ptr) = max
	}
}

func min[T int](a,b int) int {
	if(a<b) {
		return a
	}
	return b
}

func max[T int](a,b int) int {
	if(a>b) {
		return a
	}
	return b
}


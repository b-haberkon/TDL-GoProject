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


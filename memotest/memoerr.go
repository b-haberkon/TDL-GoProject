package memotest

func showError(err error) SStr {
	return NewSStr(func(s SStr) bool {
		s <- `{"message":"`
		s <- err.Error()
		s <- `"}`
		return true
	});
}

func showRetError(err error) SStr {
	return NewSStr(func(s SStr) bool {
		s <- `{"error":`
		s.Send(showError(err))
		s <- `}`
		return true
	})
}

func mapSStr[T any](els []T, fn func(el T) SStr) []SStr {
	acc := make([]SStr,len(els))
	for _, el := range els {
		acc = append( acc, fn(el) )
	}
	return acc
}

func showRetErrors(errs []error) SStr {
	if(1 == len(errs)) {
		return showRetError(errs[0]);
	}
	return NewSStr(func(s SStr) bool {
		s <- `{"error":{"message":"Multiple errors","errors":[`
		s.Join( mapSStr[error](errs, showError), "," )
		s <- `]}}`
		return true
	})
}

// Un par de algún tipo T y un error
type WithError[T any] struct {
    val T;
    err error;
}

type RetWithError[T any] chan WithError[T]

func NewRetWithError[T any]() RetWithError[T] {
	return make(RetWithError[T])
}

func (resp RetWithError[T])SendAndClose(pair WithError[T]) RetWithError[T] {
	resp <- pair
	close(resp)
	return resp
}

func (resp RetWithError[T])SendNewAndClose(ret T, err error) RetWithError[T] {
	return resp.SendAndClose(WithError[T]{ret, err})
}


func ErrToJson(resp RetWithError[SStr]) chan string {
	out := make(chan string)
	go func() {
		var recv 	SStr
		intent := <- resp
		if(intent.err != nil) {
			intent.val = showRetError(intent.err)
		} 
		// Si no era ok, intent.val es nil, que es lo buscado
		if(recv != nil) {
			for chunk := range recv {
				out <- chunk
			}
		} // if
		close(out)
	} ()
	return out
}
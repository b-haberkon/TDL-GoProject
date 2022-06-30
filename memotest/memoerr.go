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


/** @brief Recibe un canal resp RetWithError[SStr] (que recibe un stream o un error),
 * y devuelve un canal de strings con una respuesta en Json (la recibida o el error).
 * Lee una única vez el canal resp, si es un error, envía por la salida el error en
 * un strig con un JSON; si no, copia los mismos fragmentos recibidos.
 */
func ErrToJson(resp RetWithError[SStr]) chan string {
	out := make(chan string)
	go func() {
		intent := <- resp
		if(intent.err != nil) {
			intent.val = showRetError(intent.err)
		} 
		// En este punto, intent.val es el stream recibido, o el generado para el error.
		if(intent.val != nil) {
			for chunk := range intent.val {
				out <- chunk
			}
		} // if
		close(out)
	} ()
	return out
}
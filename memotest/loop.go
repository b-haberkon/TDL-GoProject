package memotest

type LoopFn func(*Loop)

type LoopCall struct {
	Fn LoopFn
}

type LoopQueue chan LoopCall

type Loop struct {
	Queue LoopQueue
	Extra any 
}

/** Inicia un bucle de asincronía, y devuelve el canal para la cola de mensajes. **/
func NewLoop(extra any) *Loop {
    loop := Loop{ make(chan LoopCall), extra }
    go func() {
        for call := range loop.Queue {
            call.Fn(&loop)
        }
    } ()
    return &loop
}

func (loop *Loop)Close() {
	if(loop.Queue != nil) {
		close(loop.Queue)
		loop.Queue = nil;
	}
}

/** @brief Llamada a función asíncrona con retorno inmediato.
 * Crea una goroutine que queda bloqueada hasta que la función fn
 * es leída por el bucle. Pero Async() vuelve al momento.
 **/
func (loop *Loop)Async(fn LoopFn) {
	go func() {
		loop.Queue <- LoopCall{fn}
	} ()
}

// Para garantizar que sea la primera
func (loop *Loop)WaitTurn(fn LoopFn) {
	loop.Queue <- LoopCall{fn}
}

func LoopAwaitInto[T any](loop *Loop, fn LoopFn, resp chan T) T {
	loop.Queue <- LoopCall{fn}
	return <- resp
}

func LoopAwait[T any](loop *Loop, fn LoopFn) T {
	return LoopAwaitInto[T](loop, fn, make(chan T))
}

func AsyncVal[T any](val T) chan T {
	resp := make(chan T)
	go func() {
		resp <- val
		close(resp)
	} ()
	return resp
}
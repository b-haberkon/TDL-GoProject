package memotest

import (
	//"fmt"
	"time"
)

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
			//fmt.Printf("Bloqueando bucle %v %v…\n", loop, loop.Extra)
            call.Fn(&loop)
			//fmt.Printf("…Desbloqueando bucle %v %v.\n", loop, loop.Extra)
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

func RespOrTimeout[T any](resp chan T, expiration time.Duration, cbk func() T, discard func(T,bool)) chan T {
	to, _, _ := timeout(expiration)
	ret := make(chan T)
	go func() {
		done := false
		for {
			select {
			case val, ok := <- resp:
				if(done) {
					if (discard != nil) {
						discard(val,ok)
					}
					return
				} else {
					done = true
					ret <- val
				}
			case <- to:
				if(!done) {
					done = true
					ret <- cbk()
				} else  {
					return
				}
			} // select
		} // for
	} ()
	return ret
}

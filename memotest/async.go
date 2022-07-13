package memotest

// Un par de algo y un error
type asyncRet WithError[any]

func (vin WithError[T]) asyncRet() asyncRet {
    return asyncRet{vin.val, vin.err}
}

/*func (vin asyncRet) WithError[T]() WithError[T] {
    return WithError[T]{(vin.val).(T), vin.err}
}*/

// Una cola es un canal de llamadas
type asyncQueue chan asyncCall

// Argumentos
type asyncArgs any

// Contexto para la función
type asyncCtx struct {
    Queue asyncQueue
    Extra any
}

// Una función asíncrona recibe unos argumentos, y un canal al que enviar un par de algo y un error
type asyncFn func(ctx asyncCtx, args asyncArgs, resp chan asyncRet);

// Una llamada asíncorna consiste en:
type asyncCall struct {
    Ret     chan asyncRet       // Un canal al cual responder
    Fn      asyncFn             // Una función a la que llamar
    Args    asyncArgs           // Unos argumentos
}

func wrapAsync[RetType any](ctx asyncCtx, fn asyncFn, args asyncArgs) chan WithError[RetType] {
    resp := make(chan WithError[RetType])
    go func() {
        recv        := make(chan asyncRet)
        ctx.Queue   <- asyncCall{recv, fn, args}
        vals        := <- recv
        resp        <- WithError[RetType]{(vals.val).(RetType), vals.err}
    } ()
    return resp
}

/** Inicia un bucle de asincronía, y devuelve el canal para la cola de mensajes. **/
func asyncLoop(extra any) asyncCtx {
    ctx := asyncCtx{ make(chan asyncCall), extra }
    go func() {
        for call := range ctx.Queue {
            call.Fn(ctx, call.Args, call.Ret)
        }
    } ()
    return ctx
}
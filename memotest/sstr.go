package memotest

import (
	"io"
	"fmt"
)

type SStr chan string

/** @brief Crea y devuelve un stream, y lo pasa
 como argumento a la goruotine dada.
 Si ésta deuvelve true, cierra el stream.
 **/
func NewSStr(fn func(SStr) bool) SStr {
	s := make(SStr)
	go func() {
		if(fn(s)) {
			close(s)
		}
	} ()
	return s
}

func SStr1(val string) SStr {
	return NewSStr(func(s SStr) bool {
		s <- val
		return true
	})
}


func (s SStr)Send(t SStr) SStr {
	for chunk := range t {
		s <- chunk
	}
	return s
}

func (s SStr)Join(inputs []SStr, sep string) SStr {
	if(s == nil) {
		s = make(SStr)
	}
	for i, input := range inputs {
		if(i>0) {
			s <- sep
		}
		for chunk := range input {
			s <- chunk
		}
	}
	return s
}


func (s SStr)Flat() string {
	out := ""
	for chunk := range s {
		out = out + chunk
	}
	return out
}

func (o SStr)AsyncFlat() SStr {
	return NewSStr(func(d SStr) bool {
		d <- o.Flat()
		return true
	})
}

func (o SStr) spy(name string) SStr {
	return NewSStr( func(d SStr) bool {
		for chunk := range o {
			d <- chunk
			fmt.Printf("** SPY %v (%v): %v\n", name, len(chunk), chunk)
		}
		return true
	})
}

func (o SStr) fork(f chan string) SStr {
	return NewSStr( func(d SStr) bool {
		for chunk := range o {
			d <- chunk
			f <- chunk
		}
		return true
	})
}

func NewSStrSpy(name string, num int) chan string {
	fstr := make(chan string)
	go func() {
		S := ""
		i := 0
		for chunk := range fstr {
			S += chunk
			i ++
			if (num>0) && ( (i % num) ==0 ) {
				fmt.Printf("*+ %v (%v,%v): %v\n",name,i,len(S),S)
				S = ""
			}
		}
		fmt.Printf("** %v (%v,%v): %v\n",name,i,len(S),S)
	} ()
	return fstr
}

// Compatible con io.Stream()
type SStrIO struct {
	stream SStr			// Stram del que está leyendo
	extra  *string		// Adicional sin enviar (resto de una lectura anterior)
}

func (s SStr) AsIO() *SStrIO {
	return &SStrIO{s, nil}
}

func (S *SStrIO) helper(p []byte, pchunk *string) (n int, err error) {
	n = copy(p, []byte(*pchunk))
	if (n < len(*pchunk)) {
		out := (*pchunk)[n:]
		S.extra = &out
	} else {
		S.extra = nil
	}
	return n, nil
}

func (S *SStrIO) Read(p []byte) (n int, err error) {
	if(S.extra != nil) {					// Si había quedado algo
		return S.helper(p, S.extra)			// Leo, pero de lo que quedó
	}

	select {
	case chunk, ok := <- S.stream:
		if(ok) {
			return S.helper(p, &chunk)		// Leo lo recibido
		} else {
			return 0, io.EOF				// Fin de stream
		}
	}	// select
}


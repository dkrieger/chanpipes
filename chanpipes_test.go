package chanpipes

import (
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestNew(t *testing.T) {
	out, in := New()
	go func() {
		in <- true
	}()
	assert.Equal(t, true, <-out)
	assert.IsType(t, make(chan<- interface{}), in)
	assert.IsType(t, make(<-chan interface{}), out)
}

func TestTee(t *testing.T) {
	defer func(gmp int) {
		runtime.GOMAXPROCS(gmp)
	}(runtime.GOMAXPROCS(-1))
	runtime.GOMAXPROCS(1)
	out, in := New()
	out, foo := Tee(out)
	out, bar := Tee(out)
	out, baz := Tee(out)
	go func(input chan<- interface{}) {
		input <- "testing"
	}(in)
	<-out
	assert.Equal(t, "testing", <-foo)
	assert.Equal(t, "testing", <-bar)
	assert.Equal(t, "testing", <-baz)
}

func TestPipe(t *testing.T) {
	imap := func(x interface{}) interface{} {
		return x.(int) * 2
	}
	igrep := func(x interface{}) bool {
		switch x.(type) {
		case int:
			return true
		default:
			return false
		}
	}
	smap := func(s interface{}) interface{} {
		return s.(string) + " world"
	}
	sgrep := func(x interface{}) bool {
		switch x.(type) {
		case string:
			return true
		default:
			return false
		}
	}
	build := func() (<-chan interface{}, chan<- interface{}, <-chan interface{}, <-chan interface{}) {
		out, in := New()
		out, foo := Tee(out)
		foo, _ = Grep(foo, igrep)
		foo = Pipe(foo, imap)
		out, bar := Tee(out)
		bar, _ = Grep(bar, sgrep)
		bar = Pipe(bar, smap)
		return out, in, foo, bar
	}
	consume := func(foo <-chan interface{}, bar <-chan interface{}) {
		select {
		case msg := <-foo:
			assert.Equal(t, 6, msg)
		case msg := <-bar:
			assert.Equal(t, "hello world", msg)
		}
	}
	out, in, foo, bar := build()
	go func() {
		in <- 3
	}()
	<-out
	consume(foo, bar)
	out, in, foo, bar = build()
	go func() {
		in <- "hello"
	}()
	<-out
	consume(foo, bar)
}

func doTestGrep(input bool) interface{} {
	out, in := New()
	res := make(chan interface{})
	cond := func(msg interface{}) bool {
		return msg.(bool)
	}
	pass, fail := Grep(out, cond)
	go func() {
		select {
		case <-pass:
			res <- "foo"
		case <-fail:
			res <- "bar"
		}
	}()
	in <- input
	return <-res
}

func TestGrep(t *testing.T) {
	assert.Equal(t, "foo", doTestGrep(true))
	assert.Equal(t, "bar", doTestGrep(false))
}

func TestFanIn(t *testing.T) {
	testInput := func(input chan<- interface{}, output <-chan interface{}) {
		go func() {
			input <- "testing"
		}()
		assert.Equal(t, "testing", <-output)
	}
	foo := make(chan interface{})
	bar := make(chan interface{})
	baz := make(chan interface{})
	all := FanIn(foo, bar, baz)
	testInput(foo, all)
	testInput(bar, all)
	testInput(baz, all)
}

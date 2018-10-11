// chanpipes provides helpers for building networks of channels, using POSIX
// shell pipeline -like semantics where applicable.
package chanpipes

// New creates a new channel and returns a read-only reference ("out") and a
// write-only reference ("in")
func New() (<-chan interface{}, chan<- interface{}) {
	out := make(chan interface{})
	in := out
	return out, in
}

// Tee takes some readable channel as input, forwards it to a new readable
// channel ("out"), then forwards it to another new readable channel ("side")
// after the goroutine wakes up. The idea is to read the final "out" before any
// "side", ensuring every side channel actually gets wired.
func Tee(in <-chan interface{}) (<-chan interface{}, <-chan interface{}) {
	side := make(chan interface{})
	out := make(chan interface{})
	go func(out chan<- interface{}, in <-chan interface{}, side chan<- interface{}) {
		upstream := <-in
		out <- upstream
		side <- upstream
	}(out, in, side)
	return out, side
}

// Pipe takes some readable channel as input and transforms its contents using
// an interface{}-to-interface{} mapper writing results to the new readable
// "out" channel.
func Pipe(in <-chan interface{}, mapper func(interface{}) interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func(out chan<- interface{}, in <-chan interface{}) {
		out <- mapper(<-in)
	}(out, in)
	return out
}

// Grep is a filtering operation. Unlike the behavior of regular grep, which
// filters each line of stdin independently, Grep filters all of "stdin".
// It's more like
//     # mkfifo foo bar && <foo cat >bar &
//     # <input tee foo | grep condition >/dev/null 2>&1 && <bar cat >output
// than
//     # <input grep condition >output
//
// For long-running processes, special care should be taken to ensure
// goroutines that read the returned channel aren't left dangling in the case
// of a failed condition, by using "select" or some other mechanism.
func Grep(in <-chan interface{}, cond func(interface{}) bool) (<-chan interface{}, <-chan bool) {
	out := make(chan interface{})
	ready := make(chan bool)
	go func(out chan<- interface{}, in <-chan interface{}) {
		upstream := <-in
		ready <- true
		if cond(upstream) {
			out <- upstream
		}
	}(out, in)
	return out, ready
}

// FanIn is like a dynamic "select" statement for N readable channels
func FanIn(inputs ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	for _, input := range inputs {
		go func(ch <-chan interface{}) {
			for {
				out <- <-ch
			}
		}(input)
	}
	return out
}

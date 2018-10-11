# chanpipes
--
    import "github.com/dkrieger/chanpipes"

chanpipes provides helpers for building networks of channels, using POSIX shell
pipeline -like semantics where applicable.

Rob Pike makes repeated comparisons between POSIX shell semantics and Go
channel/goroutine semantics in his seminal "Go Concurrency Patterns" talk
(https://talks.golang.org/2012/concurrency.slide#1). Goroutines are like
background processes, and channels are like named pipes (fifos); chanpipes takes
its name from the latter analogy.

Note that, much as shell pipelines use text as the universal interface,
chanpipes pipelines use `interface{}` for maximum interoperability, at the
expense of runtime safety. If/when golang gets generic types, some runtime
safety may be restored; in the meantime, our POSIX shell analogy is even more
literal, and the same approach of validating inputs at runtime should be taken.
If generic types don't come with golang 2, code generation may be used to
implement strongly-typed pipelines.

## Usage

#### func  Cat

```go
func Cat(inputs ...<-chan interface{}) <-chan interface{}
```
Cat is like a dynamic "select" statement for N readable channels. It is a fan-in
pattern, which is like "cat" in a POSIX shell. A notable difference between Cat
and cat is that Cat doesn't care what order input chans are passed, and will
forward messages to the output chan in the order they arrive

#### func  Grep

```go
func Grep(in <-chan interface{}, cond func(interface{}) bool) (<-chan interface{}, <-chan interface{})
```
Grep is a filtering operation. Unlike the behavior of "grep", which filters each
line of stdin independently, Grep filters all of "stdin". It's more like

    # mkfifo foo bar pass fail && <foo cat >bar &
    # <input tee foo | grep condition >/dev/null 2>&1 && <bar cat >pass || <bar cat >fail &

than

    # mkfifo output
    # <input grep condition >output &

#### func  New

```go
func New() (<-chan interface{}, chan<- interface{})
```
New creates a new channel and returns a read-only reference ("out") and a
write-only reference ("in")

#### func  Pipe

```go
func Pipe(in <-chan interface{}, mapper func(interface{}) interface{}) <-chan interface{}
```
Pipe takes some readable channel as input and transforms its contents using an
interface{}-to-interface{} mapper writing results to the new readable "out"
channel.

#### func  Tee

```go
func Tee(in <-chan interface{}) (<-chan interface{}, <-chan interface{})
```
Tee takes some readable channel as input, forwards it to a new readable channel
("out"), then forwards it to another new readable channel ("side") after the
goroutine wakes up. The idea is to read the final "out" before any "side",
ensuring every side channel actually gets wired.

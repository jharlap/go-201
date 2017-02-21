class: center, middle

# Go 201

---

# Topics

- Package Design
- Testing
- Benchmarking
- Profiling

---

# Package Design

---

# API first

`go get -u golang.org/x/tools/cmd/godoc`

`godoc -http=:6060`

???

API-centric package design is easy to see by looking at godoc. Run locally during dev and for private packages.

---

# Naming

- Package name is always lowercase
- Short, meaningful names
- No plurals

`package url`

`package flag`

---

# Organization

- Package by responsibility, not type of thing

```
package authservice

type User struct {...}

func UserByEmail(ctx context.Context, email string) (User, error)
```

- Not `package models`

---

# Multiple files

- Use multiple files within one package to separate logical sub-components

From `github.com/golang/dep`:

```
analyzer.go
analyzer_test.go
context.go
context_test.go
fs.go
fs_test.go
lock.go
lock_test.go
manifest.go
manifest_test.go
project.go
project_test.go
```

???

- Each file generally contains a single type, name matches the filename, and functions relating to that type.
- Tests for each file.

---

# Runnable vs Importable

- Separate runnable commands into `cmd` subtree

    `my.org/newsfeed/cmd/newsfeed` contains `package main`

- Importable packages outside `cmd` tree

    `my.org/newsfeed/authservice`

---

# Documentation

- Package-level documentation

    `Package squeegie implements petroleum removal for penguins`

- Long doc block in separate `doc.go`
- `package main` doc explains the binary

    `Command chainsaw finds and kills zombies`

---

# Hide Async

- Packages expose synchronous API
- Internal async behaviour completely hidden from exposed API
- Caller can manage goroutines/channels to manage concurrent behaviour

---

# Accept interfaces, return structs

```
package sync

type Locker interface { ... }

type Cond struct { ... }
func NewCond(l Locker) *Cond { ... }
```

???

- `Locker` interface is accepted - anyone can make something that satisfies the interface and make a Cond from it
- `Cond` struct is returned - note it doesn't return an interface satisfied by `*Cond` but a concrete value
- Unnecessary complexity if functions are not exported

---

# Design Exercise

- Create a package API for a streaming compression library that will implement run-length encoding (RLE)

- Create a second package for a command line program using the library to compress or uncompress files

- Do not actually implement either package (yet), just design the APIs and document them - all functions should be stubs

---

# Testing

---

# stdlib testing

```
import "testing"

func TestSomething(t *testing.T) {
	want := 3
	got := something()
	if got != want {
		t.Errorf("something() = %d; want %d", got, want)
	}
}
```

???

note the order is actual then expected, both in `if` and `t.Error` output

---

# Table-driven tests

[embedmd]:# (tests/square_test.go /func TestSquareOld/ /\n}/)
```go
func TestSquareOld(t *testing.T) {
	cases := []struct {
		in, want int
	}{
		{1, 1},
		{2, 4},
		{8, 64},
		{0, 0},
		{-1, 1},
	}

	for _, tc := range cases {
		got := Square(tc.in)
		if got != tc.want {
			t.Errorf("Square(%d) = %d; want %d", tc.in, got, tc.want)
		}
	}
}
```

```
--- FAIL: TestSquareOld (0.00s)
	square_test.go:22: Square(2) = 2; want 4
	square_test.go:22: Square(8) = 8; want 64
	square_test.go:22: Square(-1) = -1; want 1
FAIL
```
???

- Useful for repetitive testing of input-output pairs

---

# Table-driven tests with subtests

[embedmd]:# (tests/square_test.go /func TestSquareSubtests/ /\n}/)
```go
func TestSquareSubtests(t *testing.T) {
	cases := []struct {
		in, want int
	}{
		{1, 1},
		{2, 4},
		{8, 64},
		{0, 0},
		{-1, 1},
	}

	for _, tc := range cases {
		t.Run(strconv.Itoa(tc.in), func(t *testing.T) {
			got := Square(tc.in)
			if got != tc.want {
				t.Errorf("got %d; want %d", got, tc.want)
			}
		})
	}
}
```

```
--- FAIL: TestSquareSubtests (0.00s)
    --- FAIL: TestSquareSubtests/2 (0.00s)
    	square_test.go:42: got 2; want 4
    --- FAIL: TestSquareSubtests/8 (0.00s)
    	square_test.go:42: got 8; want 64
    --- FAIL: TestSquareSubtests/-1 (0.00s)
    	square_test.go:42: got -1; want 1
FAIL
```
---

# Testing with dependencies

[embedmd]:# (tests/server.go /type Logger/ /\n}/)
```go
type Logger interface {
	Log(msg string)
}
```

[embedmd]:# (tests/server.go /type Server/ /\n}/)
```go
type Server struct {
	L Logger
}
```

[embedmd]:# (tests/server.go /func/ /$/)
```go
func (s Server) Greet(name string) {
	s.L.Log(fmt.Sprintf("Greetings %s", name))
}
```

---

# Mocks (hand-rolled)

[embedmd]:# (tests/server_test.go /type/ /$/)
```go
type fakeLogger struct {
	captured string
}

func (f *fakeLogger) Log(msg string) {
	f.captured = msg
}
```

[embedmd]:# (tests/server_test.go /func TestServerLogsHand/ /\n}/)
```go
func TestServerLogsHandRolled(t *testing.T) {
	var l fakeLogger
	s := Server{L: &l}

	s.Greet("Alice")

	want := "Hello Alice"
	if l.captured != want {
		t.Errorf(`got "%s"; want "%s"`, l.captured, want)
	}
}
```

```
--- FAIL: TestServerLogsHandRolled (0.00s)
	server_test.go:29: got "Greetings Alice"; want "Hello Alice"
```

---

# Mocks (generated)

[embedmd]:# (tests/server.go /..go:gen/ /\n}/)
```go
//go:generate mockgen -source=$GOFILE -destination=./mock_tests/mock_logger.go Logger
type Logger interface {
	Log(msg string)
}
```

[embedmd]:# (tests/server_test.go /func TestServerLogsGenerated/ /\n}/)
```go
func TestServerLogsGenerated(t *testing.T) {
	ctl := gomock.NewController(t)
	m := mock_tests.NewMockLogger(ctl)
	defer ctl.Finish()

	s := Server{L: m}

	m.EXPECT().Log("Hello Alice")
	s.Greet("Alice")
}
```

```
--- FAIL: TestServerLogsGenerated (0.00s)
	controller.go:113: no matching expected call: *mock_tests.MockLogger.Log([Greetings Alice])
	controller.go:158: missing call(s) to *mock_tests.MockLogger.Log(is equal to Hello Alice)
	controller.go:165: aborting test due to missing call(s)
```

---

# Testing Exercise

- Implement the Reader from the RLE package, making sure to unit test your functions.

- Implement the command line program. Test core functionality, mocking out dependencies (eg: the RLE package). Only implement the Reader, given the RLE package only supports that.

You can generate fake data files with:
```
$ echo -ne '\x48\x01\x45\x01\x4c\x02\x4f\x01' > hello
```

---

# Fuzzing

Generates random mutations of known good inputs to find problematic inputs to your program

[embedmd]:# (fuzz/rle_fuzz.go /func Fuzz/ /\n}/)
```go
func Fuzz(data []byte) int {
	r, err := rle.NewReader(bytes.NewReader(data))
	if err != nil {
		// error handling worked
		return 0
	}

	buf := make([]byte, 64<<10)
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return 0
	}

	// empty in, empty out
	if len(data) == 0 && n == 0 {
		return 0
	}

	// output size matches expectation
	if len(data) > 0 && n <= len(buf) && n == expectedOutputSize(data) {
		return 0
	}

	return 1
}
```

---

# Fuzzing

```
$ go get github.com/dvyukov/go-fuzz/go-fuzz
$ go get github.com/dvyukov/go-fuzz/go-fuzz-build
$ go-fuzz-build github.com/jharlap/go-201/fuzz
$ go-fuzz -bin=./fuzz-fuzz.zip -workdir=fuzz/examples
2017/02/14 14:59:07 slaves: 4, corpus: 5 (3s ago), crashers: 1, restarts: 1/0, execs: 0 (0/sec), cover: 0, uptime: 3s
2017/02/14 14:59:10 slaves: 4, corpus: 5 (6s ago), crashers: 1, restarts: 1/0, execs: 0 (0/sec), cover: 42, uptime: 6s
2017/02/14 14:59:13 slaves: 4, corpus: 5 (9s ago), crashers: 1, restarts: 1/0, execs: 0 (0/sec), cover: 42, uptime: 9s
^C2017/02/14 14:59:13 shutting down...
$ 
```

1 crasher??

---

# Benchmarking

[embedmd]:# (fuzz/rle/reader_bench_test.go /func/ /\n}/)
```go
func BenchmarkReader(b *testing.B) {
	br := bytes.NewReader([]byte{87, 1, 65, 255, 32, 1, 87, 1, 65, 255, 32, 1, 87, 1, 65, 255})
	for n := 0; n < b.N; n++ {
		br.Seek(0, io.SeekStart)
		r, _ := NewReader(br)
		_, err := ioutil.ReadAll(r)
		if err != nil {
			b.Fatalf("unexpected error: %s", err)
		}
	}
}
```

```
$ go test -bench .
BenchmarkReader-4   	 3000000	       541 ns/op
PASS

$ go test -bench . -benchmem
BenchmarkReader-4   	 3000000	       539 ns/op	    1584 B/op	       3 allocs/op
PASS
```

---

# Profiling

---

# CPU time

```
$ go test -bench . -cpuprofile cpu.prof
$ go pprof rle.test cpu.prof
(pprof) web
```

---

# Flame graph

```
$ go get -u github.com/uber/go-torch
$ go test -bench . -cpuprofile cpu.prof
$ go-torch rle.test cpu.prof
```

---

# Memory allocations

```
$ go test -bench=. -memprofile mem.prof -memprofilerate 1
$ go tool pprof --alloc_space rle.test mem.prof
(pprof) web
```

---



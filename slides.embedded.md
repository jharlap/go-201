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

[embedmd]:# (tests/server.go /type Server/ /\n}/)

[embedmd]:# (tests/server.go /func/ /$/)

---

# Mocks (hand-rolled)

[embedmd]:# (tests/server_test.go /type/ /$/)

[embedmd]:# (tests/server_test.go /func TestServerLogsHand/ /\n}/)

```
--- FAIL: TestServerLogsHandRolled (0.00s)
	server_test.go:29: got "Greetings Alice"; want "Hello Alice"
```

---

# Mocks (generated)

[embedmd]:# (tests/server.go /..go:gen/ /\n}/)

[embedmd]:# (tests/server_test.go /func TestServerLogsGenerated/ /\n}/)

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

```
$ go test -bench .
BenchmarkReader-4   	 1000000	      2455 ns/op
PASS

$ go test -bench . -benchmem
BenchmarkReader-4   	 1000000	      2164 ns/op	    4144 B/op	       5 allocs/op
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



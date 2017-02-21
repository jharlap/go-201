package tests

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jharlap/go-201/tests/mock_tests"
)

func TestServerLogsGenerated(t *testing.T) {
	ctl := gomock.NewController(t)
	m := mock_tests.NewMockLogger(ctl)
	defer ctl.Finish()

	s := Server{L: m}

	m.EXPECT().Log("Hello Alice")
	s.Greet("Alice")
}

func TestServerLogsHandRolled(t *testing.T) {
	var l fakeLogger
	s := Server{L: &l}

	s.Greet("Alice")

	want := "Hello Alice"
	if l.captured != want {
		t.Errorf(`got "%s"; want "%s"`, l.captured, want)
	}
}

type fakeLogger struct {
	captured string
}

func (f *fakeLogger) Log(msg string) {
	f.captured = msg
}

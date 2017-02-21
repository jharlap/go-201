package tests

import "fmt"

//go:generate mockgen -source=$GOFILE -destination=./mock_tests/mock_logger.go Logger
type Logger interface {
	Log(msg string)
}

type Server struct {
	L Logger
}

func (s Server) Greet(name string) {
	s.L.Log(fmt.Sprintf("Greetings %s", name))
}

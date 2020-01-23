package server

type Server struct {
	Start func()
	Stop func()
}


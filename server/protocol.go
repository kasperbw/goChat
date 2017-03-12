package server

type Protocol int

const (
	None Protocol = iota
	Login
	Logout
)

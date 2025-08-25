package rocketleague

type Mode uint8

const (
	Soccer Mode = iota
	Pentathlon
)

type Players uint8

const (
	P2x2 Players = iota
	P3x3
)

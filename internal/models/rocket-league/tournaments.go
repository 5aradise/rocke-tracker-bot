package rocketleague

type Mode uint8

const (
	Soccer Mode = iota
	Pentathlon
)

func (m Mode) String() string {
	switch m {
	case Soccer:
		return "Soccer"
	case Pentathlon:
		return "Pentathlon"
	}
	panic("unknown game mode")
}

type Players uint8

const (
	P2x2 Players = iota
	P3x3
)

func (p Players) String() string {
	switch p {
	case P2x2:
		return "2x2"
	case P3x3:
		return "3x3"
	}
	panic("unknown players mode")
}

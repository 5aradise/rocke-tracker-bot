package telegram

import "strconv"

// id implements telebot.Recipient interface
type id struct {
	value int64
	s     string
}

func newID(value int64) id {
	return id{
		value: value,
		s:     strconv.FormatInt(value, 10),
	}
}

func (i id) Recipient() string {
	return i.s
}

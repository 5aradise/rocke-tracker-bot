package telegram

const ( // language codes
	english   = "en" // other (default)
	ukrainian = "ua"
	russian   = "ru"
)

type langString struct {
	other string
	uaru  string
}

func (s langString) in(lang string) string {
	switch lang {
	default:
		return s.other
	case ukrainian, russian:
		return s.uaru
	}
}

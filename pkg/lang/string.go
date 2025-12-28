package lang

type Language uint8

const ( // languages
	Other Language = iota
	English
	Ukrainian
	Russian
)

type String struct {
	en    string
	uaru  string
	other string
}

// creates language specific string, if other is not specified, other = en
func NewString(
	en string,
	uaru string,
	other ...string,
) String {
	ot := en
	if len(other) > 0 {
		ot = other[0]
	}
	return String{
		en:    en,
		uaru:  uaru,
		other: ot,
	}
}

func (s String) In(lang Language) string {
	switch lang {
	case English:
		return s.en
	case Ukrainian, Russian:
		return s.uaru
	case Other:
		return s.other
	}
	panic("unknown language")
}

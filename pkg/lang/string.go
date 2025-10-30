package lang

type Code string

const ( // language codes
	Other     Code = ""
	English   Code = "en"
	Ukrainian Code = "ua"
	Russian   Code = "ru"
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

func (s String) In(lang Code) string {
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

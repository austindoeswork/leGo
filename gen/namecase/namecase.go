package namecase

import (
	"strings"
)

// Name represents a parameter's name
type Name struct {
	Lower      string
	LowerCamel string
	UpperCamel string
}

// New instantiates a new instance of Name
func New(n string) *Name {
	if len(n) == 0 {
		return nil
	}
	if len(n) == 1 {
		return &Name{
			Lower:      strings.ToLower(n),
			LowerCamel: strings.ToLower(n),
			UpperCamel: strings.ToUpper(n),
		}
	}
	lowerCamel := strings.ToLower(n[:1]) + n[1:]
	if strings.ToUpper(n) == n {
		lowerCamel = strings.ToLower(n)
	}
	upperCamel := strings.ToUpper(n[:1]) + n[1:]
	return &Name{
		Lower:      strings.ToLower(n),
		LowerCamel: lowerCamel,
		UpperCamel: upperCamel,
	}
	return nil
}

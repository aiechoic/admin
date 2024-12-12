package errs

import "log"

type Code int

func (c Code) Error() string {
	return codeText[c]
}

var codeText = map[Code]string{
	BadRequest:          "Bad Request",
	Unauthorized:        "Unauthorized",
	InternalServerError: "Internal Server Error",
}

const (
	BadRequest Code = iota + 4000
	Unauthorized
	InternalServerError
)

func (c Code) String() string {
	return codeText[c]
}

func GetCodes() map[Code]string {
	return codeText
}

func SetCodes(codes map[Code]string) {
	for k, v := range codes {
		if _, ok := codeText[k]; ok {
			log.Panicf("code %d already exists", k)
		}
		codeText[k] = v
	}
}

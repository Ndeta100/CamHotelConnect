package api

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func NewError(code int, msg string) Error {
	return Error{
		Code: code,
		Err:  msg,
	}
}

package internal

type RuntimeError struct {
	Token   *token
	Message string
}

func (r RuntimeError) Error() string {
	return r.Token.lexeme + " " + r.Message
}

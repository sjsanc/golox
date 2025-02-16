package errors

import "github.com/sjsanc/golox/internal/token"

type ParseErr struct {
	Token   *token.Token
	Message string
}

func (r ParseErr) Error() string {
	return r.Token.Lexeme + " " + r.Message
}

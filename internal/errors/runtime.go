package errors

import "github.com/sjsanc/golox/internal/token"

type RuntimeErr struct {
	Token   *token.Token
	Message string
}

func NewRuntimeErr(token *token.Token, message string) RuntimeErr {
	return RuntimeErr{
		Token:   token,
		Message: message,
	}
}

func (r RuntimeErr) Error() string {
	return r.Token.Lexeme + " " + r.Message
}

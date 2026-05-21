package helper

import "fmt"

type ResponseStatus string

type ResponseBase struct {
	Status ResponseStatus `json:"status"`
	Error  any            `json:"error,omitempty"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (err *Error) Error() string {
	return fmt.Sprintf("code:%s,message:%s", err.Code, err.Message)
}

const (
	ResponseStatusSuccess ResponseStatus = "success"
	ResponseStatusError   ResponseStatus = "error"
)

var (
	SuccessResponse = ResponseBase{
		Status: ResponseStatusSuccess,
	}
)

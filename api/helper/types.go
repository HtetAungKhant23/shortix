package helper

type ResponseStatus string

type ResponseBase struct {
	Status ResponseStatus `json:"status"`
	Error  any            `json:"error,omitempty"`
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

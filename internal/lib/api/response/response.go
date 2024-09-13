package response

const AliasLength = 6

type Response struct {
	Status string `json: status` //Error, Ok
	Error  string `json:"error,omitempty"`
}

const (
	StatusError = "Error"
	StatusOk    = "Ok"
)

func OK() Response {
	return Response{
		Status: StatusOk,
	}
}

func Error(err string) Response {
	return Response{
		Status: StatusError,
		Error:  err,
	}
}

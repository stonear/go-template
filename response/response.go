package response

type Response struct {
	Code    string `json:"responseCode"`
	Message string `json:"responseMessage"`
	Data    any    `json:"data"`
}

func New(code Code, data any) *Response {
	return &Response{
		Code:    string(code),
		Message: DefaultRCMessage[code],
		Data:    data,
	}
}

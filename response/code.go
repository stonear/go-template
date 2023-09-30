package response

type Code string

const (
	CodeSuccess               Code = "0000"
	CodeInvalidFormat         Code = "0001"
	CodeInvalidMandatoryField Code = "0002"
	CodeNotPermitted          Code = "0003"
	CodeNotFound              Code = "0004"
	CodeGeneralError          Code = "9999"
)

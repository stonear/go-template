package response

var (
	DefaultRCMessage = map[Code]string{
		CodeSuccess:               "Successful",
		CodeInvalidFormat:         "Invalid Field Format",
		CodeInvalidMandatoryField: "Missing or Invalid Format on Mandatory Field",
		CodeNotPermitted:          "Feature Not Allowed At This Time",
		CodeNotFound:              "Data Not Found",
		CodeGeneralError:          "General Error",
	}
)

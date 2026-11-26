package COMMANDHandler

type COMMANDHandler struct{}

func (*COMMANDHandler) Handler(arrayOfBulkStrings []string) string {
	if !IsCOMMANDRequest(arrayOfBulkStrings) {
		return "ERROR: THIS IS NOT VALID COMMAND REQUEST"
	}
	return "*0\r\n"
}

package ECHOHandler

const ECHO = "ECHO"
const THIS_REQUEST_IS_NOT_ECHO_REQUEST_ERROR = "This request is not echo request"

func IsECHORequest(arrayOfBulkStrings []string) bool {
	return len(arrayOfBulkStrings) > 1 && arrayOfBulkStrings[1] == ECHO
}

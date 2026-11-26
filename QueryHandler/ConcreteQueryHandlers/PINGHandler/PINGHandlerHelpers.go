package PINGHandler

const PING = "PING"
const THIS_REQUEST_IS_NOT_PING_REQUEST_ERRROR = "ERROR: THIS REQUEST IS NOT A PING REQUEST"
const BASIC_PONG_RESPONSE = "+PONG\r\n"

func IsPINGRequest(arrayOfBulkStrings []string) bool {
	if len(arrayOfBulkStrings) > 1 && arrayOfBulkStrings[1] == PING {
		return true
	}
	return false
}

func IsEmptyPINGRequest(arrayOfBulkStrings []string) bool {
	return len(arrayOfBulkStrings) == 2
}

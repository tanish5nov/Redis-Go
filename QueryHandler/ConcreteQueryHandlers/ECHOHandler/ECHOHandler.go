package ECHOHandler

import (
	"BuildYourOwnRedis/RESPParser/Builder"
)

type ECHOHandler struct{}

func (*ECHOHandler) Handler(arrayOfBulkStrings []string) string {
	if !IsECHORequest(arrayOfBulkStrings) {
		return THIS_REQUEST_IS_NOT_ECHO_REQUEST_ERROR
	}
	responseArray := []string{}
	for i := 2; i < len(arrayOfBulkStrings); i++ {
		responseArray = append(responseArray, arrayOfBulkStrings[i])
	}
	return Builder.BuildRESPResponse(responseArray)
}

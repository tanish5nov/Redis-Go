package PINGHandler

import (
	"BuildYourOwnRedis/QueryHandler/ConcreteQueryHandlers/PONGHandler"
	"BuildYourOwnRedis/RESPParser/Builder"
	"fmt"
)

type PINGHandler struct{}

func (*PINGHandler) Handler(arrayOfBulkStrings []string) (string, error) {
	if !IsPINGRequest(arrayOfBulkStrings) {
		return "", fmt.Errorf(THIS_REQUEST_IS_NOT_PING_REQUEST_ERRROR)
	}

	responseElements := []string{PONGHandler.PONG}

	for i := 2; i < len(arrayOfBulkStrings); i++ {
		responseElements = append(responseElements, arrayOfBulkStrings[i])
	}
	return Builder.BuildRESPResponse(responseElements), nil
}

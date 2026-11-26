package Builder

import (
	"fmt"
	"strings"
)

import "strconv"

const CRLF = "\r\n"

func BuildRESPResponse(responseElements []string) string {
	responseArray := []string{"*" + strconv.Itoa(len(responseElements))}

	for i := range responseElements {
		responseArray = append(responseArray, "$"+strconv.Itoa(len(responseElements[i])))
		responseArray = append(responseArray, responseElements[i])
	}
	fmt.Println("DEBUG 6 RESP ARRAY", responseArray)
	var responseRESPString strings.Builder
	for i := 0; i < len(responseArray); i++ {
		responseRESPString.WriteString(responseArray[i])
		responseRESPString.WriteString(CRLF)
	}
	return responseRESPString.String()
}

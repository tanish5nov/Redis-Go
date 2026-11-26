package QueryHandler

import (
	"BuildYourOwnRedis/QueryHandler/ConcreteQueryHandlers/COMMANDHandler"
	"BuildYourOwnRedis/QueryHandler/ConcreteQueryHandlers/ECHOHandler"
	"BuildYourOwnRedis/QueryHandler/ConcreteQueryHandlers/PINGHandler"
	"BuildYourOwnRedis/RESPParser"
	"fmt"
)

func QueryHandlerFacade(data string) string {
	arrayOfBulkStrings, err := RESPParser.ParseFacade(data)

	if err != nil {
		fmt.Println("ERROR IN QUERY HANDLER", err.Error())
		return err.Error()
	}
	if len(arrayOfBulkStrings) < 2 {
		return "ERROR: NOT A COMPLETE REQUEST"
	}
	if PINGHandler.IsPINGRequest(arrayOfBulkStrings) {
		pingHandler := &PINGHandler.PINGHandler{}
		responseRESPString, err := pingHandler.Handler(arrayOfBulkStrings)
		if err != nil {
			return "ERROR WHILE BUILDING PING HANDLER RESPONSE"
		}
		fmt.Println("DEBUG RESP RESPONSE: ", responseRESPString)
		return responseRESPString
	} else if COMMANDHandler.IsCOMMANDRequest(arrayOfBulkStrings) {
		commandHandler := &COMMANDHandler.COMMANDHandler{}
		return commandHandler.Handler(arrayOfBulkStrings)
	} else if ECHOHandler.IsECHORequest(arrayOfBulkStrings) {
		echoHandler := &ECHOHandler.ECHOHandler{}
		return echoHandler.Handler(arrayOfBulkStrings)
	}
	return "UPDATE: RIGHT NOW WE ONLY SUPPORT PING"
}

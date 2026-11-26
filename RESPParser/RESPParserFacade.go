package RESPParser

import (
	"BuildYourOwnRedis/RESPParser/Parser/ConcreteParsers"
	"fmt"
)

func ParseFacade(data string) ([]string, error) {
	fmt.Println("PARSING STARTED")
	arrayOfBulkStringParser := ConcreteParsers.ArrayOfBulkStringParser{}
	arrayOfStrings, err := arrayOfBulkStringParser.Parse(data)
	if err != nil {
		fmt.Println("ERROR WHILE PARSING")
		return nil, err
	}
	return arrayOfStrings, nil
}

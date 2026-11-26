package ConcreteParsers

import (
	"fmt"
	"strconv"
)

type ArrayOfBulkStringParser struct {
}

// TODO: Ideally this parsing should be recusive, but just for the sake of simplicity keeping it flat parsing only
func (*ArrayOfBulkStringParser) Parse(data string) ([]string, error) {
	fmt.Println("ARRAY OF BULK PARSING: ", data)
	if !DataIsArray(data) {
		return []string{}, fmt.Errorf(PROVIDED_DATA_IS_NOT_ARRAY_ERROR)
	}

	i := 1
	arrayOfBulkStrings := []string{}

	arrayOfBulkStringParser := &ArrayOfBulkStringParser{}

	lengthString := ""
	lengthString, i = arrayOfBulkStringParser.parseLength(&data, &i)
	fmt.Println("DEBUG 3, ", i)
	if QueryResultStringContainsError(&lengthString) {
		return []string{}, fmt.Errorf(lengthString)
	}

	arrayOfBulkStrings = append(arrayOfBulkStrings, lengthString)

	lengthOfArrayOfBulkString, errorParsingLengthOfArrayOfBulkStrings := strconv.Atoi(lengthString)

	if errorParsingLengthOfArrayOfBulkStrings != nil {
		return []string{}, fmt.Errorf(errorParsingLengthOfArrayOfBulkStrings.Error())
	}
	fmt.Println("DEBUG ARRAY OF STRINGS: ", arrayOfBulkStrings)

	for j := 1; j <= lengthOfArrayOfBulkString; j++ {
		queryResult := ""
		queryResult, i = arrayOfBulkStringParser.getOperationType(&data, &i)
		fmt.Println("DEBUG: ", queryResult, arrayOfBulkStrings)
		if QueryResultStringContainsError(&queryResult) {
			return []string{}, fmt.Errorf(queryResult)
		}
		arrayOfBulkStrings = append(arrayOfBulkStrings, queryResult)
	}
	fmt.Println("ARRAY OF BULK PARSING: ", " ERROR NOT TILL HERE", arrayOfBulkStrings)

	return arrayOfBulkStrings, nil
}

func (*ArrayOfBulkStringParser) parseLength(Data *string, I *int) (string, int) {
	queryResult := ""
	i := *I
	data := *Data
	fmt.Println("DEBUG 4", i)
	if IndexOutOfBounds(Data, I) {
		fmt.Println("DEBUG 6, ERROR HERE")
		return INDEX_OUT_OF_BOUNDS_ERROR, i
	}

	j := i
	lengthString := ""
	for {
		if j == len(data) {
			queryResult = ENTIRE_DATA_IS_LENGTH_ERROR
			return queryResult, i
		}
		if !(data[j] >= '0' && data[j] <= '9') {
			break
		} else {
			lengthString = lengthString + string(data[j])
			j++
		}
	}
	// this will always leave on the start of next element
	j += 2
	return lengthString, j
}

func (*ArrayOfBulkStringParser) getOperationType(Data *string, I *int) (string, int) {
	i := *I
	data := *Data
	if IndexOutOfBounds(Data, I) {
		fmt.Println("DEBUG 5, ERROR HERE")
		return INDEX_OUT_OF_BOUNDS_ERROR, i
	}

	arrayOfBulkStringParser := &ArrayOfBulkStringParser{}

	if data[i] == '$' {
		i++
		lengthString, i := arrayOfBulkStringParser.parseLength(Data, &i)

		if QueryResultStringContainsError(&lengthString) {
			return lengthString, i
		}

		lengthOfString, lengthStringToIntegerError := strconv.Atoi(lengthString)
		fmt.Println("DEBUG 2 ", i, lengthString, lengthOfString)
		if lengthStringToIntegerError != nil {
			return lengthStringToIntegerError.Error(), i
		}
		element := ""
		for j := 1; j <= lengthOfString; j++ {
			if IndexOutOfBounds(Data, &i) {
				return INDEX_OUT_OF_BOUNDS_ERROR, i
			}
			element = element + string(data[i])
			i++
		}
		i += 2
		return element, i
	} else {
		return NOT_BULK_STRING_IDENTIFIER_ERROR, i
	}
}

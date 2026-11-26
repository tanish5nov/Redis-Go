package ConcreteParsers

import "strings"

const ERROR = "ERROR: "
const INDEX_OUT_OF_BOUNDS_ERROR = ERROR + "Index out of Bounds"
const NOT_BULK_STRING_IDENTIFIER_ERROR = ERROR + "No Bulk String identifier"
const PROVIDED_DATA_IS_NOT_ARRAY_ERROR = ERROR + "Provided data is not array"
const ENTIRE_DATA_IS_LENGTH_ERROR = ERROR + "Data is not correct, till the last of string there was just the length of data and nothing else"

func IndexOutOfBounds(Data *string, I *int) bool {
	return *I >= len(*Data) || *I < 0
}

func QueryResultStringContainsError(Data *string) bool {
	return strings.Contains(*Data, ERROR)
}

func DataIsArray(data string) bool {
	return len(data) > 0 && (data[0]) == '*'
}

package COMMANDHandler

const COMMAND = "COMMAND"

func IsCOMMANDRequest(arrayOfBulkString []string) bool {
	return (len(arrayOfBulkString) > 1) && arrayOfBulkString[1] == COMMAND
}

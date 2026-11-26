package ConcreteQueryHandlers

type QueryHandler interface {
	Handler(arrayOfBulkStrings []string) string
}

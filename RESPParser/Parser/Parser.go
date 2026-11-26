package Parser

type Parser interface {
	Parse(data string) ([]string, error)
}

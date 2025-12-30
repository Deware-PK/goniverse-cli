package processor

type Converter interface {
	Process(path string) error
}
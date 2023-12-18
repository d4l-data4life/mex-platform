package codings

type Codingset interface {
	Info() (map[string]string, error)
	Count() int
	GetMainHeadings() ([]string, error)
	ResolveMainHeadings(descriptorIDs []string, language string) ([]string, error)
	ResolveTreeNumbers(descriptorIDs []string) ([]string, error)
	Close()
}

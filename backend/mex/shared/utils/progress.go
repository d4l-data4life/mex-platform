package utils

type Progressor interface {
	Done()
	Progress(step string, details string)
}

type NopProgressor struct{}

func (p *NopProgressor) Progress(step string, details string) {}
func (p *NopProgressor) Done()                                {}

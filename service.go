package gengo

type Option func(*Service)

type Service struct {
	Inputs  []Input
	Outputs []Output

	Config *Config

	Routes []Route

	ToBeProcessed []Processable
	typesByName   map[string][]string

	tempMainPackageName string

	ReturnTypes []Field
}

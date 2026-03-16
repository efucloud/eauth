package options

type ServerRunOptions struct {
	Config string
}

func NewServerRunOptions() *ServerRunOptions {
	s := ServerRunOptions{}
	return &s
}

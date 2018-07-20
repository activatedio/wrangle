package plugin

type Context interface {
	Next() error
}

type Plugin interface {
	GetConfig() interface{}
	Filter(c Context) error
}

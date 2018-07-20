package config

import (
	"io"
	"reflect"
	"testing"
)

type registry struct {
	items map[string]WithConfig
}

func (self *registry) Get(name string) (WithConfig, bool) {

	r, ok := self.items[name]
	return r, ok
}

func TestConfig(t *testing.T, r io.Reader, plugins map[string]WithConfig, expected interface{}) {

	pr := &registry{
		items: plugins,
	}

	u := DefaultParser{
		pr,
	}

	got, err := u.Parse(r)

	check(err)

	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("Expected %+v, got %+v", expected, got)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

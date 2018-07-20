package plugin

import (
	"testing"
)

var _ Registry = (*DefaultRegistry)(nil)

func TestDefaultRegistry_SetGet(t *testing.T) {

	var u DefaultRegistry = make(map[string]Plugin)
	name := "name1"
	plugin := &StubPlugin{}

	u.Set(name, plugin)

	got, ok := u.Get(name)
	if !ok {
		t.Fatalf("Get is not true for %s", name)
	}
	if got != plugin {
		t.Fatalf("Wanted %v, got %v", plugin, got)

	}

	invalid := "invalid"

	_, ok = u.Get(invalid)

	if ok {
		t.Fatal("Get is true for %s", invalid)
	}
}

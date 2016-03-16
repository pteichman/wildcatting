package game

import (
	"reflect"
	"testing"
)

func TestNewEntity(t *testing.T) {
	var es entities
	e1 := es.NewEntity()
	e2 := es.NewEntity()
	e3 := es.NewEntity()
	result := []entity{e1, e2, e3}

	expected := []entity{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("NewEntity() -> %v, want %v", result, expected)
	}
}

func TestName(t *testing.T) {
	var es entities
	foo := es.NewEntity()
	bar := es.NewEntity()

	assertEqual := func(expected, actual string) {
		if expected != actual {
			t.Fatalf("%v, want %v", actual, expected)
		}
	}

	var m nameManager
	assertEqual("", m.Name(foo))
	assertEqual("", m.Name(bar))

	m.SetName(foo, "foo")
	m.SetName(bar, "bar")
	assertEqual("foo", m.Name(foo))
	assertEqual("bar", m.Name(bar))

	m.DelName(foo)
	m.DelName(bar)
	assertEqual("", m.Name(foo))
	assertEqual("", m.Name(bar))
}

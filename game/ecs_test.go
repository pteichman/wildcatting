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

func TestNameManager(t *testing.T) {
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

func TestPlayerManager(t *testing.T) {
	var es entities
	foo := es.NewEntity()
	bar := es.NewEntity()

	assertEqual := func(expected, actual interface{}) {
		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("%v, want %v", actual, expected)
		}
	}

	var m playerManager
	assertEqual(false, m.IsPlayer(foo))
	assertEqual(false, m.IsPlayer(bar))
	assertEqual(None, m.PlayerOne())

	m.AddPlayer(foo)
	assertEqual(true, m.IsPlayer(foo))
	assertEqual(false, m.IsPlayer(bar))
	assertEqual([]entity{foo}, m.Players())
	assertEqual(foo, m.PlayerOne())

	m.AddPlayer(bar)
	assertEqual(true, m.IsPlayer(foo))
	assertEqual(true, m.IsPlayer(bar))
	assertEqual([]entity{foo, bar}, m.Players())
}

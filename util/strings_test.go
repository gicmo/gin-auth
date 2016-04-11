package util

import "testing"

func TestToSnakeCase(t *testing.T) {
	var str string
	str = ToSnakeCase("Foo")
	if str != "foo" {
		t.Error("String expected to be 'foo'")
	}
	str = ToSnakeCase("FooBar")
	if str != "foo_bar" {
		t.Error("String expected to be 'foo_bar'")
	}
	str = ToSnakeCase("Bond007")
	if str != "bond007" {
		t.Error("String expected to be 'bond007'")
	}
	str = ToSnakeCase("MyUUID")
	if str != "my_uuid" {
		t.Error("String expected to be 'my_uuid'")
	}
	str = ToSnakeCase("i_hate_camels")
	if str != "i_hate_camels" {
		t.Error("String expected to be 'i_hate_camels'")
	}
}

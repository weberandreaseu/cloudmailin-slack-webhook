package main

import "testing"

func TestParseFromRegularName(t *testing.T) {
	input := "Hans Josef <hans@gmail.com>"
	name, email := parseFrom(input)
	if name != "Hans Josef" || email != "hans@gmail.com" {
		t.Errorf("parseFrom(%q) = %s, %s", input, name, email)
	}
}
func TestParseFromMultipleNames(t *testing.T) {
	input := "Ingrid-Bärbel Stinkelbröck <ingrid@gmail.com>"
	name, email := parseFrom(input)
	if name != "Ingrid-Bärbel Stinkelbröck" || email != "ingrid@gmail.com" {
		t.Errorf("parseFrom(%q) = %s, %s", input, name, email)
	}
}

func TestParseFromMailOnly(t *testing.T) {
	input := "<ingrid@gmail.com>"
	name, email := parseFrom(input)
	if name != "" || email != "ingrid@gmail.com" {
		t.Errorf("parseFrom(%q) = %s, %s", input, name, email)
	}
}

func TestParseFromInvalidInput(t *testing.T) {
	input := "invalid"
	name, email := parseFrom(input)
	if name != "" || email != "" {
		t.Errorf("parseFrom(%q) = %s, %s", input, name, email)
	}
}

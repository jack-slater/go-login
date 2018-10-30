package helpers_test

import (
	"testing"
	"github.com/jack-slater/go-login/app/helpers"
)

func TestShouldReturnDifferentStringToPassword(t *testing.T) {

	p := "password"
	if result := helpers.IncryptPassword(p); p == result {
		t.Errorf("The incrypted password: %s is the same as the submitted password: %s", result, p)
	}

}

func TestShouldIncryptEachPasswordDifferently(t *testing.T) {

	firstP := "password"
	secondP := "secondPassword"

	firstI := helpers.IncryptPassword(firstP)
	secondI := helpers.IncryptPassword(secondP)

	if firstI == secondI {
		t.Errorf("first password: %v incrypted: %v should not match second password: %v incrypted: %v", firstP, firstI, secondP, secondI)
	}

}

func TestShouldReturn64CharacterHash(t *testing.T) {

	p := "password"
	if result := helpers.IncryptPassword(p); len(result) != 64 {
		t.Errorf("The incrypted password is %d characters not the expected 64", len(result))
	}

}
package helpers_test

import (
	"testing"
	"github.com/jack-slater/go-login/app/helpers"
)

func TestVerifyEmailFormat(t *testing.T) {
	type args struct {
		e string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"correct formatted email", args{"j@j.com"}, true},
		{"email missing @ symbol", args{"jj.com"}, false},
		{"email missing . symbol", args{"jj@com"}, false},
		{"@ symbol in incorrect position", args{"@jj.com"}, false},
		{". symbol in incorrect position", args{"jj@com."}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := helpers.VerifyEmailFormat(tt.args.e); got != tt.want {
				t.Errorf("VerifyEmailFormat(%v) = %v, want %v", tt.args.e, got, tt.want)
			}
		})
	}
}

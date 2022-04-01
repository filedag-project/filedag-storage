package restapi

import (
	"fmt"
	"testing"
)

func TestNewConsoleCredentials(t *testing.T) {
	got, err := NewConsoleCredentials("test", "test", "us-east-1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(got)
	tokens, err := got.Get()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tokens)
}

package controller

import (
	"fmt"
	"testing"
)
var guid = "1234-5678-9100-5566"
func TestGenerateAccessToken(t *testing.T) {
	token,e:=GenerateAccessToken(guid)
	if e!= nil {
		t.Error(e)
	}
	fmt.Println(token)
	fmt.Println(token == "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDQ2ODA4MzcsInV1aWQiOiIxMjM0LTU2NzgtOTEwMC01NTY2In0.NlU-NCp33IIQ2KkbQbq700yqGtHMRJIP-qIX4HA9IovDRwrWNIGft6wxTWRz6YLwOWHtLIcz2VU1dHaWoVnwWQ")
}
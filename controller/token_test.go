package controller

import (
	"fmt"
	"testing"
)

func TestEncodeBase64(t *testing.T) {
	token := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIxNjA1MDEwOTE3IiwidXVpZCI6ImNhOGZkYWU2NDhmYzhlYzEyOGU5ZTQwNGIwMjZhYzg1In0._0l_sLx4fyJN429Go4WA-ZYL4OzXgu4GK-hM-97M3JUlqxIwHc-Cgg00_Gz5B3myUKT_AgBN6aqbcff6XhAF8Q"
	str := EncodeBase64(token)
	fmt.Println(str)
	fmt.Println(DecodeBase64(str) == token)
}

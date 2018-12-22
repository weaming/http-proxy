package main

import "testing"

func TestAES(t *testing.T) {
	s := "hello world"
	key := "weaming"
	encrypted := AES([]byte(s), key)

	ss := string(AES(encrypted, key))
	println(ss)
	if ss != s {
		t.Fail()
	}
}

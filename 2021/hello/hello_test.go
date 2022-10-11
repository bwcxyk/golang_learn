/*
@Author : yaokun
@Time : 2020/7/15 18:08
*/

package main

import "testing"

func TestHello(t *testing.T) {
	got := Hello()
	want := "Hello, world"

	if got != want {
		t.Errorf("got '%q' want '%q'", got, want)
	}
}

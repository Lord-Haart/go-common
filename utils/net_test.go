package utils

import "testing"

func TestSplitHostPort(t *testing.T) {
	testcases := []struct {
		p1 string
		r1 string
		r2 uint64
	}{
		{"localhost:80", "localhost", 80},
		{"localhost:8080", "localhost", 8080},
		{"127.0.0.1", "127.0.0.1", 0},
		{"127.0.0.1:8080", "127.0.0.1", 8080},
		{"[::1]", "[::1]", 0},
		{"[::1]:8080", "[::1]", 8080},
	}

	for _, testcase := range testcases {
		r1, r2 := SplitHostPort(testcase.p1)
		if r1 != testcase.r1 || r2 != testcase.r2 {
			t.Errorf("SplitHostAndPort(%#v) => %#v, %#v, wants %#v, %#v", testcase.p1, r1, r2, testcase.r1, testcase.r2)
		}
	}
}

func TestGetLocalIpAddress(t *testing.T) {
	t.Logf("%v", GetLocalIpAddress())
}

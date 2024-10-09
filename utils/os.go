package utils

import (
	"os"
	"strings"
)

// ExpandEnv 执行环境变量插值。
func ExpandEnv(s string) string {
	return os.Expand(s, func(v string) string {
		if strings.ContainsRune(v, ':') {
			vv := strings.SplitN(v, ":", 2)
			if vs, ok := os.LookupEnv(strings.TrimSpace(strings.ToUpper(vv[0]))); ok {
				return vs
			} else {
				return vv[1]
			}
		} else {
			return os.Getenv(strings.ToUpper(v))
		}
	})
}

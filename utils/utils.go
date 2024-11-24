package utils

import (
	"flag"
	"fmt"
	"os"
)

func UsageAndExitt(msg string) {
	if msg != "" {
		fmt.Println(msg)
	}
	flag.Usage()
	fmt.Println("")
	os.Exit(1)
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

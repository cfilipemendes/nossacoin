package main

import "os"

func isDev() bool {
	return os.Getenv("ENV") == "dev"
}

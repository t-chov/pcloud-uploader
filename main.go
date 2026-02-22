package main

import "os"

const (
	APP_NAME            = "pcloud-uploader"
	ENV_PCLOUD_USERNAME = "PCLOUD_USERNAME"
	ENV_PCLOUD_PASSWORD = "PCLOUD_PASSWORD"
)

// VERSION is the application version, injected at build time via ldflags.
var VERSION = "dev"

func main() {
	os.Exit(run())
}

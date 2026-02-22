package main

import "os"

const (
	APP_NAME            = "pcloud-uploader"
	VERSION             = "0.0.3"
	ENV_PCLOUD_USERNAME = "PCLOUD_USERNAME"
	ENV_PCLOUD_PASSWORD = "PCLOUD_PASSWORD"
)

func main() {
	os.Exit(run())
}

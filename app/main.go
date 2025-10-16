package main

import "usp-management-device-api/application"

var (
	buildTime string
	version   string
)

func main() {
	application.StartApplication(buildTime, version)
}

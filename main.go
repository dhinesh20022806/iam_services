package main

import (
	"iam_service/app/cmd"
	"os"
)



func main() {

	args := os.Args[1:]

	if len(args) > 0 && args[0] == "migrate" {
		os.Exit(cmd.RunMigrate())
	} else {
		os.Exit(cmd.RunServer())
	}
}
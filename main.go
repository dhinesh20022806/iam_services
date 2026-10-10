package main

import (
	"iam_service/app/cmd"
	"os"
)



func main() {


		os.Exit(cmd.RunServer())
	
}
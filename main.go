package main

import (
	"os"

	"github.com/dhinesh20022806/iam_service/app/cmd"
)



func main() {


		os.Exit(cmd.RunServer())
	
}
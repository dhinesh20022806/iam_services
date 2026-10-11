package cmd

import (
	"log/slog"

	"github.com/dhinesh20022806/iam_service/app/services"
)



func RunServer() int {

	server := services.NewGrpcServer(":50050")

	err := server.Run()

	if err != nil {
		slog.Error("Error creating server", "err", err)
		return 0
	}



	return 1
}
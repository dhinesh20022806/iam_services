package cmd

import (
	"iam_service/app/pkg/dbx"
	"log/slog"
)



func RunMigrate() int {
	err := dbx.Migrate("/migrations")
	if err != nil {
		slog.Error("Migrate error", "error", err)
		return 1
	}

	return 0
}
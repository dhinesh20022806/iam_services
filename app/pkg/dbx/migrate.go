package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"iam_service/app/pkg/env"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
)

var conn *sql.DB


func init() {

	var err error
	conn, err := sql.Open("postgres", env.Config.Database.URL)

	if err != nil {
		fmt.Printf("failed to open connection to the database %v", err)
	}

	conn.SetMaxIdleConns(env.Config.Database.MaxIdleConn)
	conn.SetMaxOpenConns(env.Config.Database.MaxOpenConn)
}

func Migrate(path string) error {


	dir, err := os.Open(path)

	if err != nil {
		return fmt.Errorf("failed to open %v", err)
	}

	files, err := dir.Readdir(0)

	if err != nil {
		fmt.Errorf("Failed to open read %s and %v", path, err)
	}

	versions := make([]int, len(files))
	versionFiles := make(map[int]string, len(files))

	for i, file := range files {

		fileName := file.Name()
		
		parts := strings.Split(fileName, "_")

		
		versions[i], err = strconv.Atoi(parts[0])
        versionFiles[versions[i]] = fileName

		if err != nil {
	     return	fmt.Errorf("Failed to open read %v", err)
		}
	}

	sort.Ints(versions)

	lastVersion, err := getLastMigrate()

	if err != nil {
		return fmt.Errorf("problem with getting lastversion number %v", err)
	}

	slog.Info("Current values", "#lastversion", lastVersion)

	totalMigrationExecuted := 0

	for _, version := range versions {
		if version > lastVersion {
			fileName := versionFiles[version]
			slog.Info("Running Version:", "version", version, "filename", fileName)
			err := runMigration(version, path, fileName)

			if err != nil {
				return fmt.Errorf("Running Migrate %v", err)
			}
			totalMigrationExecuted++
		}

	}


	if totalMigrationExecuted > 0 {
		slog.Info("migration have been applied", "totoal Migration", totalMigrationExecuted)
	} else {
		slog.Info("No Migration have been applied")
	}


	return nil

}

type Trx struct {
	tx *sql.Tx
	ctx context.Context
}


func BeginTx(ctx context.Context) (*Trx, error) {

	tx, err := conn.BeginTx(ctx, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to start the transaction %v", err)
	}

	return &Trx{tx: tx, ctx: ctx}, nil
}

func runMigration(version int, path string, filename string) error {

	content, err := os.ReadFile(path + "/" + filename)

	if err != nil {

	}

	ctx := context.Background()

	trx, err := BeginTx(ctx)

	if err != nil {
	  	return err
	}

	_, err = trx.tx.Exec(string(content))

	if err != nil {
		return err
	}

	_, err = trx.tx.Exec("INSERT INTO migrations_history (version, filename) VALUES ($1, $2)", version, filename)

	if err != nil {

		return err 
	}

	return trx.tx.Commit()
}

func getLastMigrate() (int, error) {

	_, err := conn.Exec(`CREATE TABLE IF NOT EXISTS migrations_history (
	     versions BIGINT PRIMARY KEY,
		 filename VARCHAR(100) null,
		 date TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`)

	if err != nil {
		return 0, err
	}

	var lastVersion sql.NullInt64

	row := conn.QueryRow("SELECT MAX(version) FROM migrations_history LIMIT 1")

	err = row.Scan(&lastVersion)

	if err != nil {
		return 0, err
	}
	return int(lastVersion.Int64), nil
	
}
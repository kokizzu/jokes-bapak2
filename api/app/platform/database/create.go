package database

import (
	"context"
	"log"

	"github.com/aldy505/bob"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Setup the table connection, create table if not exists
func Setup(db *pgxpool.Pool, ctx *context.Context) error {
	conn, err := db.Acquire(*ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(*ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(*ctx)

	// administrators table
	var tableAuthExists bool
	err = db.QueryRow(*ctx, `SELECT EXISTS (
		SELECT FROM information_schema.tables 
		WHERE  table_schema = 'public'
		AND    table_name   = 'administrators'
		);`).Scan(&tableAuthExists)
	if err != nil {
		log.Fatalln("16 - failed on checking table: ", err)
		return err
	}

	if !tableAuthExists {
		sql, _, err := bob.
			CreateTable("administrators").
			AddColumn(bob.ColumnDef{Name: "id", Type: "SERIAL", Extras: []string{"PRIMARY KEY"}}).
			StringColumn("key", "NOT NULL", "UNIQUE").
			TextColumn("token").
			StringColumn("last_used").
			ToSql()
		if err != nil {
			log.Fatalln("17 - failed on table creation: ", err)
			return err
		}

		_, err = tx.Exec(*ctx, sql)
		if err != nil {
			log.Fatalln("18 - failed on table creation: ", err)
			return err
		}
	}

	// Jokesbapak2 table

	// Check if table exists
	var tableJokesExists bool
	err = db.QueryRow(*ctx, `SELECT EXISTS (
		SELECT FROM information_schema.tables 
		WHERE  table_schema = 'public'
		AND    table_name   = 'jokesbapak2'
		);`).Scan(&tableJokesExists)
	if err != nil {
		log.Fatalln("10 - failed on checking table: ", err)
		return err
	}

	if !tableJokesExists {
		sql, _, err := bob.
			CreateTable("jokesbapak2").
			AddColumn(bob.ColumnDef{Name: "id", Type: "SERIAL", Extras: []string{"PRIMARY KEY"}}).
			TextColumn("link", "UNIQUE").
			AddColumn(bob.ColumnDef{Name: "creator", Type: "INT", Extras: []string{"NOT NULL", "REFERENCES \"administrators\" (\"id\")"}}).
			ToSql()
		if err != nil {
			log.Fatalln("11 - failed on table creation: ", err)
			return err
		}

		_, err = tx.Exec(*ctx, sql)
		if err != nil {
			log.Fatalln("12 - failed on table creation: ", err)
			return err
		}
	}

	// Submission table

	//Check if table exists
	var tableSubmissionExists bool
	err = db.QueryRow(*ctx, `SELECT EXISTS (
		SELECT FROM information_schema.tables 
		WHERE  table_schema = 'public'
		AND    table_name   = 'submission'
		);`).Scan(&tableJokesExists)
	if err != nil {
		log.Fatalln("13 - failed on checking table: ", err)
		return err
	}

	if !tableSubmissionExists {
		sql, _, err := bob.
			CreateTable("submission").
			AddColumn(bob.ColumnDef{Name: "id", Type: "SERIAL", Extras: []string{"PRIMARY KEY"}}).
			TextColumn("link", "UNIQUE", "NOT NULL").
			StringColumn("created_at").
			StringColumn("author", "NOT NULL").
			AddColumn(bob.ColumnDef{Name: "status", Type: "SMALLINT", Extras: []string{"DEFAULT 0"}}).
			ToSql()
		if err != nil {
			log.Fatalln("14 - failed on table creation: ", err)
		}

		_, err = tx.Query(*ctx, sql)
		if err != nil {
			log.Fatalln("15 - failed on table creation: ", err)
		}
	}

	err = tx.Commit(*ctx)
	if err != nil {
		return err
	}

	return nil
}
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func query(db *sql.DB, quer string) error {
	queries := strings.Split(quer, ";")
	for _, q := range queries[:len(queries)-1] {
		//fmt.Println(q)
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}

	//fmt.Println(queries[len(queries)-1])
	rows, err := db.Query(queries[len(queries)-1])
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return errors.New("No rows!")
	}

	var names []string
	names, err = rows.Columns()
	if err != nil {
		return err
	}
	values := make([]interface{}, len(names))
	refs := make([]interface{}, len(names))
	for i := range values {
		refs[i] = &(values[i])
	}
	if err = rows.Scan(refs...); err != nil {
		return err
	}

	for i := range names {
		fmt.Printf("%s: %T\n", names[i], values[i])
	}

	return nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintln(os.Stderr, "usage:", os.Args[0], "<driver> <connstr> <query>")
		os.Exit(2)
	}
	db, err := sql.Open(os.Args[1], os.Args[2])
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	if err = query(db, os.Args[3]); err != nil {
		log.Println(err)
		// As os.Exit shortcircuits defers, we have to close explicitely
		db.Close()
		os.Exit(1)
	}
}

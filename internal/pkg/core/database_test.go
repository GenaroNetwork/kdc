package core

import (
	"testing"

	"database/sql"
	"log"
	"path/filepath"
	"time"
	"math/big"
	"fmt"
)

func TestDB(t *testing.T) {
	dbFileFullPath := filepath.Join(dbHome, "own.db")
	log.Print(dbFileFullPath)
	db, err := sql.Open("sqlite3", dbFileFullPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table if not exists foo (id integer not null primary key, name text);
	delete from foo;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

var q = make(chan func(), 1)

func useChan(funcs chan func()) {
	for true {
		funcs <- func() {
			log.Print(1)
		}
		time.Sleep(time.Second * 2)
	}
}

func TestChan(t *testing.T) {
	funcs := make(chan func(), 1)
	go useChan(funcs)
	for true {
		func1 := <- funcs
		func1()
	}
}

func TestBigint(t *testing.T) {
	n := new(big.Int)
	n, ok := n.SetString("ffff", 16)
	if !ok {
		fmt.Println("SetString: error")
		return
	}
	fmt.Printf("0x%x", n)
}

func TestGetOrCreateFileDB(t *testing.T) {
	getOrCreateFileDB("dsadsa")
}
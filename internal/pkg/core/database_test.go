package core

import (
	"testing"

	"database/sql"
	"log"
	"path/filepath"
	"math/big"
	"fmt"
	"time"
	"sync"
	"strconv"
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
	i := 1
	for true {
		funcs <- func() {
			i++
			log.Print(i)
			time.Sleep(time.Millisecond * 2)
		}
		fmt.Println("func given")
	}
}

func TestChan(t *testing.T) {
	funcs := make(chan func(), 10)
	go useChan(funcs)
	for true {
		func1 := <- funcs
		func1()
	}
}

func useChan2(funcs chan func()) int{
	var wg sync.WaitGroup
	wg.Add(1)
	i := 0
	funcs <- func() {
		wg.Done()
		i ++
		time.Sleep(time.Millisecond * 2)
	}
	wg.Wait()
	return i
}

func TestWaitChan(t *testing.T) {
	funcs := make(chan func(), 10)
	go func() {
		for true {
			func1 := <- funcs
			func1()
		}
	}()

	b := useChan2(funcs)
	fmt.Println(b)
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
	mt := MortgageTableT{
		"userAsss": *big.NewInt(1),
		"r":   *big.NewInt(2),
		"gri": *big.NewInt(3),
		"adg": *big.NewInt(4),
	}
	at := AllowTableT{
		"userAsss": 1,
	}
	for i := 1; i <= 100; i++ {
		go initNewFile("0xbbbb" + strconv.Itoa(i), "0xowner","{}", &at, &mt)
	}
	time.Sleep(time.Second * 5)
}

func TestIsOwner(t *testing.T) {
	b, _ := isOwner("0xbbbb1", "0xowner")
	fmt.Println(b)
}

func TestSetTerminate(t *testing.T) {
	err := setFileTerminate("0xbbbb10")
	if err != nil {
		t.Fatal(err)
	}
}
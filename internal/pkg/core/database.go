package core

import (
	"os/user"
	"github.com/op/go-logging"
	"path/filepath"
	"os"
	"database/sql"
	"fmt"
)
var dbLog = logging.MustGetLogger("database")
var dbHome string

func init() {
	usr, err := user.Current()
	if err != nil {
		dbLog.Fatal(err)
	}
	homeDir := usr.HomeDir
	dbHome = filepath.Join(homeDir, ".kdc")
	os.MkdirAll(dbHome, os.ModePerm)
	dbLog.Debug("database home dir: %s", dbHome)
	// ensure owner table exist
	{
		dbFileFullPath := filepath.Join(dbHome, "own.db")
		dbLog.Debug(dbFileFullPath)
		db, err := sql.Open("sqlite3", dbFileFullPath)
		if err != nil {
			dbLog.Fatal(err)
			return
		}
		defer db.Close()

		fileIndexSql := `create table if not exists fileIndex 
		                 (fileId text not null primary key, 
		                 owner text not null, 
		                 isopen INTEGER DEFAULT 1, 
		                 originjson text, 
		                 state text );`

		privilegeSql := `create table if not exists privilege
						 (fileId text not null, 
						 user text not null, 
					  	 privilege INTEGER not null);`

		_, err = db.Exec(fileIndexSql)
		if err != nil {
			dbLog.Fatal("%q: %s\n", err, fileIndexSql)
			return
		}
		_, err = db.Exec(privilegeSql)
		if err != nil {
			dbLog.Fatal("%q: %s\n", err, fileIndexSql)
			return
		}
	}
}


func getOrCreateFileDB(fileId string) {
	dbFileFullPath := filepath.Join(dbHome, fileId + ".db")
	dbLog.Debug("sqlite3 DB path: %s", dbFileFullPath)
	db, err := sql.Open("sqlite3", dbFileFullPath)
	if err != nil {
		dbLog.Fatal(err)
		return
	}
	defer db.Close()

	sqlStmt := fmt.Sprintf(`create table if not exists FILE_%s (fileId text not null primary key, user text not null);`, fileId)
	_, err = db.Exec(sqlStmt)
	if err != nil {
		dbLog.Fatal("%q: %s\n", err, sqlStmt)
		return
	}
}
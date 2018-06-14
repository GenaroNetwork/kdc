package core

import (
	"os/user"
	"github.com/op/go-logging"
	"path/filepath"
	"os"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"sync"
	"time"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"errors"
)

var dbLog = logging.MustGetLogger("database")
var dbHome string
var dbMutex = &sync.Mutex{}
var dbConn *sql.DB

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
		dbConn, err = sql.Open("sqlite3", dbFileFullPath)
		if err != nil {
			dbLog.Fatal(err)
			return
		}

		fileIndexSql := `create table if not exists fileIndex 
		                 (fileId text not null primary key, 
		                 owner text not null, 
		                 isopen INTEGER DEFAULT 1, 
		                 originjson text, 
		                 state text,
		                 createTime int not null);`

		privilegeSql := `create table if not exists privilege
						 (fileId text not null, 
						 user text not null, 
					  	 privilege INTEGER not null,
					  	 createTime int not null);`

		_, err = dbConn.Exec(fileIndexSql)
		if err != nil {
			dbLog.Fatal("%q: %s\n", err, fileIndexSql)
			return
		}
		_, err = dbConn.Exec(privilegeSql)
		if err != nil {
			dbLog.Fatal("%q: %s\n", err, fileIndexSql)
			return
		}
	}
}

func initNewFile(fileId string, owner string, originJson string, allow *AllowTableT, mortgage *MortgageTableT) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	nowTime := time.Now().Unix()
	tx, err := dbConn.Begin()
	defer tx.Commit()
	if err != nil {
		dbLog.Fatal(err)
	}
	// insert into fileIndex TODO:
	sqlIndex := fmt.Sprintf(`insert into fileIndex (fileId, owner, originjson, createTime) values ('%s', '%s', '%s', %d);`, fileId, owner, originJson, nowTime)
	_, err1 := tx.Exec(sqlIndex)
	if err1 != nil {
		dbLog.Error("%q: %s\n", err1, sqlIndex)
		return err1
	}
	// insert into privilege
	stmtP, err := tx.Prepare("insert into privilege(fileId, user, privilege, createTime) values(?, ?, ?, ?)")
	defer stmtP.Close()
	if err != nil {
		dbLog.Fatal(err)
	}
	for userA, privA := range *allow {
		stmtP.Exec(fileId, userA, privA, nowTime)
	}
	// create modify table and insert init value
	tableName := fmt.Sprintf("FILE_%s", fileId)
	sqlStmt := fmt.Sprintf(`create table %s (userId text not null, opration text, value text, createTime int not null);`, tableName)

	_, err2 := tx.Exec(sqlStmt)
	if err2 != nil {
		dbLog.Error("%q: %s\n", err2, sqlStmt)
		return err2
	}

	stmtM, err := tx.Prepare(fmt.Sprintf("insert into %s (userId, opration, value, createTime) values (?, ?, ?, ?);", tableName))
	defer stmtM.Close()
	if err != nil {
		dbLog.Fatal(err)
	}
	for userM, coins := range *mortgage {
		_, err = stmtM.Exec(userM, "init", hexutil.EncodeBig(&coins), nowTime)
		if err != nil {
			dbLog.Error("%q: %s\n", err, sqlStmt)
			return err
		}
	}
	return nil
}

func isOwner(fileId string, user string) (bool, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	stmt, err := dbConn.Prepare("select count(1) count from fileIndex where fileId = ? and owner = ?")
	if err != nil {
		dbLog.Error("isOwner sql err: %s", err)
		return false, err
	}
	defer stmt.Close()
	var count int
	err = stmt.QueryRow(fileId, user).Scan(&count)
	if err != nil {
		dbLog.Error("isOwner sql err: %s", err)
		return false, err
	}

	return count == 1, nil
}

func setFileTerminate(fileId string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	tx, err := dbConn.Begin()
	defer tx.Commit()
	if err != nil {
		dbLog.Fatal(err)
	}

	stmtM, err := tx.Prepare("update fileIndex set isopen = 0 where fileId = ?")
	defer stmtM.Close()
	if err != nil {
		dbLog.Fatal(err)
	}
	result, err := stmtM.Exec(fileId)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == int64(0) {
		return errors.New("terminate sql has no effect")
	}

	return nil
}
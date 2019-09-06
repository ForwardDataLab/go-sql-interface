package sqlinterface

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"
)

// InterfaceTest : tests the interface import
func InterfaceTest() int {
    return 5
}

// GetRows : fetches rows from DB
func (db DB) GetRows(rowAccess RowAccess) []RowStructure {
    if(db.DbType == "mysql") {
        return mysqlGetRows(db, rowAccess)
    } else if (db.DbType == "postgres") {
        return postgresGetRows(db, rowAccess)
    } else {
        return nil
    }
}

func mysqlGetRows(db DB, rowAccess RowAccess) []RowStructure {
    currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "@/" + rowAccess.DatabaseName)
    queryString := "SELECT * FROM " +
        rowAccess.Table +
        " WHERE " + rowAccess.Column + " = ?"
    statement, _ := currentDatabase.Prepare(queryString)
    fetchedArr := make([]RowStructure, len(rowAccess.Indices))
    for index, rowIndex := range rowAccess.Indices {
        statement.QueryRow(rowIndex).Scan(&(fetchedArr[index]).USER_NAME,
            &(fetchedArr[index]).INDEX_COL)
    }
    return fetchedArr
}

func postgresGetRows(db DB, rowAccess RowAccess) []RowStructure {
    currentDatabase, _ := sql.Open(db.DbType, "user=" + db.Username +
        " password=" + db.Password +
        " dbname=" + rowAccess.DatabaseName +
        "sslmode=disable")
    queryString := "SELECT * FROM " +
        rowAccess.Table +
        " WHERE " + rowAccess.Column + " = 1"
    fetchedArr := make([]RowStructure, len(rowAccess.Indices))
    for index, rowIndex := range rowAccess.Indices {
        row := currentDatabase.QueryRow(queryString, rowIndex)
        err := row.Scan(&(fetchedArr[index]).USER_NAME,
            &(fetchedArr[index]).INDEX_COL)
        if err != nil {
            print(err)
        }
    }
    return fetchedArr
}

package sqlinterface

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

// DB : the database struct to store information about the connection
type DB struct {
    DbType string
    Username string
    Password string
}

// RowAccess : struct to request rows from the database
type RowAccess struct {
    DatabaseName string
    Table string
    Column string
    Indices []int
}

// InterfaceTest : tests the interface import
func InterfaceTest() int {
    return 5
}

// GetRows : fetches rows from DB
func (db DB) GetRows(rowAccess RowAccess) []interface{} {
    currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password)
    queryString := "SELECT * FROM " +
        rowAccess.DatabaseName + "." + rowAccess.Table +
        " WHERE " + rowAccess.Column + " = ?"
    statement, _ := currentDatabase.Prepare(queryString)
    fetchedArr := make([]interface{}, len(rowAccess.Indices))
    for _, index := range rowAccess.Indices {
        statement.QueryRow(index).Scan(&fetchedArr[index])
    }
    return fetchedArr
}

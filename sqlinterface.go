package sqlinterface

import (
    "database/sql"
    "fmt"
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

type rowStructure struct {
    USER_NAME string
    INDEX_COL int
}

// InterfaceTest : tests the interface import
func InterfaceTest() int {
    return 5
}

// GetRows : fetches rows from DB
func (db DB) GetRows(rowAccess RowAccess) []rowStructure {
    currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "/" + rowAccess.DatabaseName)
    queryString := "SELECT * FROM " +
        rowAccess.Table +
        " WHERE " + rowAccess.Column + " = ?"
    statement, _ := currentDatabase.Prepare(queryString)
    fetchedArr := make([]rowStructure, len(rowAccess.Indices))
    for index, rowIndex := range rowAccess.Indices {
        statement.QueryRow(rowIndex).Scan(&(fetchedArr[index]))
    }
    return fetchedArr
}

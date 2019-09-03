package sqlinterface

import (
    "fmt"
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
func (rowStruct RowAccess) GetRows() {
    fmt.Println(rowStruct.DatabaseName)
    fmt.Println(rowStruct.Table)
    fmt.Println(rowStruct.Column)
    fmt.Println(rowStruct.Indices)
}

package sqlinterface

import (
        "database/sql"
)

// DB : the database struct to store information about the connection
type DB struct {
    Host string
    Port string
    DbType string
    Username string
    Password string
    DatabaseName string
    Table string
    ClusterMap map[int]int
    ClusterSize int
    NumClusters int
    fresh bool
    newConfiguration []int
}

// RowAccess : struct to request rows from the database
type RowAccess struct {
    Column string
    Indices []int
}

// Cell : struct to store an information about a cell in a row
type Cell struct {
    Type string
    Value string
}

// TableMetadata : struct to store the result of describe table
type TableMetadata struct {
    Field sql.NullString
    Type sql.NullString
    Null sql.NullString
    Key sql.NullString
    Default sql.NullString
    Extra sql.NullString
}

func max(a int, b int) int {
    if a > b {
        return a
    }
    return b
}

func min(a int, b int) int {
    if a < b {
        return a
    }
    return b
}

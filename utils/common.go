package sqlinterface

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

// RowStructure : struct to represent the row structure of the database
type RowStructure struct {
    USER_NAME string
    INDEX_COL int
}


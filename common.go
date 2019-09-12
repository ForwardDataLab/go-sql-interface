package sqlinterface

// DB : the database struct to store information about the connection
type DB struct {
    DbType string
    Username string
    Password string
    DatabaseName string
    Table string
}

// RowAccess : struct to request rows from the database
type RowAccess struct {
    Column string
    Indices []int
}

// Cell : struct to store an information about a cell in a row
type Cell struct {
    Column string
    Type string
    Value string
}

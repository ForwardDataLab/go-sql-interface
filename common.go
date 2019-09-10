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

// RowStructure : struct to represent the row structure of the database
type RowStructure struct {
    ID int
    NAME string
    Age int
    Department string
    GPA float32
}

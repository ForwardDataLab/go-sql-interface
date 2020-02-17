package sqlinterface

import (
    "database/sql"
    "strconv"
    "strings"
    "fmt"
)

// OpenSQLDB: connect to mysql and return *sql.DB which is a connection pool
func (db DB) OpenSQLDB() *sql.DB {
    currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "@tcp(" + db.Host + ":" + db.Port + ")/" + db.DatabaseName)
    return currentDatabase
}

// MysqlGetRowsBatch: query rows in mysql db
func MysqlGetRowsBatch(TableName string, currentDatabase *sql.DB, rowAccess RowAccess) [][]string {
    convertedIndices := make([]interface{}, len(rowAccess.Indices))
    for i, v := range rowAccess.Indices {
        convertedIndices[i] = v
    }
    queryString := "SELECT * FROM " +
        TableName +
        " WHERE " + rowAccess.Column + " in (?" + strings.Repeat(", ?", len(convertedIndices) - 1) + ")"
    statement, err := currentDatabase.Prepare(queryString)
    defer statement.Close()
    if (err != nil) {
        fmt.Println(err);
    }
    rows, err := statement.Query(convertedIndices...)
    defer rows.Close()
    if (err == nil) {
        columns, _ := rows.Columns()
        values := make([]sql.RawBytes, len(columns))
        defer rows.Close()
        fetchedArr := make([]interface{}, len(values))
        for i := range values {
            fetchedArr[i] = &(values[i])
        }

        returnArr := [][]string{}
        for rows.Next() {
            rows.Scan(fetchedArr...)
            currentArr := []string{}
            for _, v := range values {
                if v == nil {
                    currentArr = append(currentArr, "NULL")
                } else {
                    currentArr = append(currentArr, string(v))
                }
            }
            returnArr = append(returnArr, currentArr)
        }
        return returnArr
    }
    return nil
}

// MysqlInsertRow: insert a row into mysql db
func MysqlInsertRow(TableName string, currentDatabase *sql.DB, indexCol string, cells []Cell) int {
    // INSERT INTO table_name (col, col, col) VALUES (NULL, 'my name', 'my group')
    selectMaxQueryString := "SELECT MAX(" + indexCol + ") FROM " + TableName
    var maxIndex int
    rows, _ := currentDatabase.Query(selectMaxQueryString)
    defer rows.Close()
    for rows.Next() {
        rows.Scan(&maxIndex)
    }

    // " (?" + strings.Repeat(", ?", len(cells) - 1) + ")" +
    insertQueryString := "INSERT INTO " +
        TableName +
        " VALUES (?" + strings.Repeat(", ?", len(cells) - 1) + ")"
    insertStatement, err := currentDatabase.Prepare(insertQueryString)
    defer insertStatement.Close()
    if err != nil {
        fmt.Println(insertQueryString)
        fmt.Println(err)
    }

    // create interface and add max index
    insertCell := make([]interface{}, len(cells))

    for i, v := range cells {
        if v.Type == "ID" {
            insertCell[i] = maxIndex + 1
        } else if v.Type == "int" {
            insertCell[i], _ = strconv.ParseInt(v.Value, 10, 32)
        } else if v.Type == "string" {
            insertCell[i] = string(v.Value)
        } else if v.Type == "float" {
            insertCell[i], _ = strconv.ParseFloat(v.Value, 64)
        } else if v.Type == "bool" {
            insertCell[i], _ = strconv.ParseBool(v.Value)
        } else {
            fmt.Println(v.Type)
            return -1
        }
    }
    _, _ = insertStatement.Exec(insertCell...)
    return maxIndex + 1
}

// MysqlDeleteRow: delete a row in mydql db
func MysqlDeleteRow(TableName string, currentDatabase *sql.DB, indexCol string, index int) {
    // DELETE FROM table_name WHERE index_col = index
    deleteQueryString := "DELETE FROM " +
        TableName +
        " WHERE " + indexCol + " = ?"
    deleteStatement, _ := currentDatabase.Prepare(deleteQueryString)
    defer deleteStatement.Close()
    _, err := deleteStatement.Exec(index)
    if (err != nil) {
        fmt.Println(err);
    }
}

// CloseMySQLDB: close mysql database connection pool
func CloseMySQLDB(currentDatabase sql.DB) {
    currentDatabase.Close()
}

// GetColMap : gets the column mapping from DB
func MysqlGetColMap(TableName string, currentDatabase sql.DB) []TableMetadata {
    // colMap := make(map[string]string)
    columnQueryString := "DESCRIBE " + TableName
    var tableMetadata TableMetadata
    rows, _ := currentDatabase.Query(columnQueryString)
    defer rows.Close()
    var returnArr []TableMetadata
    for rows.Next() {
        rows.Scan(
            &(tableMetadata.Field),
            &(tableMetadata.Type),
            &(tableMetadata.Null),
            &(tableMetadata.Key),
            &(tableMetadata.Default),
            &(tableMetadata.Extra),
        )
        returnArr = append(returnArr, tableMetadata)
    }
    return returnArr
}




// GetColMap : gets the column mapping from DB
func (db DB) GetColMap() []TableMetadata {
    if(db.DbType == "mysql") {
        return mysqlGetColMap(db)
    } else if (db.DbType == "postgres") {
        return nil
    } else {
        // should panic or do proper error throwing
        return nil
    }
}

// EvaluateFormula : evaluates arbitrary sql statements
func (db DB) EvaluateFormula(formulaData string) [][]string {
    if(db.DbType == "mysql") {
        return mysqlEvaluateFormula(db, formulaData)
    } else if (db.DbType == "postgres") {
        return nil
    } else {
        // should panic or do proper error throwing
        return nil
    }
}

// GetRows : fetches rows from DB
func (db DB) GetRows(rowAccess RowAccess) [][]string {
    if(db.DbType == "mysql") {
        return mysqlGetRows(db, rowAccess)
    } else if (db.DbType == "postgres") {
        return postgresGetRows(db, rowAccess)
    } else {
        // should panic or do proper error throwing
        return nil
    }
}

// GetColumns : fetches columns from DB
func (db DB) GetColumns(columnAccess ColumnAccess) [][]string {
    if(db.DbType == "mysql") {
        return mysqlGetColumns(db, columnAccess)
    } else if (db.DbType == "postgres") {
        return postgresGetColumns(db, columnAccess)
    } else {
        // should panic or do proper error throwing
        return nil
    }
}

// InitDB : initializes the database upon initial creation of workspace
func (db *DB) InitDB() {
    // add a column called index_col
    // ALTER TABLE `myTable` ADD COLUMN `id` INT AUTO_INCREMENT UNIQUE
    if(db.DbType == "mysql") {
        mysqlInitDB(db)
    } else if (db.DbType == "postgres") {
        postgresInitDB()
    } else {
        // should panic or do proper error throwing
    }
}

// InsertRow : inserts a new row into the database
func (db DB) InsertRow(indexCol string, cells []Cell) int {
    // insert a row into db defined by rowStructure
    // INSERT INTO table_name (col, col, col) VALUES (NULL, 'my name', 'my group')
    if(db.DbType == "mysql") {
        return mysqlInsertRow(db, indexCol, cells, false)
    } else if (db.DbType == "postgres") {
        // return postgresInsertRow(db, row)
        return -1
    } else {
        // should panic or do proper error throwing
        return -1
    }
}

// InsertColumn : inserts a new column into the database
func (db DB) InsertColumn(columnName string, columnType string) int {
    // insert a column into db defined by columnStructure
    // INSERT INTO table_name (col, col, col) VALUES (NULL, 'my name', 'my group')
    if(db.DbType == "mysql") {
        return mysqlInsertColumn(db, columnName, columnType)
    } else if (db.DbType == "postgres") {
        // return postgresInsertColumn(db, column)
        return -1
    } else {
        // should panic or do proper error thcolumning
        return -1
    }
}

// DeleteRow : delets a row from the database
func (db DB) DeleteRow(indexCol string, index int) {
    // DELETE FROM table_name WHERE index_col = index
    if(db.DbType == "mysql") {
        mysqlDeleteRow(db, indexCol, index)
    } else if (db.DbType == "postgres") {
        postgresDeleteRow(db, index)
    } else {
        // should panic or do proper error throwing
    }
}

// UpdateRow : updates a row from the database
func (db DB) UpdateRow(indexCol string, cells []Cell) {
    // UPDATE table_name WHERE index_col = index
    if(db.DbType == "mysql") {
        mysqlUpdateRow(db, indexCol, cells)
    } else if (db.DbType == "postgres") {
        // update but postgres
    } else {
        // should panic or do proper error throwing
    }
}

// GetRowsBatch : fetches rows from DB in batches
func (db DB) GetRowsBatch(rowAccess RowAccess) [][]string {
    if(db.DbType == "mysql") {
        return mysqlGetRowsBatch(db, rowAccess)
    } else if (db.DbType == "postgres") {
        return postgresGetRows(db, rowAccess)
    } else {
        // should panic or do proper error throwing
        return nil
    }
}

// GetRowsSerial : fetches rows from DB in serial
func (db DB) GetRowsSerial(rowAccess RowAccess) [][]string {
    if(db.DbType == "mysql") {
        return mysqlGetRowsSerial(db, rowAccess)
    } else if (db.DbType == "postgres") {
        return postgresGetRows(db, rowAccess)
    } else {
        // should panic or do proper error throwing
        return nil
    }
}

// GetRowsCluster : fetches rows from DB in cluters
func (db DB) GetRowsCluster(rowAccess RowAccess) [][]string {
    if(db.DbType == "mysql") {
        return mysqlGetRows(db, rowAccess)
    } else if (db.DbType == "postgres") {
        return postgresGetRows(db, rowAccess)
    } else {
        // should panic or do proper error throwing
        return nil
    }
}

// OptimizeDB : optimzies the database
func (db *DB) OptimizeDB(rankToRowMapArr []map[int]int) {
    if(db.DbType == "mysql") {
        mysqlOptimizeDB(db, rankToRowMapArr)
    } else if (db.DbType == "postgres") {
        // postgresOptimizeDB(db, rankToRowMapArr)
    } else {
        // should panic or do proper error throwing
    }
}

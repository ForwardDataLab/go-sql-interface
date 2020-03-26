package sqlinterface

import (
    "database/sql"
)

func (db DB)BuildConnectionPool() *sql.DB {
    if db.DbType == "mysql" {
        return mysqlBuildConnection(&db)
    } else if db.DbType == "postgres" {
        return nil
    } else {
        // should panic or do proper error throwing
        return nil
    }
}

func (db DB)PrepareQueryMulStmt(currentDB *sql.DB, numQuery int) *sql.Stmt {
    if db.DbType == "mysql" {
        return mysqlPrepareQueryMulStmt(&db, currentDB, numQuery)
    } else if db.DbType == "postgres" {
        return nil
    } else {
        // should panic or do proper error throwing
        return nil
    }
}

func (db DB)PrepareQueryMetaData(currentDB *sql.DB) *sql.Stmt {
    if db.DbType == "mysql" {
        return mysqlPrepareQueryMetaDataStmt(&db, currentDB)
    } else if db.DbType == "postgres" {
        return nil
    } else {
        // should panic or do proper error throwing
        return nil
    }
}

func (db DB)PrepareInsertOneRow(currentDB *sql.DB, numOfCol int) *sql.Stmt {
    if db.DbType == "mysql" {
        return mysqlPrepareInsertOneRow(&db, currentDB, numOfCol)
    } else if db.DbType == "postgres" {
        return nil
    } else {
        // should panic or do proper error throwing
        return nil
    }
}

func (db DB)PrepareDeleteOneRow(currentDB *sql.DB) *sql.Stmt {
    if db.DbType == "mysql" {
        return mysqlPrepareDeleteOneRow(&db, currentDB)
    } else if db.DbType == "postgres" {
        return nil
    } else {
        // should panic or do proper error throwing
        return nil
    }
}

func (db DB)PrepareQueryMaxIndex(currentDB *sql.DB) *sql.Stmt {
    if db.DbType == "mysql" {
        return mysqlPrepareQueryMaxIndex(&db, currentDB)
    } else if db.DbType == "postgres" {
        return nil
    } else {
        // should panic or do proper error throwing
        return nil
    }
}

func (db DB)QueryMetaData(currentDB *sql.DB) *sql.Stmt {
    if db.DbType == "mysql" {
        return mysqlPrepareQueryMaxIndex(&db, currentDB)
    } else if db.DbType == "postgres" {
        return nil
    } else {
        // should panic or do proper error throwing
        return nil
    }
}

func (db DB)QueryNumOfCol(QueryMetaData *sql.Stmt) int{
    if db.DbType == "mysql" {
        return mysqlQueryNumOfCol(QueryMetaData)
    } else if db.DbType == "postgres" {
        return 0
    } else {
        // should panic or do proper error throwing
        return 0
    }
}

func (db DB)QueryMaxIndex(QueryMaxIndexStmt *sql.Stmt) int {
    if db.DbType == "mysql" {
        return mysqlQueryMaxIndex(QueryMaxIndexStmt)
    } else if db.DbType == "postgres" {
        return 0
    } else {
        // should panic or do proper error throwing
        return 0
    }
}

func ExecuteInsertOneRow(dbType string, InsertOneRowStmt *sql.Stmt, Parameters []interface{}) {
    if dbType == "mysql" {
        mysqlInsertOneRow(InsertOneRowStmt, Parameters)
    } else if dbType == "postgres" {

    } else {
        // should panic or do proper error throwing
    }
}

func ExecuteDeleteOneRow(dbType string, DeleteOneRowStmt *sql.Stmt, IDToDelete int) {
    if dbType == "mysql" {
        mysqlDeleteOneRow(DeleteOneRowStmt, IDToDelete)
    } else if dbType == "postgres" {

    } else {
        // should panic or do proper error throwing
    }
}

// InsertColumn : inserts a new column into the database
func InsertColumn(db *sql.DB, dbType string, tableName string, columnName string, columnType string) {
    // insert a column into db defined by columnStructure
    // INSERT INTO table_name (col, col, col) VALUES (NULL, 'my name', 'my group')
    if(dbType == "mysql") {
        mysqlInsertColumn(db, tableName, columnName, columnType)
    } else if (dbType == "postgres") {
        // return postgresInsertColumn(db, column)
        return
    } else {

    }
}

// UpdateRow : updates a row from the database
func UpdateRow(dbType string, DeleteOneRowStmt *sql.Stmt, IDToDelete int, InsertOneRowStmt *sql.Stmt, Parameters []interface{}) {
    // UPDATE table_name WHERE index_col = index
    if dbType == "mysql" {
        mysqlUpdateRow(DeleteOneRowStmt, IDToDelete, InsertOneRowStmt, Parameters)
    } else if dbType == "postgres" {
        // update but postgres
    } else {
        // should panic or do proper error throwing
    }
}


func ExecuteMetaDataStmt(dbType string, MetaDataStmt *sql.Stmt) (*sql.Rows, error) {
    if(dbType == "mysql") {
        return mysqlExecuteMetaDataStmt(MetaDataStmt)
    } else if (dbType == "postgres") {
        // return postgresInsertColumn(db, column)
        return nil, nil
    } else {
        // should panic or do proper error thcolumning
        return nil, nil
    }
}


func ExecuteQueryMulStmt(dbType string, QueryMulStmt *sql.Stmt, QueryIDs []interface{}) (*sql.Rows, error) {
    if(dbType == "mysql") {
        return mysqlExecuteQueryMulStmt(QueryMulStmt, QueryIDs)
    } else if (dbType == "postgres") {
        return nil, nil
    } else {
        return nil, nil
    }
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

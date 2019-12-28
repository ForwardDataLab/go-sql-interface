package sqlinterface


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

package sqlinterface

// GetRows : fetches rows from DB
func (db DB) GetRows(rowAccess RowAccess) []RowStructure {
    if(db.DbType == "mysql") {
        return mysqlGetRows(db, rowAccess)
    } else if (db.DbType == "postgres") {
        return postgresGetRows(db, rowAccess)
    } else {
        // should panic or do proper error throwing
        return nil
    }
}

// InitDB : initializes the database upon initial creation of workspace
func (db DB) InitDB() {
    // add a column called index_col
    // ALTER TABLE `myTable` ADD COLUMN `id` INT AUTO_INCREMENT UNIQUE
    if(db.DbType == "mysql") {
        mysqlInitDB()
    } else if (db.DbType == "postgres") {
        postgresInitDB()
    } else {
        // should panic or do proper error throwing
    }
}

// InsertRow : inserts a new row into the database
func (db DB) InsertRow(rowStructure RowStructure) int {
    // insert a row into db defined by rowStructure
    // INSERT INTO table_name (col, col, col) VALUES (NULL, 'my name', 'my group')
    if(db.DbType == "mysql") {
        return mysqlInsertRow(rowStructure)
    } else if (db.DbType == "postgres") {
        return postgresInsertRow(rowStructure)
    } else {
        // should panic or do proper error throwing
        return -1
    }
}

// DeleteRow : delets a row from the database
func (db DB) DeleteRow(index int) {
    // DELETE FROM table_name WHERE index_col = index
    if(db.DbType == "mysql") {
        mysqlDeleteRow(index)
    } else if (db.DbType == "postgres") {
        postgresDeleteRow(index)
    } else {
        // should panic or do proper error throwing
    }
}

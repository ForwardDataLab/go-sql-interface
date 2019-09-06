package sqlinterface

// GetRows : fetches rows from DB
func (db DB) GetRows(rowAccess RowAccess) []RowStructure {
    if(db.DbType == "mysql") {
        return mysqlGetRows(db, rowAccess)
    } else if (db.DbType == "postgres") {
        return postgresGetRows(db, rowAccess)
    } else {
        return nil
    }
}

// InitDB : initializes the database upon initial creation of workspace
func (db DB) InitDB() {
    // add a column called index_col
    // ALTER TABLE `myTable` ADD COLUMN `id` INT AUTO_INCREMENT UNIQUE
    if(db.DbType == "mysql") {
        return mysqlInitDB()
    } else if (db.DbType == "postgres") {
        return postgresInitDB()
    } else {
        return nil
    }
}

// InsertRow : inserts a new row into the database
func (db DB) InsertRow(rowStructure RowStructure) {
    // insert a row into db defined by rowStructure
    // INSERT INTO table_name (col, col, col) VALUES (NULL, 'my name', 'my group')
    if(db.DbType == "mysql") {
        return mysqlInsertRow(rowStructure)
    } else if (db.DbType == "postgres") {
        return postgresInsertRow(rowStructure)
    } else {
        return nil
    }
}

// DeleteRow : delets a row from the database
func (db DB) DeleteRow(index int) {
    // DELETE FROM table_name WHERE index_col = index
    if(db.DbType == "mysql") {
        return mysqlDeleteRow(index)
    } else if (db.DbType == "postgres") {
        return postgresDeleteRow(index)
    } else {
        return nil
    }
}

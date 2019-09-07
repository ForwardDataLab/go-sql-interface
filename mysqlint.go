package sqlinterface

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func mysqlInitDB(){
}

func mysqlGetRows(db DB, rowAccess RowAccess) []RowStructure {
    currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "@/" + db.DatabaseName)
    queryString := "SELECT * FROM " +
        db.Table +
        " WHERE " + rowAccess.Column + " = ?"
    statement, _ := currentDatabase.Prepare(queryString)
    fetchedArr := make([]RowStructure, len(rowAccess.Indices))
    for index, rowIndex := range rowAccess.Indices {
        statement.QueryRow(rowIndex).Scan(&(fetchedArr[index]).USER_NAME,
            &(fetchedArr[index]).INDEX_COL)
    }
    return fetchedArr
}

func mysqlInsertRow(rowStructure RowStructure) int {
    // INSERT INTO table_name (col, col, col) VALUES (NULL, 'my name', 'my group')
    currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "@/" + db.DatabaseName)
    insertQueryString := "INSERT INTO " +
        db.Table +
        " (USER_NAME, INDEX_COL) VALUES (?, ?)"
    insertStatement, _ := currentDatabase.Prepare(insertQueryString)
    _, _ = insertStatement.Exec(rowStructure.USER_NAME, rowStructure.INDEX_COL)

    selectQueryString := "SELECT INDEX_COL FROM " +
        db.Table +
        " WHERE USER_NAME = ?"
    selectStatement, _ := currentDatabase.Prepare(selectQueryString)

    autoIncrementIndex := -1;
    selectStatement.QueryRow(rowStructure.USER_NAME).Scan(&autoIncrementIndex)
    return autoIncrementIndex
}

func mysqlDeleteRow(index int) {
    // DELETE FROM table_name WHERE index_col = index
    currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "@/" + db.DatabaseName)
    deleteQueryString := "DELETE FROM " +
        db.Table +
        "WHERE INDEX_COL = ?"
    deleteStatement, _ := currentDatabase.Prepare(deleteQueryString)
    _, _ = deleteStatement.Exec(index)
}

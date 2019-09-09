package sqlinterface

import (
    "database/sql"
    "strings"
    _ "github.com/go-sql-driver/mysql"
)

func mysqlInitDB(){
}

func mysqlGetRows(db DB, rowAccess RowAccess) []RowStructure {
    currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "@/" + db.DatabaseName)
    queryString := "SELECT * FROM " +
        db.Table +
        " WHERE " + rowAccess.Column + " in (/" + strings.Repeat(", ?", len(rowAccess.Indices) - 1) + ")"
    statement, _ := currentDatabase.Prepare(queryString)
    fetchedArr := make([]RowStructure, len(rowAccess.Indices))
    rows, _ := statement.Query(rowAccess.Indices...)
    defer rows.Close()
    index := 0
    for rows.Next() {
        rows.Scan(
            &(fetchedArr[index]).ID,
            &(fetchedArr[index]).NAME,
            &(fetchedArr[index]).Age,
            &(fetchedArr[index]).Department,
            &(fetchedArr[index]).GPA,
        )
        index++;
    }
    return fetchedArr
}

func mysqlInsertRow(db DB, rowStructure RowStructure) int {
    // INSERT INTO table_name (col, col, col) VALUES (NULL, 'my name', 'my group')
    // currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password +
    //     "@/" + db.DatabaseName)
    // insertQueryString := "INSERT INTO " +
    //     db.Table +
    //     " (USER_NAME, INDEX_COL) VALUES (?, ?)"
    // insertStatement, _ := currentDatabase.Prepare(insertQueryString)
    // _, _ = insertStatement.Exec(rowStructure.USER_NAME, rowStructure.INDEX_COL)

    // selectQueryString := "SELECT INDEX_COL FROM " +
    //     db.Table +
    //     " WHERE USER_NAME = ?"
    // selectStatement, _ := currentDatabase.Prepare(selectQueryString)

    // autoIncrementIndex := -1;
    // selectStatement.QueryRow(rowStructure.USER_NAME).Scan(&autoIncrementIndex)
    // return autoIncrementIndex
    return -1
}

func mysqlDeleteRow(db DB, index int) {
    // DELETE FROM table_name WHERE index_col = index
    // currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password +
    //     "@/" + db.DatabaseName)
    // deleteQueryString := "DELETE FROM " +
    //     db.Table +
    //     "WHERE INDEX_COL = ?"
    // deleteStatement, _ := currentDatabase.Prepare(deleteQueryString)
    // _, _ = deleteStatement.Exec(index)
}

package sqlinterface

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func mysqlGetRows(db DB, rowAccess RowAccess) []RowStructure {
    currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "@/" + rowAccess.DatabaseName)
    queryString := "SELECT * FROM " +
        rowAccess.Table +
        " WHERE " + rowAccess.Column + " = ?"
    statement, _ := currentDatabase.Prepare(queryString)
    fetchedArr := make([]RowStructure, len(rowAccess.Indices))
    for index, rowIndex := range rowAccess.Indices {
        statement.QueryRow(rowIndex).Scan(&(fetchedArr[index]).USER_NAME,
            &(fetchedArr[index]).INDEX_COL)
    }
    return fetchedArr
}

func mysqlInsertRow(rowStructure RowStructure) {
}

func mysqlDeleteRow(index int) {
}

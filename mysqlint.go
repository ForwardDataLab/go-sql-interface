package sqlinterface

import (
    "database/sql"
    "strings"
    _ "github.com/go-sql-driver/mysql"
)

func mysqlInitDB(){
}

func mysqlGetRowsSerial(db DB, rowAccess RowAccess) [][]string {
    convertedIndices := make([]interface{}, len(rowAccess.Indices))
    for i, v := range rowAccess.Indices {
        convertedIndices[i] = v
    }
    currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "@/" + db.DatabaseName)
    queryString := "SELECT * FROM " +
        db.Table +
        " WHERE " + rowAccess.Column + " = ?"
    statement, _ := currentDatabase.Prepare(queryString)

    columns := nil
    values := nil
    fetchedArr := nil
    returnArr := [][]string{}
    for index, value := range convertedIndices {
        row, _ := statement.QueryRow(value)
        if columns == nil {
            columns, _ := row.Columns()
            values := make([]sql.RawBytes, len(columns))
            fetchedArr := make([]interface{}, len(values))
            for i := range values {
                fetchedArr[i] = &(values[i])
            }
        }
        row.Scan(fetchedArr...)
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

func mysqlGetRowsBatch(db DB, rowAccess RowAccess) [][]string {
    currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "@/" + db.DatabaseName)
    queryString := "SELECT * FROM " +
        db.Table +
        " WHERE " + rowAccess.Column + " in (?" + strings.Repeat(", ?", len(rowAccess.Indices) - 1) + ")"
    statement, _ := currentDatabase.Prepare(queryString)
    rows, _ := statement.Query(rowAccess.Indices...)
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

func mysqlGetRows(db DB, rowAccess RowAccess) [][]string {
    currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "@/" + db.DatabaseName)
    queryString := "SELECT * FROM " +
        db.Table +
        " WHERE " + rowAccess.Column + " in (?" + strings.Repeat(", ?", len(rowAccess.Indices) - 1) + ")"
    statement, _ := currentDatabase.Prepare(queryString)
    rows, _ := statement.Query(rowAccess.Indices...)
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

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

    rows, _ := statement.Query(convertedIndices[0])
    columns, _ := rows.Columns()
    values := make([]sql.RawBytes, len(columns))
    fetchedArr := make([]interface{}, len(values))
    for i := range values {
        fetchedArr[i] = &(values[i])
    }
    returnArr := [][]string{}
    for _, value := range convertedIndices {
        row := statement.QueryRow(value)
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
    convertedIndices := make([]interface{}, len(rowAccess.Indices))
    for i, v := range rowAccess.Indices {
        convertedIndices[i] = v
    }
    currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "@/" + db.DatabaseName)
    queryString := "SELECT * FROM " +
        db.Table +
        " WHERE " + rowAccess.Column + " in (?" + strings.Repeat(", ?", len(convertedIndices) - 1) + ")"
    statement, _ := currentDatabase.Prepare(queryString)
    rows, _ := statement.Query(convertedIndices...)
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
    return mysqlGetRowsBatch(db, rowAccess)
}

func mysqlGetColMap(db DB) map[string]string {
    colMap := make(map[string]string)
    currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "@/" + db.DatabaseName)
    columnQueryString := "SELECT * FROM " + db.Table+ " LIMIT 1"
    rows, _ := currentDatabase.Query(columnQueryString)
    columns, _ := rows.Columns()
    print(columns)
    return colMap
}

func mysqlInsertRow(db DB, indexCol string, cells []Cell) int {
    // INSERT INTO table_name (col, col, col) VALUES (NULL, 'my name', 'my group')
    currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "@/" + db.DatabaseName)
    selectMaxQueryString := "SELECT MAX(" + indexCol + ") FROM " + db.Table
    var maxIndex int
    row, _ := currentDatabase.Query(selectMaxQueryString)
    row.Scan(&maxIndex)

    insertQueryString := "INSERT INTO " +
        db.Table +
        " (?" + strings.Repeat(", ?", len(cells) - 1) + ")" +
        " VALUES (?" + strings.Repeat(", ?", len(cells) - 1) + ")"
    insertStatement, _ := currentDatabase.Prepare(insertQueryString)

    // create interface and add max index
    insertCell := make([]interface{}, len(cells) * 2)
    for i, v := range cell {
        insertCell[i] = v.Column
    }

    for i, v := range cells {
        if v.Column == "int" {
            insertCell[len(cells) + i] = int(v.Value)
        } else if v.Column == "string" {
            insertCell[len(cells) + i] = string(v.Value)
        } else if v.Column == "float" {
            insertCell[len(cells) + i] = float64(v.Value)
        } else if v.Column == "bool" {
            insertCell[len(cells) + i] = bool(v.Value)
        } else {
            return -1
        }
    }
    _, _ = insertStatement.Exec(insertCell...)

    return maxIndex + 1
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

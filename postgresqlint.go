package sqlinterface

import (
    // "database/sql"
    // _ "github.com/lib/pq"
)

func postgresInitDB(){
}

func postgresGetRows(db DB, rowAccess RowAccess) [][]string {
    // currentDatabase, _ := sql.Open(db.DbType, "user=" + db.Username +
    //     " password=" + db.Password +
    //     " dbname=" + db.DatabaseName +
    //     "sslmode=disable")
    // queryString := "SELECT * FROM " +
    //     db.Table +
    //     " WHERE " + rowAccess.Column + " = 1"
    // fetchedArr := make([]RowStructure, len(rowAccess.Indices))
    // for index, rowIndex := range rowAccess.Indices {
    //     row := currentDatabase.QueryRow(queryString, rowIndex)
    //     err := row.Scan(
    //             &(fetchedArr[index]).ID,
    //             &(fetchedArr[index]).NAME,
    //             &(fetchedArr[index]).Age,
    //             &(fetchedArr[index]).Department,
    //             &(fetchedArr[index]).GPA,
    //     )
    //     if err != nil {
    //         print(err)
    //     }
    // }
    // return fetchedArr
    return nil
}

func postgresInsertRow(db DB, row []string) int {
    return -1
}

func postgresDeleteRow(db DB, index int) {
}

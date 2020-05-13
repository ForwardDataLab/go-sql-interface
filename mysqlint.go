package sqlinterface

import (
    "database/sql"
    "fmt"
    "math/rand"
    "strconv"
    "strings"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

func mysqlInitDB(db *DB){
    db.fresh = true
}

func mysqlBuildConnection(db *DB) *sql.DB{
    currentDatabase, err := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "@tcp(" + db.Host + ":" + db.Port + ")/" + db.DatabaseName)
    if err != nil {
        fmt.Println(err)
    }
    return currentDatabase
}

func mysqlPrepareQueryMulStmt(db *DB, currentDB *sql.DB, numQuery int, idColumn string, tableMetadata []TableMetadata) *sql.Stmt{
    QueryMulString := "SELECT " + tableMetadata[0].Field.String
    for i := 1; i < len(tableMetadata); i ++ {
        QueryMulString += ("," + tableMetadata[i].Field.String)
    }
    QueryMulString += " FROM " + db.Table + " WHERE " + idColumn + " in (?"
    for i := 1; i < numQuery; i++ {
        QueryMulString += ", ?"
    }
    QueryMulString += ")"
    QueryMulStmt, err := currentDB.Prepare(QueryMulString)
    if err != nil {
        fmt.Println(err)
    }
    return QueryMulStmt
}

// func mysqlPrepareQueryMetaDataStmt(db *DB, currentDB *sql.DB) *sql.Stmt {
//     QueryMetaDataString := "DESCRIBE " + db.Table
//     MetaDataStmt, err := currentDB.Prepare(QueryMetaDataString)
//     if err != nil {
//         fmt.Println(err)
//     }
//     return MetaDataStmt
// }

func mysqlPrepareUpdateOneRow(db *DB, currentDB *sql.DB, columnNames []string, idColumnName string) *sql.Stmt {
    updateString := "UPDATE " + db.Table + " SET " + columnNames[0] + " = ?"
    for i := 1; i < len(columnNames); i++ {
        updateString += ("," + columnNames[i] + " = ?")
    }
    updateString += (" WHERE " + idColumnName + " = ?");

    UpdateOneRowStmt, err := currentDB.Prepare(updateString)
    if err != nil {
        fmt.Println("update string:", updateString)
        fmt.Println("Prepare update one row:", err)
    }
    return UpdateOneRowStmt
}

func mysqlPrepareUpdateCell(db *DB, currentDB *sql.DB, columnName string, idColumnName string) *sql.Stmt {
    updateString := "UPDATE " + db.Table + " SET " + columnName + " = ?"
    updateString += (" WHERE " + idColumnName + " = ?");

    UpdateCellStmt, err := currentDB.Prepare(updateString)
    if err != nil {
        fmt.Println("update string:", updateString)
        fmt.Println("Prepare update cell", err)
    }
    return UpdateCellStmt
}

func mysqlPrepareInsertOneRow(db *DB, currentDB *sql.DB, numOfCol int) *sql.Stmt {
    insertString := "INSERT INTO " + db.Table + " VALUES (?"
    for i := 1; i < numOfCol; i++ {
        insertString += ", ?"
    }
    insertString += ")"
    InsertOneStmt, err := currentDB.Prepare(insertString)
    if err != nil {
        fmt.Println(err)
    }
    return InsertOneStmt
}

func mysqlPrepareDeleteOneRow(db *DB, currentDB *sql.DB) *sql.Stmt {
    deleteString := "DELETE FROM " + db.Table + " WHERE ID = ?"
    DeleteOneStmt, err := currentDB.Prepare(deleteString)
    if err != nil {
        fmt.Println(err)
    }
    return DeleteOneStmt
}

func mysqlPrepareQueryMaxIndex(db *DB, currentDB *sql.DB) *sql.Stmt {
    QueryMaxIndexString := "SELECT MAX(ID) FROM " + db.Table
    QueryMaxIndexStmt, err := currentDB.Prepare(QueryMaxIndexString)
    if err != nil {
        fmt.Println(err)
    }
    return QueryMaxIndexStmt
}

func mysqlQueryNumOfCol(MetaDataStmt *sql.Stmt) int {
    ColMap, err := MetaDataStmt.Query()
    if err != nil {
        fmt.Println(err)
    }
    numOfCol := 0
    for ColMap.Next() {
        numOfCol++
    }
    return numOfCol
}

func mysqlQueryMaxIndex(QueryMaxIndexStmt *sql.Stmt) int {
    rows, err := QueryMaxIndexStmt.Query()
    defer rows.Close()
    if err != nil {
        fmt.Println(err)
        return 0
    } else {
        var values sql.RawBytes
        var fetchedArr interface{}
        fetchedArr = &values
        var maxIndex int
        for rows.Next() {
            rows.Scan(fetchedArr)
            maxIndex, err = strconv.Atoi(string(values))
            if err != nil {
                fmt.Println(err)
            }
        }
        return maxIndex
    }
}

func mysqlInsertOneRow(InsertOneRowStmt *sql.Stmt, InsertPara []interface{}) {
    InsertOneRowStmt.Exec(InsertPara...)
}

func mysqlDeleteOneRow(DeleteOneRowStmt *sql.Stmt, IDToDelete interface{}) {
    DeleteOneRowStmt.Exec(IDToDelete)
}

func mysqlDeleteOneColumn(db *DB, currentDB *sql.DB, ColumnName string) int {
    queryString := "ALTER TABLE " + db.Table + " DROP COLUMN " + ColumnName;
    _, err := currentDB.Exec(queryString)
    if (err != nil) {
        fmt.Println(err);
    }
    return 0
}

func mysqlUpdateCell(db DB, cell Cell, rowIDToUpdate interface{}, UpdateCellStmt *sql.Stmt) {
    updateCellParameters := make([]interface{}, 2)
    updateCellParameters[0] = cell.Value
    updateCellParameters[1] = rowIDToUpdate
    fmt.Println(UpdateCellStmt)
    fmt.Println(updateCellParameters)
    if _, err := UpdateCellStmt.Exec(updateCellParameters...); err != nil {
        fmt.Println("Update cell error: ", err)
    }
}

func mysqlUpdateRow(db DB, cells []Cell, UpdateOneRowStmt *sql.Stmt) {
   updateRow := make([]interface{}, len(cells) + 1)
   var idColumnValue interface{};
   for i, v := range cells {
       if v.Type == "ID" {
           updateRow[i] = v.Value
           idColumnValue = v.Value
       } else {
           updateRow[i] = v.Value
       }

   }
   updateRow[len(updateRow) - 1] = idColumnValue
   if _, err := UpdateOneRowStmt.Exec(updateRow...); err != nil {
       fmt.Println(err)
   }
}

func mysqlInsertColumn(db DB, DBPool *sql.DB, columnName string, columnType string) int {
    queryString := "ALTER TABLE " + db.Table + " ADD COLUMN " + columnName + " " + columnType;
    _, err := DBPool.Exec(queryString)
    if (err != nil) {
        fmt.Println(err);
    }
    return 0
}

func mysqlGetMetadata(db DB, currentDB *sql.DB) []TableMetadata {
    QueryMetaDataString := "DESCRIBE " + db.Table
    preparedMetadataQuery, err := currentDB.Prepare(QueryMetaDataString)
    rows, err := preparedMetadataQuery.Query()
    if err != nil {
        fmt.Println(err)
        return nil
    } else {
        var tableMetadata TableMetadata
        var returnArr []TableMetadata
        for rows.Next() {
            rows.Scan(
                &(tableMetadata.Field),
                &(tableMetadata.Type),
                &(tableMetadata.Null),
                &(tableMetadata.Key),
                &(tableMetadata.Default),
                &(tableMetadata.Extra),
            )
            returnArr = append(returnArr, tableMetadata)
        }
        return returnArr
    }
}

func mysqlExecuteQueryMulStmt(QueryMulStmt *sql.Stmt, QueryIDs []interface{}) [][]string {
    rows, err := QueryMulStmt.Query(QueryIDs...)
    if (err == nil) {
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
    } else {
        fmt.Println(err)
        return nil
    }
}


func mysqlEvaluateFormula(db DB, formulaData string) [][]string {
    currentDatabase, err := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "@tcp(" + db.Host + ":" + db.Port + ")/" + db.DatabaseName)
    queryString := formulaData
    if (err != nil) {
        fmt.Println(err);
    }
    statement, err := currentDatabase.Prepare(queryString)
    if (err != nil) {
        fmt.Println(err);
    }
    rows, err := statement.Query()
    if (err == nil) {
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
    return nil
}

func cost(clusterConfiguration []int, rowToRankMapArr []map[int]int) int {
    weights := []int{1, 1, 1}
    coefficients := []int{0, 0, 0}

    clusterMapping := make(map[int]int)
    for _, v := range clusterConfiguration {
        if currentNum, ok := clusterMapping[v]; ok {
            clusterMapping[v] = currentNum + 1
        } else {
            clusterMapping[v] = 1
        }
    }
    currentMaxSize := 0
    for _, v := range clusterMapping {
        if v > currentMaxSize {
            currentMaxSize = v
        }
    }
    sumDifference := 0
     for _, mapping := range rowToRankMapArr {
         for i, v := range clusterConfiguration {
             // fmt.Print("i: ")
             // fmt.Print(i)
             // fmt.Print("v: ")
             // fmt.Print(v)
             // fmt.Print("mapping[i + 1]: ")
             // fmt.Println(mapping[i + 1])
             // fmt.Println(v - mapping[i + 1])
             sumDifference += ((v - mapping[i + 1]) * (v - mapping[i + 1]))
         }
     }

     coefficients[0] = currentMaxSize
     coefficients[1] = len(clusterMapping)
     coefficients[2] = sumDifference
     fmt.Print("coefficients: ")
     fmt.Println(coefficients)
     // fmt.Print("Associated configuration: ")
     // fmt.Println(clusterConfiguration)
     finalCost := 0
     for i := 0; i < 3; i ++ {
         finalCost += weights[i] * coefficients[i]
     }
     return finalCost
}

func calculateOptimalClusterSize(numRows int) int {
    if numRows < 10000 {
        return numRows / 3
    }
    return numRows / 10
}

func mysqlOptimizeDB(db *DB, rowToRankMapArr []map[int]int) {
    currentMinimumConfiguration := make([]int, len(rowToRankMapArr[0]))
    if db.fresh {
        db.ClusterMap = make(map[interface{}]int)
        db.ClusterSize = calculateOptimalClusterSize(len(rowToRankMapArr[0]))
        db.NumClusters = len(rowToRankMapArr[0]) / db.ClusterSize
        for i := 0; i < len(rowToRankMapArr[0]); i ++ {
            currentMinimumConfiguration[i] = -9999
        }
        db.fresh = false
    } else {
        for i := 0; i < len(rowToRankMapArr[0]); i ++ {
            currentMinimumConfiguration[i] = db.ClusterMap[i]
        }
    }
    pickMinimumCost(db, currentMinimumConfiguration, 0, db.NumClusters, db.ClusterSize, rowToRankMapArr)
    fmt.Print("New configuration: ")
    fmt.Println(db.newConfiguration)
    for i, v := range db.newConfiguration {
        db.ClusterMap[i + 1] = v
    }

    currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "@tcp(" + db.Host + ":" + db.Port + ")/" + db.DatabaseName)
    currentTransaction, _ := currentDatabase.Begin()
    fmt.Println("Beginning transaction")
    for i, v := range db.newConfiguration {
        updateString := "UPDATE " + db.Table + " SET cluster_id = ? WHERE id = ?"
        updateStatement, _ := currentTransaction.Prepare(updateString)
        updateStatement.Exec(v, i + 1)
    }
    fmt.Println("End transaction")
    currentTransaction.Commit()
    // set the db.ClusterMap to the minimum configuratino found
}

func pickMinimumCost(db *DB, currentConfiguration []int, numIter int, numClusters int, clusterSize int, rowToRankMapArr []map[int]int) {
    if numIter < 100 {
        currentCost := cost(currentConfiguration, rowToRankMapArr)
        newConfiguration := getConfiguration(len(currentConfiguration), numIter, numClusters, clusterSize)
        newCost := cost(newConfiguration, rowToRankMapArr)
        if newCost < currentCost {
            fmt.Print("proposed configuration: ")
            fmt.Println(newConfiguration)
            fmt.Print("cost of proposed configuration: ")
            fmt.Println(newCost)
            db.newConfiguration = newConfiguration
            pickMinimumCost(db, newConfiguration, numIter + 1, numClusters, clusterSize, rowToRankMapArr)
        } else {
            pickMinimumCost(db, currentConfiguration, numIter + 1, numClusters, clusterSize, rowToRankMapArr)
        }
    }
}

func getConfiguration(lengthConfiguration int, numIter int, numClusters int, clusterSize int) []int {
    rand.Seed(time.Now().UTC().UnixNano())
    newConfiguration := make([]int, lengthConfiguration)
    // currentValue := numIter
    // for i := lengthConfiguration - 1; i > -1; i -- {
    //     newConfiguration[i] = int(math.Mod(float64(currentValue), float64(numClusters)))
    //     currentValue = (currentValue - int(math.Mod(float64(currentValue), float64(numClusters)))) / numClusters
    for i := 0; i < lengthConfiguration; i ++ {
        newConfiguration[i] = i / clusterSize
    }
    shuffledConfiguration := make([]int, len(newConfiguration))
    perm := rand.Perm(len(newConfiguration))
    for i, v := range perm {
        shuffledConfiguration[v] = newConfiguration[i]
    }
    // fmt.Println(shuffledConfiguration)
    return shuffledConfiguration
}

func mysqlGetRowsSerial(db DB, rowAccess RowAccess) [][]string {
    convertedIndices := make([]interface{}, len(rowAccess.Indices))
    for i, v := range rowAccess.Indices {
        convertedIndices[i] = v
    }
    currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "@tcp(" + db.Host + ":" + db.Port + ")/" + db.DatabaseName)
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
    currentDatabase, err := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "@tcp(" + db.Host + ":" + db.Port + ")/" + db.DatabaseName)
    queryString := "SELECT * FROM " +
        db.Table +
        " WHERE " + rowAccess.Column + " in (?" + strings.Repeat(", ?", len(convertedIndices) - 1) + ")"
    if (err != nil) {
        fmt.Println(err);
    }
    statement, err := currentDatabase.Prepare(queryString)
    if (err != nil) {
        fmt.Println(err);
    }
    rows, err := statement.Query(convertedIndices...)
    if (err == nil) {
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
    return nil
}

func mysqlGetRowsCluster(db DB, rowAccess RowAccess) [][]string {
    clusterIDs := []int{}
    subsetClusterMap := make(map[int]bool)
    rowIDMap := make(map[interface{}]bool)
    for _, v := range rowAccess.Indices {
        rowIDMap[v] = true
    }

    for _, v := range rowAccess.Indices {
        subsetClusterMap[db.ClusterMap[v]] = true
    }

    for k := range subsetClusterMap {
        clusterIDs = append(clusterIDs, k)
    }

    convertedIndices := make([]interface{}, len(clusterIDs))
    for i, v := range clusterIDs {
        convertedIndices[i] = v
    }
    currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "@tcp(" + db.Host + ":" + db.Port + ")/" + db.DatabaseName)
    queryString := "SELECT * FROM " +
        db.Table +
        " WHERE cluster_id in (?" + strings.Repeat(", ?", len(convertedIndices) - 1) + ")"
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
    filteredArr := [][]string{}
    // // PERFORM FILTERING BASED ON rowAccess.Indices
    fmt.Println(returnArr)
    fmt.Println(rowIDMap)
    for _, v := range returnArr {
        index, _ := strconv.ParseInt(v[0], 10, 64)
        fmt.Println(index)
        if _, ok := rowIDMap[int(index)]; ok {
            filteredArr = append(filteredArr, v)
        }
    }

    return filteredArr
}

func mysqlGetRows(db DB, rowAccess RowAccess) [][]string {
    return mysqlGetRowsBatch(db, rowAccess)
}

func mysqlGetColumns(db DB, columnAccess ColumnAccess) [][]string {
    return mysqlGetColumnsBatch(db, columnAccess)
}

func mysqlGetColumnsBatch(db DB, columnAccess ColumnAccess) [][]string {
    colunmNames := columnAccess.ColumnNames
    numOfColumn := len(colunmNames)

    currentDatabase, err := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "@tcp(" + db.Host + ":" + db.Port + ")/" + db.DatabaseName)
    queryString := "SELECT " + colunmNames[0]
    for i := 1; i < numOfColumn; i ++ {
        queryString += ", " + colunmNames[i]
    }

    queryString += " FROM " + db.Table
    if (err != nil) {
        fmt.Println(err);
    }
    statement, err := currentDatabase.Prepare(queryString)
    if (err != nil) {
        fmt.Println(err);
    }
    rows, err := statement.Query()
    if (err == nil) {
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
    return nil
}


func mysqlInsertRow(cells []Cell, maxIndex int, InsertOneStmt *sql.Stmt) int {
    insertCell := make([]interface{}, len(cells))

    for i, v := range cells {
        if v.Type == "ID" {
            insertCell[i] = maxIndex + 1
        } else {
            insertCell[i] = v.Value
        }
    }
    _, _ = InsertOneStmt.Exec(insertCell...)
    return maxIndex + 1
}


func mysqlDeleteRow(db DB, indexCol string, index int) {
    // DELETE FROM table_name WHERE index_col = index
    currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "@tcp(" + db.Host + ":" + db.Port + ")/" + db.DatabaseName)
    deleteQueryString := "DELETE FROM " +
        db.Table +
        " WHERE " + indexCol + " = ?"
    deleteStatement, _ := currentDatabase.Prepare(deleteQueryString)
    _, err := deleteStatement.Exec(index)
    if (err != nil) {
        fmt.Println(err);
    }
}


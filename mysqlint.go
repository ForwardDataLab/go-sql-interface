package sqlinterface

import (
        "database/sql"
        "fmt"
        "math/rand"
        "strings"
        "strconv"
        "time"
        _ "github.com/go-sql-driver/mysql"
       )

func mysqlInitDB(db *DB){
    db.fresh = true
}

func cost(clusterConfiguration []int, rankToRowMapArr []map[int]int) int {
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
    for _, v := range clusterConfiguration {
        if v > currentMaxSize {
            currentMaxSize = v
        }
    }
    sumDifference := 0
     for _, mapping := range rankToRowMapArr {
         for i, v := range clusterConfiguration {
             sumDifference += (v - mapping[i]) * (v - mapping[i])
         }
     }

     coefficients[0] = currentMaxSize
     coefficients[1] = len(clusterMapping)
     coefficients[2] = sumDifference
     finalCost := 0
     for i := 0; i < 3; i ++ {
         finalCost += weights[i] * coefficients[i]
     }
     return finalCost
}

func calculateOptimalClusterSize(numRows int) int {
    if numRows < 100 {
        return numRows / 3
    }
    return -1
}

func mysqlOptimizeDB(db *DB, rankToRowMapArr []map[int]int) {
currentMinimumConfiguration := make([]int, len(rankToRowMapArr[0]))
    if db.fresh {
        db.ClusterMap = make(map[int]int)
        db.ClusterSize = calculateOptimalClusterSize(len(rankToRowMapArr[0]))
        db.NumClusters = len(rankToRowMapArr[0]) / db.ClusterSize
        for i := 0; i < len(rankToRowMapArr[0]); i ++ {
            currentMinimumConfiguration[i] = 0
        }
    } else {
        for i := 0; i < len(rankToRowMapArr[0]); i ++ {
            currentMinimumConfiguration[i] = db.ClusterMap[i]
        }
    }
    pickMinimumCost(db, currentMinimumConfiguration, 0, db.NumClusters, rankToRowMapArr)
    fmt.Print("New configuration: ")
    fmt.Println(db.newConfiguration)
    // set the db.ClusterMap to the minimum configuratino found
}

func pickMinimumCost(db *DB, currentConfiguration []int, numIter int, numClusters int, rankToRowMapArr []map[int]int) int {
    if numIter > 50 {
        return -1;
    }
    currentCost := cost(currentConfiguration, rankToRowMapArr)
    newConfiguration := getConfiguration(len(currentConfiguration), numIter, numClusters)
    newCost := pickMinimumCost(db, newConfiguration, numIter + 1, numClusters, rankToRowMapArr)
    if newCost == -1 {
        return currentCost
    }
    if newCost < currentCost {
        db.newConfiguration = newConfiguration
    }
    return min(currentCost, newCost)
}

func getConfiguration(lengthConfiguration int, numIter int, numClusters int) []int {
    rand.Seed(time.Now().UTC().UnixNano())
    newConfiguration := make([]int, lengthConfiguration)
    // currentValue := numIter
    // for i := lengthConfiguration - 1; i > -1; i -- {
    //     newConfiguration[i] = int(math.Mod(float64(currentValue), float64(numClusters)))
    //     currentValue = (currentValue - int(math.Mod(float64(currentValue), float64(numClusters)))) / numClusters
    for i := 0; i < lengthConfiguration; i ++ {
        newConfiguration[i] = i / numClusters
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

func mysqlGetRowsCluster(db DB, rowAccess RowAccess) [][]string {
    // clusterIDs := []int{}
    // subsetClusterMap := make(map[int]bool)
    // rowIDMap := make(map[int]bool)
    // for i, v := range rowAccess.Indices {
    //     rowIDMap[v] = true
    // }

    // for i, v := range rowAccess.Indices {
    //     subsetClusterMap[db.ClusterMap[v]] = true
    // }

    // for k, v := range subsetClusterMap {
    //     clusterIDs = append(clusterIDs, k)
    // }

    // convertedIndices := make([]interface{}, len(clusterIDs))
    // for i, v := range clusterIDs {
    //     convertedIndices[i] = v
    // }
    // currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password +
    //     "@/" + db.DatabaseName)
    // queryString := "SELECT * FROM " +
    //     db.Table +
    //     " WHERE " + rowAccess.Column + " in (?" + strings.Repeat(", ?", len(convertedIndices) - 1) + ")"
    // statement, _ := currentDatabase.Prepare(queryString)
    // rows, _ := statement.Query(convertedIndices...)
    // columns, _ := rows.Columns()
    // values := make([]sql.RawBytes, len(columns))
    // defer rows.Close()
    // fetchedArr := make([]interface{}, len(values))
    // for i := range values {
    //     fetchedArr[i] = &(values[i])
    // }

    // returnArr := [][]string{}
    // for rows.Next() {
    //     rows.Scan(fetchedArr...)
    //     currentArr := []string{}
    //     for _, v := range values {
    //         if v == nil {
    //             currentArr = append(currentArr, "NULL")
    //         } else {
    //             currentArr = append(currentArr, string(v))
    //         }
    //     }
    //     returnArr = append(returnArr, currentArr)
    // }
    filteredArr := [][]string{}
    // // PERFORM FILTERING BASED ON rowAccess.Indices
    // for i, v := range returnArr {
    //     if val, ok = rowIDMap[v[0]]; ok {
    //         filteredArr = append(filteredArr, v)
    //     }
    // }

    return filteredArr
}

func mysqlGetRows(db DB, rowAccess RowAccess) [][]string {
    return mysqlGetRowsCluster(db, rowAccess)
}

func mysqlGetColMap(db DB) []string {
    // colMap := make(map[string]string)
    currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "@/" + db.DatabaseName)
    columnQueryString := "SELECT * FROM " + db.Table+ " LIMIT 1"
    rows, _ := currentDatabase.Query(columnQueryString)
    columns, _ := rows.Columns()
    return columns
}

func mysqlInsertRow(db DB, indexCol string, cells []Cell) int {
    // INSERT INTO table_name (col, col, col) VALUES (NULL, 'my name', 'my group')
    currentDatabase, _ := sql.Open(db.DbType, db.Username + ":" + db.Password +
        "@/" + db.DatabaseName)
    selectMaxQueryString := "SELECT MAX(" + indexCol + ") FROM " + db.Table
    var maxIndex int
    rows, _ := currentDatabase.Query(selectMaxQueryString)
    for rows.Next() {
        rows.Scan(&maxIndex)
    }

    // " (?" + strings.Repeat(", ?", len(cells) - 1) + ")" +
    insertQueryString := "INSERT INTO " +
        db.Table +
        " VALUES (?" + strings.Repeat(", ?", len(cells) - 1) + ")"
    insertStatement, err := currentDatabase.Prepare(insertQueryString)
    if err != nil {
        fmt.Println(insertQueryString)
        fmt.Println(err)
    }

    // create interface and add max index
    insertCell := make([]interface{}, len(cells))

    for i, v := range cells {
        if v.Type == "ID" {
            insertCell[i] = maxIndex + 1
        } else if v.Type == "int" {
            insertCell[i], _ = strconv.ParseInt(v.Value, 10, 32)
        } else if v.Type == "string" {
            insertCell[i] = string(v.Value)
        } else if v.Type == "float" {
            insertCell[i], _ = strconv.ParseFloat(v.Value, 64)
        } else if v.Type == "bool" {
            insertCell[i], _ = strconv.ParseBool(v.Value)
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

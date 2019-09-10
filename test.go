package sqlinterface

import (
    "fmt"
    "time"
)

// InterfaceTest : tests the interface import
func InterfaceTest() int {
    return 5
}

func InterfaceBatchTime() {
    fmt.Println("Starting batch test")
    totalDuration := 0.0
    numIterations := 10

    for i := 0; i < numIterations; i ++ {
        start := time.Now()
        mysqlDb := DB{
            DbType: "mysql",
            DatabaseName: "dataspread_local",
            Table: "second_sheet",
            Username: "dataspread-user-1",
            Password: "dataspread-pass-1",
        }
        ra := RowAccess{
            Column: "ID",
            Indices: []int{1, 2, 3, 4, 5},
        }
        mysqlDb.GetRowsBatch(ra)
        currentTime := time.Since(start)
        totalDuration += currentTime.Seconds()
        fmt.Println(currentTime)
    }
    fmt.Println("total average time is: %f", totalDuration / float64(numIterations))
}

func InterfaceSerialTime() {
    fmt.Println("Starting Serial test")
    totalDuration := 0.0
    numIterations := 10

    for i := 0; i < numIterations; i ++ {
        start := time.Now()
        mysqlDb := DB{
            DbType: "mysql",
            DatabaseName: "dataspread_local",
            Table: "second_sheet",
            Username: "dataspread-user-1",
            Password: "dataspread-pass-1",
        }
        ra := RowAccess{
            Column: "ID",
            Indices: []int{1, 2, 3, 4, 5},
        }
        mysqlDb.GetRowsSerial(ra)
        currentTime := time.Since(start)
        totalDuration += currentTime.Seconds()
        fmt.Println(currentTime)
    }
    fmt.Println("total average time is: %f", totalDuration / float64(numIterations))
}

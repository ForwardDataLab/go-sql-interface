package sqlinterface

import (
    "fmt"
    "time"
)

// InterfaceTest : tests the interface import
func InterfaceTest() int {
    return 5
}

// InterfaceBatchTimeTest : tests the batch interface get rows times
func InterfaceBatchTimeTest(db DB, ra RowAccess) {
    fmt.Println("Starting batch test")
    totalDuration := int64(0)
    numIterations := 10

    for i := 0; i < numIterations; i ++ {
        start := time.Now()
        db.GetRowsBatch(ra)
        currentTime := time.Since(start)
        totalDuration += currentTime.Milliseconds()
        fmt.Println(currentTime)
    }
    fmt.Println("total average time is: ", totalDuration / int64(numIterations))
}

// InterfaceSerialTimeTest : tests the serial interface get rows times
func InterfaceSerialTimeTest(db DB, ra RowAccess) {
    fmt.Println("Starting Serial test")
    totalDuration := int64(0)
    numIterations := 10

    for i := 0; i < numIterations; i ++ {
        start := time.Now()
        db.GetRowsSerial(ra)
        currentTime := time.Since(start)
        totalDuration += currentTime.Milliseconds()
        fmt.Println(currentTime)
    }
    fmt.Println("total average time is: ", totalDuration / int64(numIterations))
}

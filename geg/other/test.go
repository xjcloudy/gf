package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gbinlog"
    "gitee.com/johng/gf/g/encoding/gbinary"
)

func main() {
    b, err := gbinlog.New("/tmp/gbinlog", "gbinlog")
    fmt.Println(err)
    //b.SetCap(102400)
    for i := 0; i < 100; i++ {
        fmt.Println(b.Push(gbinary.EncodeInt(i)))
    }
    //events2 := make(chan int, 100)
    //go func() {
    //    for{
    //        v := <- events1
    //        fmt.Println(v)
    //    }
    //
    //}()

    //go func() {
    //    time.Sleep(2*time.Second)
    //    events1 <- 1
    //    events2 <- 2
    //    time.Sleep(2*time.Second)
    //    close(events1)
    //    close(events2)
    //    events1 <- 1
    //    events2 <- 2
    //}()
    //
    //select {
    //
    //}
}
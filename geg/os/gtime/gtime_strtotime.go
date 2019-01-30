package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gtime"
    "time"
)

func main() {
    array := []string{
        "2017-12-14 04:51:34 +0805 LMT",
        "2006-01-02T15:04:05Z07:00",
        "2014-01-17T01:19:15+08:00",
        "2018-02-09T20:46:17.897Z",
        "2018-02-09 20:46:17.897",
        "2018-02-09T20:46:17Z",
        "2018-02-09 20:46:17",
        "2018.02.09 20:46:17",
        "2018-02-09",
        "2017/12/14 04:51:34 +0805 LMT",
        "2018/02/09 12:00:15",
        "18/02/09 12:16",
        "18/02/09 12",
        "18/02/09 +0805 LMT",
        "01/Nov/2018:13:28:13 +0800",
        "01-Nov-2018 11:50:28 +0805 LMT",
        "01-Nov-2018T15:04:05Z07:00",
        "01-Nov-2018T01:19:15+08:00",
        "01-Nov-2018 11:50:28 +0805 LMT",
        "01/Nov/18 11:50:28",
        "01/Nov/2018 11:50:28",
        "01/Nov/2018:11:50:28",
        "01.Nov.2018:11:50:28",
        "01/Nov/2018",
    }
    cstLocal, _ := time.LoadLocation("Asia/Shanghai")
    for _, s := range array {
        fmt.Println(s)
        if t, err := gtime.StrToTime(s); err == nil {
            fmt.Println(t.String())
            fmt.Println(t.In(cstLocal).String())
        } else {
            panic(err)
        }
        fmt.Println()
    }
}

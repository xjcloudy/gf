package main

import (
    "gitee.com/johng/gf/third/github.com/fsnotify/fsnotify"
    "log"
    "gitee.com/johng/gf/g/os/glog"
)

func main() {
    //创建一个监控对象
    watch, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watch.Close()
    //添加要监控的对象，文件或文件夹
    err = watch.Add("D:\\Workspace\\Go\\GOPATH\\src\\gitee.com\\johng\\gf\\geg\\other\\test.go")
    if err != nil {
        log.Fatal(err)
    }
    //我们另启一个goroutine来处理监控对象的事件
    go func() {
        for {
            select {
            case ev := <-watch.Events:
                glog.Println(ev)

            case err := <-watch.Errors:
                log.Println("error : ", err)
                return

            }
        }
    }()

    //循环
    select {}
}
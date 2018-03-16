package main

import (
    "gitee.com/johng/gf/g/container/gdrh"
    "fmt"
)

func main () {
    m := gdrh.New(2, 2)
    m.Set(1, 11)
    m.Set(2, 22)
    m.Set(3, 33)
    m.Set(4, 44)
    fmt.Println(m.Get(4))
}

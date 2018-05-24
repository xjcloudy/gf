package main

import (
    "gitee.com/johng/gf/g/encoding/gxml"
    "fmt"
)

func main() {
    xmlContent := `
<?xml version="1.0" encoding="gbk"?>
<config>
    <Test>我爱gf</Test>
</config>
`
    //mxj.XmlCharsetReader("gbk", nil)
    bytes, err := gxml.ToJson([]byte(xmlContent))
    fmt.Println(err)
    fmt.Println(string(bytes))
}
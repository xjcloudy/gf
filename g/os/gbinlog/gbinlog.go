// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 高可用的binlog.
// BinLog文件结构：[数据状态(8bit) 数据长度(32bit) 事务编号(64bit) 数据字段(变长) 事务编号(64bit)] 注意binlog中的事务编号不是递增的，但是是唯一的
package gbinlog

import (
    "os"
    "sync"
    "errors"
    "gitee.com/johng/gf/g/encoding/gbinary"
    "gitee.com/johng/gf/g/container/gtype"
    "strings"
    "strconv"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/gfilepool"
    "fmt"
)

const (
    gDEFAULT_MAX_CAP    = 1024*1024*1024        // 默认允许的最大binlog文件大小(1GB)
    gDEFAULT_FILE_FLAGS = os.O_RDWR|os.O_CREATE // 文件操作默认flag
    gDEFAULT_BLOCK_SIZE = 1024*4                // 默认读取的数据块大小
)

// binlog操作对象
type BinLog struct {
    mu      sync.RWMutex     // 互斥锁
    cap     *gtype.Int       // 单binlog文件允许的最大容量，超过该容量则新建文件
    path    string           // 日志目录(绝对路径)
    name    string           // binlog文件名
    head    *Log             // 最老日志文件操作对象
    tail    *Log             // 最新日志文件操作对象
    logs    *gmap.IntInterfaceMap
    popitem *LogItem         // 最后一条Pop的LogItem
}

// log文件操作对象(一个binlog文件)
type Log struct {
    mu      sync.RWMutex     // 互斥锁
    fp      *gfilepool.Pool  // 文件指针池
    num     int              // 文件编号
    size    *gtype.Int       // 当前文件大小
    blog    *BinLog          // 所属binlog对象
}

// log数据项
type LogItem struct {
    log     *Log             // 所log对象
    Id      int64            // 数据项ID(48bit文件偏移量 + 16bit文件编号)
    Status  int              // 数据项状态
    Index   int64            // 事务在binlog文件的开始位置
    Buffer  []byte           // 二进制数据
}

// 创建binlog对象
func New(path string, name string, cap...int) (*BinLog, error) {
    // 目录权限检测
    if !gfile.Exists(path) {
        gfile.Mkdir(path)
    }
    if !gfile.IsWritable(path) {
        return nil, errors.New(path + " is not writable for saving binlog")
    }
    // 检索文件编号及对象指针
    files := gfile.ScanDir(path)
    min   := 0
    max   := 0
    for _, file := range files {
        array := strings.Split(file, ".")
        if strings.Compare(name, array[0]) == 0 {
            if num, err := strconv.Atoi(array[1]); err == nil {
                if num > max {
                    max = num
                }
            }
        }
    }
    // 初始化文件打开指针
    blog := &BinLog {
        cap  : gtype.NewInt(gDEFAULT_MAX_CAP),
        path : path,
        name : name,
        logs : gmap.NewIntInterfaceMap(),
    }
    if len(cap) > 0 {
        blog.cap.Set(cap[0])
    }
    blog.head = blog.getLog(min)
    blog.tail = blog.head
    if max != min {
        blog.tail = blog.getLog(max)
    }
    return blog, nil
}

// 设置单log文件最大存储大小
func (blog *BinLog) SetCap(cap int) {
    blog.cap.Set(cap)
}

// 创建log对象
func (blog *BinLog) getLog(num int) *Log {
    if v := blog.logs.Get(num); v != nil {
        return v.(*Log)
    }
    path := blog.path + gfile.Separator + blog.name + "." + strconv.Itoa(num)
    log  := &Log{
        num  : num,
        blog : blog,
        size : gtype.NewInt(int(gfile.Size(path))),
    }
    log.fp = gfilepool.New(path, gDEFAULT_FILE_FLAGS, 60)
    blog.logs.Set(num, log)
    return log
}

// 关闭binlog
func (blog *BinLog) Close() {
    blog.logs.Iterator(func(k int, v interface{}) {
        v.(*Log).close()
    })
}

// 关闭log
func (log *Log) close() {
    log.mu.Lock()
    log.fp.Close()
    log.mu.Unlock()
}

// 写入数据
func (blog *BinLog) Append(data []byte, sync...bool) (int64, error) {
    return blog.AppendWithStatus(data, 0, sync...)
}

// 写入数据，自定义status字段
func (blog *BinLog) AppendWithStatus(data []byte, status int, sync...bool) (int64, error) {
    // 判断是否需要新生成log文件
    if blog.tail.size.Val() >= blog.cap.Val() {
        blog.mu.Lock()
        blog.tail = blog.getLog(blog.tail.num + 1)
        blog.mu.Unlock()
    }
    // 执行数据写入
    length  := len(data)
    txidbuf := gbinary.EncodeInt64(gtime.Nanosecond())
    buffer  := make([]byte, len(data) + 21)
    copy(buffer[0:],           gbinary.EncodeInt8(int8(status)))
    copy(buffer[1:],           gbinary.EncodeInt32(int32(length)))
    copy(buffer[5:],           txidbuf)
    copy(buffer[13:],          data)
    copy(buffer[13 + length:], txidbuf)
    return blog.tail.append(buffer, sync...)
}

func (blog *BinLog) Pop() (*LogItem, error) {


}

func (blog *BinLog) GetById(id int64) (*LogItem, error) {
    num     := int(id & 0xffff)
    offset  := id >> 16
    for {
        log := blog.getLog(num)
        if log.size.Val() == 0 {
            break
        }
        pf, err := log.fp.File()
        if err != nil {
            return nil, err
        }
        defer pf.Close()
        for {
            buffer := gfile.GetBinContentByTwoOffsets(pf.File(), offset, gDEFAULT_BLOCK_SIZE)
            if buffer == nil {
                break
            }
            length   := gbinary.DecodeToInt32(buffer[1 : 5])
            itemsize := int(length + 21)
            if itemsize > len(buffer) {
                buffer = append(buffer, gfile.GetBinContentByTwoOffsets(pf.File(),
                    offset + int64(len(buffer)), int64(itemsize - len(buffer)))...)
                if itemsize != len(buffer) {
                    return nil, errors.New(fmt.Sprintf("invalid log id: %d", id))
                }
            }
            return &LogItem{

            }

        }

    }
    return nil, errors.New(fmt.Sprintf("data not found for id: %d", id))
}

func bufferToLogItem(buffer []byte) *LogItem {

}

// 根据文件偏移量计算id
func (log *Log) idByOffset(offset int64) int64 {
    return offset << 16 | int64(log.num)
}

// 写入数据
func (log *Log) append(buffer []byte, sync...bool) (int64, error) {
    pf, err := log.fp.File()
    if err != nil {
        return 0, err
    }
    defer pf.Close()

    log.mu.Lock()
    defer log.mu.Unlock()

    // 写到文件末尾
    start, err := pf.File().Seek(0, 2)
    if err != nil {
        return 0, err
    }
    // 执行数据写入
    if _, err := pf.File().WriteAt(buffer, start); err != nil {
        return 0, err
    }
    // 是否执行文件sync
    if len(sync) > 0 && sync[0] {
        if err := pf.File().Sync(); err != nil {
            return 0, err
        }
    }
    size := int(start) + len(buffer)
    log.size.Set(size)
    return log.idByOffset(int64(size)), nil
}
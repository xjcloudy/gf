// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// KV嵌入式数据库，底层使用的是LevelDB.
package gkv

import "github.com/syndtr/goleveldb/leveldb"

type DB struct {
    db *leveldb.DB
}

// 创建/读取一个KV数据库
func New(path string) (*DB, error) {
    if db, err := leveldb.OpenFile(path, nil); err != nil {
        return nil, err
    } else {
        return &DB{
            db : db,
        }, nil
    }
}

func(db *DB) Set(key, value []byte) error {
    return db.db.Put(key, value, nil)
}

func(db *DB) Get(key []byte) []byte {
    value, _ := db.db.Get(key, nil)
    return value
}

func(db *DB) Delete(key []byte) error {
    return db.db.Delete(key, nil)
}
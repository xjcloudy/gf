// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 字符编码转换.
package gcharset

import (
    "github.com/axgle/mahonia"
)

func Convert(content, toEncoding, fromEncoding string) string {
    return mahonia.NewEncoder(fromEncoding).ConvertString(content)
}
package gdrh



// DRH算法操作对象
type Map struct {
    size   int     // 分区基数
    degree int     // 分区阶数
    root   *table  // 根哈希表
}

// 哈希表
type table struct {
    m      *Map    // 所属Map对象
    deep   int     // 表深度
    size   int     // 表分区(除根节点外，必须为奇数)
    parts  []*part // 分区数组
}

// 哈希表中的分区
type part struct {
    items  []*item // 数据项列表，必须按照值进行从小到大排序，便于二分查找
}

// 分区中的数据项
type item struct {
    key   int   // 数据项键名，这里设置为int，便于演示
    value int   // 数据项键值，这里设置为int，便于演示
    part  *part // 深度分区标识，指向另外一个分区
}

// 创建DRH对象
func New(size, degree int) *Map {
    m := &Map {
        size   : size,
        degree : degree,
    }
    m.root = &table {
        m     : m,
        deep  : 0,
        size  : size,
        parts : make([]*part, size),
    }
    return m
}

// 设置键值对数据
func (m *Map) Set(key, value int) {
    m.root.set(key, value)
}

func (t *table) set(key, value int) {
    part  := t.parts[key % len(t.parts)]
    if part == nil {
        part = &part {
            items:[]*item{{
                key   : key,
                value : value,
            }},
        }
    } else {
        index, cmp := part.search(key)
        if cmp == 0 {
            part.items[index].value = value
        } else {
            // 首先进行进行数据项插入
            part.save(&item{key : key, value : value }, index, cmp)
            // 接着再判断是否需要进行DRH算法处理
            if len(part.items) == t.m.degree {

            }
        }
    }
}

// 添加一项, cmp < 0往前插入，cmp >= 0往后插入
func (p *part) save(item *item, index int, cmp int) {
    if cmp == 0 {
        p.items[index] = item
    }
    pos := index
    if cmp == -1 {
        // 添加到前面
    } else {
        // 添加到后面
        pos = index + 1
        if pos >= len(p.items) {
            pos = len(p.items)
        }
    }
    rear   := append([]*item{}, p.items[pos : ]...)
    p.items = append(p.items[0 : pos], item)
    p.items = append(p.items, rear...)
}

// 返回值1: 二分查找中最近对比的数组索引
// 返回值2: -2表示压根什么都未找到，-1表示最近一个索引对应的值比key小，1最近一个索引对应的值比key大
// 两个值的好处是即使匹配不到key,也能进一步确定插入的位置索引
func (p *part) search(key int) (int, int) {
    min := 0
    max := len(p.items) - 1
    mid := 0
    cmp := -2
    for min < max {
        mid = int((min + max) / 2)
        if key < p.items[mid].key {
            max = mid - 1
            cmp = -1
        } else if key > p.items[mid].key {
            min = mid + 1
            cmp = 1
        } else {
            return mid, 0
        }
    }
    return mid, cmp
}
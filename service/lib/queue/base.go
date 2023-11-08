package queue

// 队列器
type Queuer interface {
	// 左侧插入
	LPush(value ...interface{}) error
	// 右侧插入
	RPush(value ...interface{}) error
	// 删除元素
	Delete(value interface{}) error
	// 使用下标获取值
	GetByIndex(index int64, v interface{}) error
	// 左侧读取并删除
	LPop(v interface{}) error
	// 右侧读取并删除
	RPop(v interface{}) error
	// 队列长度
	Length() (int64, error)
	// 清空队列
	Flush() error
}

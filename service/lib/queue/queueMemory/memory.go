package queueMemory

import (
	"encoding/json"
	"errors"
	"reflect"
	"sync"
)

type Pool struct {
	Values [][]byte
	Lock   sync.RWMutex
}

func New() *Pool {
	return &Pool{}
}

func (k *Pool) LPush(value ...interface{}) error {
	k.Lock.Lock()
	defer k.Lock.Unlock()
	for i := 0; i < len(value); i++ {
		v, _ := json.Marshal(value[i])
		k.Values = append([][]byte{v}, k.Values...)
	}
	return nil
}

func (k *Pool) RPush(value ...interface{}) error {
	k.Lock.Lock()
	defer k.Lock.Unlock()
	for i := 0; i < len(value); i++ {
		v, _ := json.Marshal(value[i])
		k.Values = append(k.Values, v)
	}
	return nil
}

func (k *Pool) Delete(value interface{}) error {
	k.Lock.Lock()
	defer k.Lock.Unlock()
	var index int64
	v, _ := json.Marshal(value)

	for i, item := range k.Values {
		if reflect.DeepEqual(item, v) {
			index = int64(i)
			k.removeIndex(index)
			return nil
		}
	}

	return nil
}

// 取出值
func (k *Pool) GetByIndex(index int64, v interface{}) error {
	k.Lock.RLock()
	defer k.Lock.RUnlock()
	if int64(len(k.Values)) >= index {
		json.Unmarshal(k.Values[index], v)
		return nil
	} else {
		return errors.New("index non-existent")
	}
}

// 左-取出并删除
func (k *Pool) LPop(v interface{}) error {
	if err := k.GetByIndex(0, v); err != nil {
		return err
	} else {
		k.removeIndex(0)
	}
	return nil
}

// 右-取出并删除
func (k *Pool) RPop(v interface{}) error {
	index := int64(len(k.Values) - 1)
	if err := k.GetByIndex(index, v); err != nil {
		return err
	} else {
		k.removeIndex(index)
	}
	return nil
}

func (k *Pool) removeIndex(index int64) error {
	k.Lock.Lock()
	defer k.Lock.Unlock()
	k.Values = append(k.Values[:index], k.Values[index+1:]...)
	return nil
}

func (k *Pool) Length() (int64, error) {
	return int64(len(k.Values)), nil
}

func (k *Pool) Flush() error {
	k.Lock.Lock()
	defer k.Lock.Unlock()
	k.Values = nil
	return nil
}

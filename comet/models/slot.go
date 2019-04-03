package models

import (
	"sync"
)

/**
 * 在线信息节点
 * 关于空结构体的内存占用问题@see http://www.golangtc.com/t/575442b8b09ecc02f7000057
 * @link http://misfra.me/optimizing-concurrent-map-access-in-go
 */
type Slot struct {
	lock sync.RWMutex
	dict map[string]*Session
}

/**
 * 初始化
 */
func NewSlot(size int) *Slot{
	s := &Slot{}
	s.dict = make(map[string]*Session, size)
	s.lock = sync.RWMutex{}
	return s
}

/**
 * 整个槽的容量
 */
func (slot *Slot) Len() int {
	return len(slot.dict)
}

/**
 * 是否包含某个uuid
 */
func (slot *Slot) Has(uuid string) bool {
	slot.lock.RLock()
	defer slot.lock.RUnlock()
	_, ok := slot.dict[uuid]
	return ok
}

/**
 * 获取某个Session信息
 *
 * @param string udid
 * @return *Session
 */
func (slot *Slot) Get(udid string) *Session {
	slot.lock.RLock()
	defer slot.lock.RUnlock()
	if m, ok := slot.dict[udid]; ok {
		return m
	}
	return nil
}

/**
 * 添加Session
 *
 * @link http://dave.cheney.net/2014/03/25/the-empty-struct
 */
func (slot *Slot) Add(udid string, session *Session) {
	slot.lock.Lock()
	defer slot.lock.Unlock()
	if _, ok := slot.dict[udid]; !ok {
		slot.dict[udid] = session
	}
	slot.dict[udid] = session
}

/**
 * 删除Session
 *
 * @param string udid
 * @return string
 */
func (slot *Slot) Del(udid string) {
	slot.lock.Lock()
	defer slot.lock.Unlock()
	if _, ok := slot.dict[udid]; ok {
		delete(slot.dict, udid)
	}
}

/**
 * 获取某个槽下的全部Session
 */
func (slot *Slot) All() map[string]*Session {
	return slot.dict
}

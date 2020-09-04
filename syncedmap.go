package maps

import (
	"sync"
)

type SyncedMap struct {
	Data Map
	mu   sync.RWMutex
}

func (this *SyncedMap) Set(key interface{}, value interface{}) *SyncedMap {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.Data.Set(key, value)
	return this
}

func (this *SyncedMap) Update(values ...map[string]interface{}) *SyncedMap {
	this.mu.Lock()
	defer this.mu.Unlock()

	for _, value := range values {
		for k, v := range value {
			this.Data.Set(k, v)
		}
	}
	return this
}

func (this *SyncedMap) Del(key ...interface{}) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.Data.Del(key...)
}

func (this *SyncedMap) GetDel(key interface{}) (value interface{}, ok bool) {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.Data.GetDel(key)
}

func (this SyncedMap) Has(key ...interface{}) bool {
	this.mu.RLock()
	defer this.mu.RUnlock()
	return this.Data.Has(key...)
}

func (this SyncedMap) Get(key interface{}) (value interface{}, ok bool) {
	if this.Data == nil {
		return
	}
	this.mu.RLock()
	defer this.mu.RUnlock()
	if this.Data == nil {
		return
	}
	value, ok = this.Data[key]
	return
}

func (this SyncedMap) GetBool(key interface{}, defaul ...bool) bool {
	v, _ := this.Get(key)
	if v != nil {
		return v.(bool)
	}
	if len(defaul) > 0 {
		return defaul[0]
	}
	return false
}

func (this SyncedMap) GetString(key interface{}, defaul ...string) string {
	v, _ := this.Get(key)
	if v != nil {
		return v.(string)
	}
	if len(defaul) > 0 {
		return defaul[0]
	}
	return ""
}

func (this SyncedMap) GetInt(key interface{}, defaul ...int) int {
	v, _ := this.Get(key)
	if v != nil {
		return v.(int)
	}
	if len(defaul) > 0 {
		return defaul[0]
	}
	return 0
}

func (this SyncedMap) GetSlice(key interface{}, defaul ...[]interface{}) (r []interface{}) {
	v, _ := this.Get(key)
	if v != nil {
		return v.([]interface{})
	}
	if len(defaul) > 0 {
		return defaul[0]
	}
	return
}

func (this SyncedMap) GetStrings(key interface{}, defaul ...[]string) (r []string) {
	v, _ := this.Get(key)
	if v != nil {
		return v.([]string)
	}
	if len(defaul) > 0 {
		return defaul[0]
	}
	return
}

func (this SyncedMap) MustGetInterface(key interface{}, defaul ...interface{}) interface{} {
	if v, ok := this.Get(key); ok {
		return v
	}
	if len(defaul) > 0 {
		return defaul[0]
	}
	return nil
}

func (this SyncedMap) GetInterface(key string, defaul ...interface{}) interface{} {
	if v, ok := this.Get(key); ok {
		return v
	}
	if len(defaul) > 0 {
		return defaul[0]
	}
	return nil
}

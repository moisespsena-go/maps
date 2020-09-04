package maps

import (
	"bytes"
	"database/sql/driver"
	"regexp"

	"github.com/moisespsena-go/aorm"

	"github.com/moisespsena-go/getters"
	"gopkg.in/yaml.v2"
)

type MapSiSlice []MapSI

func (this *MapSiSlice) Append(m ...MapSI) {
	*this = append(*this, m...)
}

func (this *MapSiSlice) AppendMap(m ...map[string]interface{}) {
	for _, m := range m {
		*this = append(*this, m)
	}
}

func (this MapSiSlice) Get(key string) (value interface{}, ok bool) {
	l := len(this)
	for i := l; i > 0; i-- {
		if value, ok = this[i-1].Get(key); ok {
			return
		}
	}
	return
}

func (this MapSiSlice) GetMapS(key string) (value MapSI, ok bool) {
	l := len(this)
	for i := l; i > 0; i-- {
		if value, ok = this[i-1].GetMapS(key); ok {
			return
		}
	}
	return
}

func (this MapSiSlice) GetMap(key string) (value Map, ok bool) {
	l := len(this)
	for i := l; i > 0; i-- {
		if value, ok = this[i-1].GetMap(key); ok {
			return
		}
	}
	return
}

func (this MapSiSlice) Has(key string) (ok bool) {
	l := len(this)
	for i := l; i > 0; i-- {
		if this[i-1].Has(key) {
			return true
		}
	}
	return
}

type MapSI map[string]interface{}

func (this MapSI) AormDataType(dialect aorm.Dialector) string {
	if ok, _ := regexp.MatchString(`postgres`, dialect.GetName()); ok {
		return "JSONB"
	}
	return "TEXT"
}

func (this MapSI) Value() (driver.Value, error) {
	if this == nil || len(this) == 0 {
		return nil, nil
	}
	return yaml.Marshal(this)
}

func (this *MapSI) Scan(src interface{}) (err error) {
	if src == nil {
		*this = nil
		return
	}
	switch t := src.(type) {
	case string:
		return this.Scan([]byte(t))
	case []byte:
		*this = MapSI{}
		return yaml.NewDecoder(bytes.NewBuffer(t)).Decode(this)
	}
	return
}

func (this *MapSI) Set(key string, value interface{}) *MapSI {
	if *this == nil {
		*this = map[string]interface{}{
			key: value,
		}
	} else {
		(*this)[key] = value
	}
	return this
}

func (this *MapSI) Del(key ...string) {
	if this == nil {
		return
	}

	for _, key := range key {
		delete(*this, key)
	}
	return
}

func (this *MapSI) GetDel(key string) (value interface{}, ok bool) {
	if this == nil {
		return
	}

	if value, ok = (*this)[key]; ok {
		delete(*this, key)
	}
	return
}

func (this MapSI) Has(key ...string) (ok bool) {
	if this == nil {
		return
	}

	for _, key := range key {
		if _, ok = (this)[key]; !ok {
			return false
		}
	}
	return
}

func (this *MapSI) HasOne(key ...string) (ok bool) {
	if *this == nil {
		return
	}

	for _, key := range key {
		if _, ok = (*this)[key]; ok {
			return
		}
	}
	return
}

func (this *MapSI) Update(m_ ...map[string]interface{}) *MapSI {
	for _, m_ := range m_ {
		for k, v := range m_ {
			this.Set(k, v)
		}
	}
	return this
}

func (this MapSI) Get(key string) (value interface{}, ok bool) {
	if this == nil {
		return
	}
	value, ok = (this)[key]
	return
}

func (this MapSI) Find(pth ...string) (value interface{}, ok bool) {
	if this == nil {
		return
	}
	var m interface{} = this
	for _, key := range pth {
		switch mt := m.(type) {
		case interface {
			Get(string) (interface{}, bool)
		}:
			if m, ok = mt.Get(key); !ok {
				return
			}
		case interface {
			Get(interface{}) (interface{}, bool)
		}:
			if m, ok = mt.Get(key); !ok {
				return
			}
		case map[string]interface{}:
			if m, ok = mt[key]; !ok {
				return
			}
		case map[interface{}]interface{}:
			if m, ok = mt[key]; !ok {
				return
			}
		}
	}
	return m, m != nil
}

func (this *MapSI) GetMapS(key string) (value MapSI, ok bool) {
	if *this == nil {
		return
	}
	var vi interface{}
	if vi, ok = (*this)[key]; ok {
		switch t := vi.(type) {
		case MapSI:
			return t, true
		case map[string]interface{}:
			return MapSI(t), true
		default:
			return
		}
	}
	return
}

func (this *MapSI) GetMap(key string) (value Map, ok bool) {
	if *this == nil {
		return
	}
	var vi interface{}
	if vi, ok = (*this)[key]; ok {
		return Map(vi.(map[interface{}]interface{})), true
	}
	return
}

func (this *MapSI) Getter() getters.Getter {
	return getters.New(func(key interface{}) (value interface{}, ok bool) {
		switch kt := key.(type) {
		case string:
			return this.Get(kt)
		case []string:
			return this.Find(kt...)
		}
		return
	})
}

func (this *MapSI) ReadKey(getter getters.Getter, key interface{}) (ok bool) {
	var v interface{}
	if v, ok = getter.Get(key); ok {
		*this = v.(MapSI)
	}
	return
}

func GetMapSI(getter getters.Getter, key string) (m MapSI, ok bool) {
	if v, ok := getter.Get(key); ok {
		return v.(MapSI), true
	}
	return 
}
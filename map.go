package maps

type Map map[interface{}]interface{}

func (this *Map) Set(key interface{}, value interface{}) *Map {
	if *this == nil {
		*this = map[interface{}]interface{}{
			key: value,
		}
	} else {
		(*this)[key] = value
	}
	return this
}

func (this *Map) Del(key ...interface{}) {
	if this == nil {
		return
	}

	for _, key := range key {
		delete(*this, key)
	}
	return
}

func (this *Map) GetDel(key interface{}) (value interface{}, ok bool) {
	if this == nil {
		return
	}

	if value, ok = (*this)[key]; ok {
		delete(*this, key)
	}
	return
}

func (this *Map) Has(key ...interface{}) (ok bool) {
	if *this == nil {
		return
	}

	for _, key := range key {
		if _, ok = (*this)[key]; !ok {
			return false
		}
	}
	return
}

func (this *Map) HasOne(key ...interface{}) (ok bool) {
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

func (this *Map) Update(m_ ...map[interface{}]interface{}) *Map {
	for _, m_ := range m_ {
		for k, v := range m_ {
			this.Set(k, v)
		}
	}
	return this
}

func (this Map) Get(key interface{}) (value interface{}, ok bool) {
	if this == nil {
		return
	}
	value, ok = this[key]
	return
}

func (this Map) Value(key interface{}) (value interface{}) {
	if this == nil {
		return nil
	}
	return this[key]
}

func (this Map) GetString(key interface{}, defaul ...string) (value string) {
	if this == nil {
		for _, value = range defaul {
			return
		}
		return
	}
	if v, ok := this[key]; ok {
		value = v.(string)
	} else {
		for _, value = range defaul {
			return
		}
	}
	return
}

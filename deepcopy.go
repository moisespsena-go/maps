package maps

import "errors"

// DeepCopyTo deep copy the config into a dest.
func (this MapSI) DeepCopy(dest interface{}, cb ...func(m MapSI, key string, value interface{}) (interface{}, error)) (err error) {
	if this == nil {
		return errors.New("this::MapSI is nil. You cannot read from a nil map")
	}
	if dest == nil {
		return errors.New("dest is nil. You cannot insert to a nil map")
	}

	var cp = func(m MapSI, k string, v interface{}) (err error) {
		for _, cb := range cb {
			if v, err = cb(m, k, v); err != nil {
				return
			}
		}
		switch dt := dest.(type) {
		case MapSI:
			dt[k] = v
		case map[string]interface{}:
			dt[k] = v
		case map[interface{}]interface{}:
			dt[k] = v
		}
		return
	}

	var deepCopy = func(src MapSI, k string) (err error) {
		var destSM MapSI
		switch destType := dest.(type) {
		case map[string]interface{}:
			destSM = destType
		case MapSI:
			destSM = destType
		case map[interface{}]interface{}:
			if _, ok := destType[k]; ok {
				switch t := destType[k].(type) {
				case MapSI:
					return src.DeepCopy(t, cb...)
				case map[interface{}]interface{}:
					return src.DeepCopy(t, cb...)
				case map[string]interface{}:
				}
			}
			m := make(MapSI)
			if err = src.DeepCopy(m, cb...); err == nil {
				destType[k] = m
			}
			return
		}

		if _, ok := destSM[k]; ok {
			switch t := destSM[k].(type) {
			case MapSI:
				return src.DeepCopy(t, cb...)
			case map[string]interface{}:
				return src.DeepCopy(t, cb...)
			case map[interface{}]interface{}:
				return src.DeepCopy(t, cb...)
			default:
				m := make(MapSI)
				if err = src.DeepCopy(m, cb...); err == nil {
					destSM[k] = m
				}
			}
		} else {
			m := make(MapSI)
			if err = src.DeepCopy(m, cb...); err == nil {
				destSM[k] = m
			}
		}
		return
	}

	for k, v := range this {
		switch vt := v.(type) {
		case MapSI:
			if err = deepCopy(vt, k); err != nil {
				return
			}
		case map[string]interface{}:
			if err = deepCopy(vt, k); err != nil {
				return
			}
		case map[interface{}]interface{}:
			var vSI = make(MapSI)
			for vk, vv := range vt {
				if s, ok := vk.(string); ok {
					vSI[s] = vv
				}
			}
			if err = deepCopy(vSI, k); err != nil {
				return
			}
		default:
			if err = cp(this, k, v); err != nil {
				return
			}
		}
	}
	return
}

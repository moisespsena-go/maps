package maps

import "sync"

func Diff(a, b map[string]interface{}, sameA bool) (add, del, same SyncedMap) {
	var w sync.WaitGroup
	w.Add(2)
	go func() {
		defer w.Done()
		if sameA {
			for key, value := range a {
				if _, ok := b[key]; ok {
					same.Set(key, value)
				} else {
					del.Set(key, value)
				}
			}
		} else {
			for key, value := range a {
				if _, ok := b[key]; !ok {
					del.Set(key, value)
				}
			}
		}
	}()
	go func() {
		defer w.Done()
		if sameA {
			for key, value := range b {
				if _, ok := a[key]; !ok {
					add.Set(key, value)
				}
			}
		} else {
			for key, value := range b {
				if _, ok := a[key]; ok {
					same.Set(key, value)
				} else {
					add.Set(key, value)
				}
			}
		}
	}()
	w.Wait()
	return
}

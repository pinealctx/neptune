package async

import (
	"context"
	"reflect"
	"sync"
)

//type validation cache
var (
	typeCache     = make(map[reflect.Type]struct{})
	typeCacheLock sync.RWMutex
	//error interface
	errorInterface = reflect.TypeOf((*error)(nil)).Elem()
	//context interface
	contextInterface = reflect.TypeOf((*context.Context)(nil)).Elem()
)

//add type validation
func addType2Validation(t reflect.Type) {
	typeCacheLock.Lock()
	typeCache[t] = struct{}{}
	typeCacheLock.Unlock()
}

//is in type validation cache
func inValidateCache(t reflect.Type) bool {
	var ok bool
	typeCacheLock.RLock()
	_, ok = typeCache[t]
	typeCacheLock.RUnlock()
	return ok
}

//validate function
func validateFn(fn interface{}) (reflect.Type, bool) {
	if fn == nil {
		return nil, false
	}
	var fnType = reflect.TypeOf(fn)
	if inValidateCache(fnType) {
		return fnType, true
	}
	if fnType.Kind() != reflect.Func {
		return fnType, false
	}
	if fnType.NumIn() != 2 {
		return fnType, false
	}
	if fnType.NumOut() != 2 {
		return fnType, false
	}
	if !fnType.In(0).Implements(contextInterface) {
		return fnType, false
	}
	if !fnType.Out(1).Implements(errorInterface) {
		return fnType, false
	}
	addType2Validation(fnType)
	return fnType, true
}

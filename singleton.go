package singlet

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrTypeMismatch = errors.New("the requested type does not match the singleton type")

// Singleton can store a single value atomically when used with GetOrDo.
type Singleton struct {
	p   atomic.Pointer[any]
	mtx sync.Mutex
}

// GetOrDo will create and store a value in the Singleton, if one doesn't already exist, by calling the fn, else it will return the existing value.
// If the fn returns an error, that error will be returned and no value will be stored in the singleton.
// If the requested type does not match a type the existing singleton type, ErrTypeMismatch is returned.
// This function is threadsafe.
func GetOrDo[T any](singleton *Singleton, fn func() (T, error)) (result T, err error) {
	maybeResult := singleton.p.Load()
	if maybeResult == nil {
		// Lock to guard against applying fn twice
		singleton.mtx.Lock()
		defer singleton.mtx.Unlock()
		maybeResult = singleton.p.Load()

		// Double check
		if maybeResult == nil {
			result, err = fn()
			if err == nil {
				var resultAny any = result
				singleton.p.Store(&resultAny)
			}
			return result, err
		}
	}

	var ok bool
	result, ok = (*maybeResult).(T)
	if !ok {
		return *new(T), ErrTypeMismatch
	}
	return result, nil
}

// Get returns a previously created value for the singleton, else the default value for T.
// Returns ErrTypeMismatch if the requested type does not match a type the existing singleton type.
// This function is threadsafe.
func Get[T any](singleton *Singleton) (T, error) {
	maybeResult := singleton.p.Load()
	if maybeResult == nil {
		return *new(T), nil
	}
	result, ok := (*maybeResult).(T)
	if !ok {
		return *new(T), ErrTypeMismatch
	}
	return result, nil
}

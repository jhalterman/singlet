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

// Get returns a previously created value for the singleton, else the default value for T. This function is threadsafe.
// Returns ErrTypeMismatch if the requested type does not match a type that was previously stored in the singleton.
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

// GetOrDo returns a previously created value for the singleton, else it creates and returns a new one by calling fn. This function is threadsafe.
// Returns ErrTypeMismatch if the requested type does not match a type that was previously stored in the singleton.
func GetOrDo[T any](singleton *Singleton, fn func() T) (result T, err error) {
	maybeResult := singleton.p.Load()
	if maybeResult == nil {
		// Lock to guard against applying fn twice
		singleton.mtx.Lock()
		defer singleton.mtx.Unlock()
		maybeResult = singleton.p.Load()

		// Double check
		if maybeResult == nil {
			result = fn()
			var resultAny any = result
			singleton.p.Store(&resultAny)
			return result, nil
		}
	}

	var ok bool
	result, ok = (*maybeResult).(T)
	if !ok {
		return *new(T), ErrTypeMismatch
	}
	return result, nil
}

package singlet

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Foo[T any] struct {
	t T
}

var (
	fn1 = func() Foo[int] {
		return Foo[int]{t: 1}
	}
	fn2 = func() Foo[int] {
		return Foo[int]{t: 2}
	}
	fn1Ptr = func() *Foo[int] {
		return &Foo[int]{t: 1}
	}
	fn2Ptr = func() *Foo[int] {
		return &Foo[int]{t: 2}
	}
)

func TestGetOrDo(t *testing.T) {
	t.Run("should return same value for T", func(t *testing.T) {
		// Asseert fn1 and fn2 produce different results
		assert.NotEqual(t, fn1(), fn2())

		// Assert GetOrDo(s, fn1) and GetOrDo(s, fn2) produce the same result
		s := &Singleton{}
		foo1, err1 := GetOrDo(s, fn1)
		foo2, err2 := GetOrDo(s, fn2)
		assert.Equal(t, foo1, foo2)
		assert.NoError(t, err1)
		assert.NoError(t, err2)
	})

	t.Run("should return same value for *T", func(t *testing.T) {
		// Asseert fn1 and fn2 produce different results
		assert.NotEqual(t, fn1Ptr(), fn2Ptr())

		// Assert GetOrDo(s, fn1) and GetOrDo(s, fn2) produce the same result
		s := &Singleton{}
		foo1, err1 := GetOrDo(s, fn1Ptr)
		foo2, err2 := GetOrDo(s, fn2Ptr)
		assert.Equal(t, foo1, foo2)
		assert.NoError(t, err1)
		assert.NoError(t, err2)
	})

	t.Run("should error when requested type doesn't match singleton type", func(t *testing.T) {
		s := &Singleton{}
		GetOrDo(s, fn1)
		_, err := GetOrDo(s, func() string {
			return "wrong type"
		})
		assert.Error(t, ErrTypeMismatch, err)
	})
}

func TestGet(t *testing.T) {
	t.Run("should return stored value", func(t *testing.T) {
		singleton := &Singleton{}

		result, err := Get[Foo[int]](singleton)
		assert.NoError(t, err)
		assert.Equal(t, Foo[int]{}, result)

		GetOrDo(singleton, fn1)
		result, err = Get[Foo[int]](singleton)
		assert.NoError(t, err)
		assert.Equal(t, Foo[int]{t: 1}, result)
	})

	t.Run("should error when requested type doesn't match singleton type", func(t *testing.T) {
		s := &Singleton{}
		GetOrDo(s, fn1)
		_, err := Get[string](s)
		assert.Error(t, ErrTypeMismatch, err)
	})
}

func BenchmarkGetOrDo(b *testing.B) {
	b.Run("sync.Once", func(b *testing.B) {
		once := sync.Once{}
		for i := 0; i < b.N; i++ {
			once.Do(func() {
				_ = Foo[int]{t: 1}
			})
		}
	})

	b.Run("singleton.GetOrDo", func(b *testing.B) {
		singleton := &Singleton{}
		for i := 0; i < b.N; i++ {
			_, _ = GetOrDo(singleton, func() Foo[int] {
				return Foo[int]{t: 1}
			})
		}
	})
}

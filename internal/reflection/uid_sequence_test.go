package reflection

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_UidSequence_Next_Current_returns_initial_value_plus_one(t *testing.T) {

	// Arrange
	initialValue := uint64(0)
	sut := NewUidSequence(WithInitialValue(initialValue))

	// Act
	actual := sut.Next()

	// Assert
	assert.Equal(t, uint64(initialValue+1), actual)
}

func Test_UidSequence_Next_concurrent(t *testing.T) {

	// Arrange
	sut := NewUidSequence()

	const delta = 100

	wg := sync.WaitGroup{}
	identifiers := make(chan uint64, delta)

	wg.Add(delta)

	go func() {
		for i := 0; i < delta; i++ {
			go func() {
				defer wg.Done()
				duration := time.Duration(5 + rand.Intn(45))
				time.Sleep(duration * time.Millisecond)
				identifiers <- sut.Next()
			}()
		}
	}()

	// Act
	wg.Wait()

	// Assert
	close(identifiers)

	seen := make(map[uint64]struct{})
	for id := range identifiers {
		_, found := seen[id]
		seen[id] = struct{}{}
		assert.False(t, found)
	}
}

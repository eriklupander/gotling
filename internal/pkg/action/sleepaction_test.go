package action

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSleepAction_ExecuteWithDurationFormat(t *testing.T) {
	action := NewSleepAction(map[interface{}]interface{}{"duration":"100ms"})
	start := time.Now()
	action.Execute(nil, nil)
	assert.Greater(t, time.Since(start).Milliseconds(), int64(99))
}

func TestSleepAction_ExecuteWithIntegerFormat(t *testing.T) {
	action := NewSleepAction(map[interface{}]interface{}{"duration":1})
	start := time.Now()
	action.Execute(nil, nil)
	assert.Greater(t, time.Since(start).Milliseconds(), int64(999))
}
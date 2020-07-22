package incrementor //nolint: testpackage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const maxLoops int64 = 42 // NOTE(toby3d): максимальное число итераций в тестах

func TestGetNumber(t *testing.T) {
	i := New()
	assert.Zero(t, i.GetNumber())

	for j := int64(1); j < maxLoops; j++ {
		i.IncrementNumber()

		if !assert.Equal(t, j, i.GetNumber(), "number must be equal to loop step") {
			t.FailNow()
		}
	}
}

func TestIncrementNumber(t *testing.T) {
	i := New()

	for j := int64(1); j < maxLoops; j++ {
		i.IncrementNumber()

		if !assert.Equal(t, j, i.number, "number must be equal to loop step index") {
			t.FailNow()
		}
	}

	i.SetMaximumValue(maxLoops)
	i.IncrementNumber()
	assert.Zero(t, i.number, "number should be reset when the maximum is reached")
}

func TestSetMaximumValue(t *testing.T) {
	for _, tc := range []struct {
		name        string // NOTE(toby3d): имя кейса
		inputMax    int64  // NOTE(toby3d): устанавливаемый порог
		expMax      int64  // NOTE(toby3d): ожидаемое порог
		inputNumber int64  // NOTE(toby3d): изначальное значение числа
		expNumber   int64  // NOTE(toby3d): ожидаемое значение числа
	}{
		{
			name:     "zero max to default",
			inputMax: 0,
			expMax:   DefaultMaximumValue,
		}, {
			name:     "negative max to default",
			inputMax: -24,
			expMax:   DefaultMaximumValue,
		}, {
			name:     "valid max input",
			inputMax: 42,
			expMax:   42,
		}, {
			name:        "number below new max",
			inputMax:    42,
			inputNumber: 24,
			expMax:      42,
			expNumber:   24,
		}, {
			name:        "number above new max",
			expNumber:   0,
			inputNumber: 420,
			expMax:      42,
			inputMax:    42,
		},
		// NOTE(toby3d): число больше чем DefaultMaximumValue невозможно установить, так как оно приведёт к
		// ошибке компиляции: "constant 9223372036854775808 overflows int"
	} {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			i := New()
			i.number = tc.inputNumber

			i.SetMaximumValue(tc.inputMax)

			if !assert.Equal(t, i.maxValue, tc.expMax) || !assert.Equal(t, i.number, tc.expNumber) {
				t.FailNow()
			}
		})
	}
}

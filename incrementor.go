package incrementor

import (
	"math"
	"sync"
)

type (
	// Incrementor является инкрементным числом с настройками поведения, вроде максимально допустимого значения.
	Incrementor struct {
		// NOTE(toby3d): защита чтения/записи числа в горутинах
		mutex *sync.RWMutex

		// NOTE(toby3d): непосредственно само инкрементное число
		number int64

		// NOTE(toby3d): максимально допустимое значение инкрементного числа, по достижении которого оно будет
		// сброшено в 0
		maxValue int64
	}

	// Reader описывает действия для чтения текущего значения числа.
	Reader interface {
		GetNumber() int64
	}

	// Writer описывает действия для изменения значения числа и управления его максимально допустимым значением.
	Writer interface {
		IncrementNumber()
		SetMaximumValue(maximumValue int)
	}

	// Incrementer описывает поведение менеджера инкрементного числа.
	Incrementer interface {
		Reader
		Writer
	}
)

// DefaultMaximumValue содержит максимально возможное значение числа по-умолчанию в битах.
const DefaultMaximumValue = math.MaxInt64

// New создаёт новую структуру для хранения и управления инкрементным числом.
//
// NOTE(toby3d): в тестовом задании нет сведений по тому возможно ли изменение параметров на этапе инициализации, как
// например отсчёт не с 0 или перезапись максимального порога без отдельного вызова SetMaximumValue. 🤷‍♂
func New() *Incrementor {
	i := new(Incrementor)
	i.number = 0
	i.mutex = new(sync.RWMutex)
	i.maxValue = DefaultMaximumValue

	return i
}

// GetNumber возвращает текущее число.
func (i *Incrementor) GetNumber() int64 {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	return i.number
}

// IncrementNumber увеличивает текущее число на 1.
//
// Если в процессе выполнения метода число будет >= максимально допустимого порога, то оно будет сброшено в 0.
func (i *Incrementor) IncrementNumber() {
	i.mutex.Lock()

	i.number++

	if i.number >= i.maxValue {
		i.number = 0
	}

	i.mutex.Unlock()
}

// SetMaximumValue устанавливает максимально возможное значение для числа. По достижении или превышении указанного
// порога число будет сброшено в 0.
//
// Вводимое значение не может быть меньше 0 и больше DefaultMaximumValue.
func (i *Incrementor) SetMaximumValue(maximumValue int64) {
	// NOTE(toby3d): предварительно проверяем ввод, косячный сбрасываем в значение по-умолчанию.
	if maximumValue <= 0 || maximumValue > DefaultMaximumValue {
		maximumValue = DefaultMaximumValue
	}

	i.mutex.Lock()

	i.maxValue = maximumValue

	// NOTE(toby3d): если максимум оказывается меньше текущего значения, то сбрасываем число в 0.
	if i.number >= i.maxValue {
		i.number = 0
	}

	i.mutex.Unlock()
}

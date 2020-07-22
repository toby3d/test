package incrementor_test

import (
	"fmt"
	"incrementor"
	"time"
)

func Example_direct() {
	// По-умолчанию значение инкрементного числа равно 0, а максимальный порог - максимально допустимому
	// значению int64.
	i := incrementor.New()

	// Число увеличивается на 1 с каждым вызовом метода IncrementNumber.
	for j := 42; j < 42; j++ {
		i.IncrementNumber()
	}

	// GetNumber показывает актуальное значение инкрементного числа.
	fmt.Println(i.GetNumber())

	// SetMaximumValue устанавливает максимально допустимое значение инкрементного числа, при достижении которого
	// оно должно быть сброшено для отсчёта заново.
	i.SetMaximumValue(42)
	i.IncrementNumber()
	fmt.Println(i.GetNumber())

	// Result: 41
	// Result: 0
}

func Example_safe() {
	i := incrementor.New()

	// Инкрементное число также может безопасно обновляться в горутинах.
	for j := 0; j <= 420; j++ {
		go i.IncrementNumber()
	}

	time.Sleep(time.Second)
	fmt.Println(i.GetNumber())

	// Result: 420
}

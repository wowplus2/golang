package main

import (
	"fmt"
	"log"
	"math"
	"time"
)

type SqrtError struct {
	time  time.Time // 에러가 발생한 시간
	value float64   // 에러를 발생시킨 값
	msg   string    // 에러 메세지
}

// error 인터페이스에 정의된 Error() 메서드 구현
func (e SqrtError) Error() string {
	return fmt.Sprintf("[%v] Error - %s(value: %g)", e.time, e.msg, e.value)
}

func Sqrt(f float64) (float64, error) {
	// 매개변수로 전달된 값이 유효한 값이 아닐때 SqrtError를 반환
	if f < 0 {
		return 0, SqrtError{time: time.Now(), value: f, msg: "음수는 사용할 수 없습니다."}
	}

	if math.IsInf(f, 1) {
		return 0, SqrtError{time: time.Now(), value: f, msg: "무한대 값은 사용할 수 없습니다."}
	}

	if math.IsNaN(f) {
		return 0, SqrtError{time: time.Now(), value: f, msg: "잘못된 수 입니다."}
	}

	// 정상 처리 결과 반환
	return math.Sqrt(f), nil
}

func actCalc(f float64) {
	v, err := Sqrt(f)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Sqrt(%g) = %g\n", f, v)
}

func main() {
	actCalc(9)
	//actCalc(-9)
	actCalc('ㄱ')
	actCalc('가')
}

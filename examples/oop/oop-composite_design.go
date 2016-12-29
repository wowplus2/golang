package main

import "fmt"

// 자전거의 Part
type Part struct {
	Name string
	Desc string
	NeedsSpare bool
}

type Parts []Part

// Part와 연결된 Spares 메서드
// NeedsSpare가 true인 part들을 반환한다.
func (parts Parts) Spares() (spares Parts) {
	for _, part := range parts {
		if part.NeedsSpare {
			spares = append(spares, part)
		}
	}
	return spares
}

// 자전거 구조체. Class역활을 한다.
// 크기(Size)와 파트(Parts)로 분류할 수 있다.
type Bicycle struct {
	Size string
	Parts
}

// 자전거 타입별 파트를 정의했다.
var (
	RoadBikeParts = Parts{
		{"chain", "10-speed", true},
		{"tire_size", "23", true},
		{"tape_color", "red", true},
	}
	MountainBikeParts = Parts{
		{"chain", "10-speed", true},
		{"tire_size", "21", true},
		{"front_shok", "Manitou", false},
		{"rear_shok", "Fox", true},
	}
	RecumbentBikeParts = Parts{
		{"chain", "9-speed", true},
		{"tire_size", "28", true},
		{"flag", "tall and orange", true},
	}
)


func main() {
	roadBike := Bicycle{Size: "L", Parts: RoadBikeParts}
	mountBike := Bicycle{Size: "L", Parts: MountainBikeParts}
	recumBike := Bicycle{Size: "L", Parts: RecumbentBikeParts}

	roadBike.Spares()
	fmt.Println(roadBike)
	fmt.Println(mountBike)
	fmt.Println(recumBike)
}
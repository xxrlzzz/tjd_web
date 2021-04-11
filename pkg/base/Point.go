package base

import "fmt"

type Point struct {
	Latitude  float64 `json:"latitude" valid:"Required"`
	Longitude float64 `json:"longitude" valid:"Required"`
}

func (p *Point) String() string {
	return fmt.Sprint(p.Latitude, ",", p.Longitude)
}

package types

import "fmt"

// Price 价格,单位厘,0.001元
type Price int64

func (this Price) String() string {
	return fmt.Sprintf("%.2f 元", this.Yuan())
}

func (this Price) Yuan() float64 {
	return float64(this) / 1000
}

func (this Price) Int64() int64 {
	return int64(this)
}

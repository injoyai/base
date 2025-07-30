package types

import "fmt"

// Price 价格,单位厘,0.001元
type Price int64

func (this Price) Int() int64 {
	return int64(this)
}

func (this Price) Float() float64 {
	return float64(this) / 1e3
}

func (this Price) String() string {
	f, unit := this.FloatUnit()
	//f += 0.005,已经四舍五入了
	return fmt.Sprintf("%.2f %s", f, unit)
}

// Yuan 价格,单位元
func (this Price) Yuan() float64 {
	return float64(this) / 1e3
}

// Thousand 价格,单位万元
func (this Price) Thousand() float64 {
	return this.Yuan() / 1e4
}

func (this Price) FloatUnit() (float64, string) {
	switch {
	case this < 1e2:
		return float64(this) / 10, "分"
	case this < 1e3:
		return float64(this) / 1e2, "毛"
	case this < 1e7:
		return float64(this) / 1e3, "元"
	case this < 1e11:
		return float64(this) / 1e7, "万元"
	case this < 1e15:
		return float64(this) / 1e11, "亿元"
	default:
		return float64(this) / 1e15, "万亿元"
	}
}

func Yuan[T Number](f T) Price {
	return Price(int64(f)) * 1000
}

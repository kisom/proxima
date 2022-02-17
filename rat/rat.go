package rat

import (
	"math"
	"math/big"
	"time"
)

var (
	Half   = big.NewRat(1, 2)
	One    = UInt64(1)
	Zero   = UInt64(0)
	K      = UInt64(1000)
	Second = Int64(int64(time.Second))
)

func Rat() *big.Rat {
	return new(big.Rat)
}

func UInt64(v uint64) *big.Rat {
	return Rat().SetUint64(v)
}

func Int64(v int64) *big.Rat {
	return Rat().SetInt64(v)
}

func Float(f float64) *big.Rat {
	return Rat().SetFloat64(f)
}

func AsFloat(x *big.Rat) float64 {
	f, _ := x.Float64()
	return f
}

func FromString(s string) *big.Rat {
	r := Rat()
	r, _ = r.SetString(s)
	return r
}

func Duration(x *big.Rat) time.Duration {
	f := AsFloat(x)
	f = math.Round(f)
	n := int64(f)
	return time.Duration(n)
}

func DurationSeconds(x *big.Rat) time.Duration {
	t := Mul(x, Second)
	return Duration(t)
}

func Mul(x, y *big.Rat) *big.Rat {
	return Rat().Mul(x, y)
}

func Add(x, y *big.Rat) *big.Rat {
	return Rat().Add(x, y)
}

func Div(x, y *big.Rat) *big.Rat {
	return Rat().Quo(x, y)
}

func Sub(x, y *big.Rat) *big.Rat {
	return Rat().Sub(x, y)
}

func Sqr(x *big.Rat) *big.Rat {
	return Mul(x, x)
}

func Sqrt(x *big.Rat) *big.Rat {
	num := new(big.Float).SetInt(x.Num())
	denom := new(big.Float).SetInt(x.Denom())
	num = num.Sqrt(num)
	denom = denom.Sqrt(denom)
	root := new(big.Float).Quo(num, denom)
	r, _ := root.Rat(nil)
	return r
}

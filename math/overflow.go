package math

import (
	builtinMath "math"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

var log = logger.GetOrCreate("arwen/overflow")

// AddUint64 performs addition on uint64 and logs an error if the addition overflows
func AddUint64(a, b uint64) uint64 {
	res, err := AddUint64WithErr(a, b)
	if err != nil {
		log.Error("AddUint64 overflow", "a", a, "b", b)
		return builtinMath.MaxUint64
	}

	return res
}

// AddUint64WithErr performs addition on uint64 and returns an error if the addition overflows
func AddUint64WithErr(a, b uint64) (uint64, error) {
	s := a + b
	if s >= a && s >= b {
		return s, nil
	}

	return s, ErrAdditionOverflow
}

// AddInt64 performs addition on int64 and logs an error if the addition overflows
func AddInt64(a, b int64) int64 {
	res := a + b
	if (res > a) == (b > 0) {
		return res
	}

	log.Error("AddInt64 overflow", "a", a, "b", b)
	return builtinMath.MaxInt64
}

// SubUint64 performs subtraction on uint64, in case of underflow returns the difference
func SubUint64(a, b uint64) (uint64, uint64) {
	if a >= b {
		return a - b, 0
	}

	return 0, b - a
}

// MulUint64 performs multiplication on uint64 and logs an error if the multiplication overflows
func MulUint64(a, b uint64) uint64 {
	res, err := MulUint64WithErr(a, b)
	if err != nil {
		log.Error("MulUint64 overflow", "a", a, "b", b)
		return builtinMath.MaxUint64
	}

	return res
}

// MulUint64WithErr performs multiplication on uint64 and returns an error if the multiplication overflows
func MulUint64WithErr(a, b uint64) (uint64, error) {
	res := a * b
	if a == 0 || b == 0 || a == res/b {
		return res, nil
	}

	return 0, ErrMultiplicationOverflow
}

// AddInt32 performs addition on int32 and logs an error if the addition overflows
func AddInt32(a, b int32) int32 {
	res, err := AddInt32WithError(a, b)
	if err != nil {
		log.Error("AddInt32 overflow", "a", a, "b", b)
		return builtinMath.MaxInt32
	}

	return res
}

// AddInt32WithError performs addition on int32 and returns an error if the addition overflows
func AddInt32WithError(a, b int32) (int32, error) {
	res := a + b
	if (res > a) == (b > 0) {
		return res, nil
	}

	return 0, ErrAdditionOverflow
}

// SubInt performs subtraction on int and logs an error if the subtraction overflows
func SubInt(a, b int) int {
	res := a - b
	if (res < a) == (b > 0) {
		return res
	}

	log.Error("SubInt underflow", "a", a, "b", b)
	return builtinMath.MinInt64
}

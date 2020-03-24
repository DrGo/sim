package rng

import "unsafe"

func sfrand(seed *int) float32 {
	var f float32
	*seed *= 16807
	*(*uint32)(unsafe.Pointer(&f)) = uint32((uint(*seed) >> 9) | 0x40000000)
	return f - 3.0
}

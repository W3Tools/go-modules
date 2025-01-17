package utils

func UnsignedIntegerToUint64(v any) uint64 {
	switch v := v.(type) {
	case uint64:
		return v
	case uint32:
		return uint64(v)
	case uint16:
		return uint64(v)
	case uint8:
		return uint64(v)
	case uint:
		return uint64(v)
	default:
		return 0
	}
}

package tool

import (
	"strconv"
)

func ItoFloat64(i interface{}) (float64, bool) {
	switch v := i.(type) {
	case float64:
		return v, true
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, true
		}
		return 0, false
	default:
		return 0, false
	}
}

package models

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cyrushiker/pigo/pkg/tool"
)

type Doct struct {
	Id         string    `json:"id" dt:"keyword"`
	Name       string    `json:"name,omitempty" dt:"keyword"`
	Key        string    `json:"key" dt:"keyword"`
	Desc       string    `json:"desc,omitempty" dt:"text"`
	CreateDate time.Time `json:"createDate,omitempty" dt:"date"`
	Position   int
}

func (d *Doct) esTypeName() string {
	return "docut"
}

func (d *Doct) Add() error {
	return nil
}

type DKey struct {
}

type KGroup struct {
}

type KeyDefine struct {
	Id     string
	Code   string
	Groups []struct {
		Id       string
		Key      string
		Name     string
		Rank     int
		MetaKeys []struct {
			Id       string
			Key      string
			Name     string
			Keyword  string
			Etype    string
			Optional string
		}
	}
}

type DocVerify struct {
	kd           *KeyDefine
	doc          map[string]interface{}
	raiseOnError bool
	verifies     []interface{}
}

func (dv *DocVerify) Verify() {
	for _, g := range dv.kd.Groups {
		if d, ok := dv.doc[g.Key]; ok {
			logger.Println(d)
			if d, ok := d.([]interface{}); ok {
				logger.Println(d)
			}
		}
	}
}

type verifyError struct {
	K string
	V interface{}
	T string
}

func (e verifyError) Error() string {
	return fmt.Sprintf(`%s: "%v" is not wanted %s`, e.K, e.V, e.T)
}

var valueUnitKM = [...]string{"value", "value_str", "unit", "isnormal", "origin", "gt", "gte", "lt", "lte", "eq"}

// verify group key-value with json unmarshal rules
// see `https://golang.org/src/encoding/json/decode.go#L50`
func verifySwitch(val interface{}, k, kt, rr, unit string) (interface{}, error) {
	if val == nil {
		return nil, fmt.Errorf("value of %s is null", k)
	}
	switch kt {
	case "Normal", "String":
		if v, ok := val.(string); ok {
			return strings.TrimSpace(v), nil
		}
		if v, ok := val.(map[string]interface{}); ok {
			return fmt.Sprintf("%v", v), nil
		}
		return nil, verifyError{k, val, "string"}
	case "Number":
		if f, ok := tool.ItoFloat64(val); ok {
			return f, nil
		}
		return nil, verifyError{k, val, "number"}
	case "RegExp":
		if s, ok := val.(string); !ok {
			return nil, fmt.Errorf("regexp value must be string")
		} else {
			rr = strings.Trim(rr, "/")
			
			reg, err := regexp.Compile(rr)
			if err != nil {
				return nil, err
			}
			if !reg.MatchString(s) {
				return nil, fmt.Errorf(`"%s" is not match the ruls "/%s/"`, s, rr)
			}
			return s, nil
		}
	case "Date":
		if val == float64(-1) || val == "-1" {
			return nil, fmt.Errorf("date cannot be -1")
		}
		if f, ok := val.(float64); ok {
			return int64(f), nil
		}
		if s, ok := val.(string); ok {
			s = strings.TrimSpace(s)
			if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
				return t.UnixNano() / 1e6, nil
			}
			if t, err := time.Parse("2006-01-02 15:04", s); err == nil {
				return t.UnixNano() / 1e6, nil
			}
			if t, err := time.Parse("2006-01-02", s); err == nil {
				return t.UnixNano() / 1e6, nil
			}
			if f, err := strconv.ParseFloat(s, 64); err == nil {
				return int64(f), nil
			}
		}
		return nil, verifyError{k, val, "date"}
	case "ArrayString":
		if s, ok := val.([]interface{}); ok {
			ss := []string{}
			for _, i := range s {
				ss = append(ss, fmt.Sprintf("%s", i))
			}
			return ss, nil
		}
		return nil, verifyError{k, val, "array-string"}
	case "Multiple":
		if s, ok := val.([]interface{}); ok {
			ss := []string{}
			for _, i := range s {
				ss = append(ss, fmt.Sprintf("%v", i))
			}
			return strings.Join(ss, ","), nil
		}
		return nil, verifyError{k, val, "multiple-value"}
	case "ValueUnit":
		vu := make(map[string]interface{})
		switch val := val.(type) {
		case float64:
			vu["value"] = val
			vu["unit"] = unit
			return vu, nil
		case map[string]interface{}:
			for _, k := range valueUnitKM {
				if v, ok := val[k]; ok {
					if k == "value" {
						if f, ok := tool.ItoFloat64(v); ok {
							vu[k] = f
						}
					} else {
						if v != nil && v != "" && v != "null" {
							vu[k] = v
						}
					}
				}
			}
			_, ok1 := vu["value"]
			_, ok2 := vu["value_str"]
			if !ok1 && !ok2 {
				return nil, fmt.Errorf("no one of value and value_str exists")
			}
			return vu, nil
		case string:
			val = strings.TrimSpace(val)
			if len(val) < 1 {
				return nil, fmt.Errorf("value-unit can not be empty string")
			}
			div := strings.Index(val, " ")
			if div > -1 {
				vu["value"] = val[:div]
				vu["unit"] = val[div+1:]
			} else {
				if f, ok := tool.ItoFloat64(val); ok {
					vu["value"] = f
				} else {
					vu["value_str"] = val
				}
			}
			return vu, nil
		default:
			return nil, verifyError{k, val, "value-unit"}
		}
	default:
		return fmt.Sprintf("%v", val), nil
	}
}

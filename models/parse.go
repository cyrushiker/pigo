package models

import (
	"fmt"
	"strings"
)

type bnode struct {
	elem  string
	left  *bnode
	right *bnode
}

var symbol = []string{
	"(",
	")",
	"and", // -> &
	"or",  // -> |
	"not", // -> ^
	">",
	"<",
	">=", // -> ]
	"<=", // -> [
	"=",
}

// eg. (姓名 = 吴秉礼 or (男 肌酐 > 0)) and (白细胞<20 红细胞 > 10)
func parse(query string) *bnode {
	if len(query) < 0 {
		return nil
	}
	query = strings.ReplaceAll(query, "and", "∩")
	query = strings.ReplaceAll(query, "or", "∪")
	query = strings.ReplaceAll(query, "not", "≠")
	query = strings.ReplaceAll(query, ">=", "≥")
	query = strings.ReplaceAll(query, "<=", "≤")
	logger.Println(query)

	// var ss []rune
	// var mv []rune
	// for _, r := range []rune(query) {
	// 	switch r {
	// 	default:

	// 	}
	// }
	return nil
}

func convert(input string) {
	var s []bnode
	for _, e := range input {
		switch e {
		case '+', '-', '*':
			n := len(s) - 1
			r := s[n]
			l := s[n-1]
			s = s[:n-1]
			s = append(s, bnode{elem: fmt.Sprintf("%c", e), left: &l, right: &r})
		default:
			s = append(s, bnode{elem: fmt.Sprintf("%c", e)})
		}
	}
	logger.Printf("%#v", s[0])
}

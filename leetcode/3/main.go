package main

import (
	"fmt"
)

func maxInt(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func lengthOfLongestSubstring(s string) int {
	l, ans := len(s), 0
	m := make(map[byte]int)
	for i, j := 0, 0; j < l; j++ {
		if _i, ok := m[s[j]]; ok == true {
			i = maxInt(_i, i)
		}
		ans = maxInt(ans, j-i+1)
		m[s[j]] = j + 1
		if ans >= l-i {
			break
		}
	}
	return ans
}

func main() {
	r := lengthOfLongestSubstring("abcabcbb")
	fmt.Printf("%#v", r)
}

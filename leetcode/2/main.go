package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	ret := new(ListNode)
	p, q, c := l1, l2, ret
	carry := 0
	for p != nil || q != nil {
		var x, y int
		if p != nil {
			x, p = p.Val, p.Next
		}
		if q != nil {
			y, q = q.Val, q.Next
		}
		sum := x + y + carry
		carry = sum / 10
		c.Next = &ListNode{Val: sum % 10}
		c = c.Next
	}
	if carry > 0 {
		c.Next = &ListNode{Val: carry}
	}
	return ret.Next
}

func main() {
	r := addTwoNumbers(&ListNode{Val: 4}, &ListNode{Val: 9})
	fmt.Printf("%#v", r)
}

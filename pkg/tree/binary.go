package tree

import (
	"fmt"
	"regexp"
	"strings"
)

// 一般树的实现，
type treeNode struct {
	element     interface{}
	firstChild  *treeNode
	nextSibling *treeNode
}

func preorderTraversal(t *treeNode) {
	var listAll func(int, *treeNode)
	listAll = func(deep int, n *treeNode) {
		fmt.Printf("%*s %s\n", deep*4, " ", n.element)
		for c := n.firstChild; c != nil; c = c.nextSibling {
			deep++
			listAll(deep, c)
		}
	}
	listAll(0, t)
}

type binaryNode struct {
	element interface{}
	left    *binaryNode
	right   *binaryNode
}

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
	fmt.Println(query)

	termRe := regexp.MustCompile(`([^()∩∪≠≥≤=><\s])+`)
	terms := termRe.FindAllString(query, -1)
	fmt.Println(terms)

	cf := make([]string, 0, 10)
	rep := termRe.ReplaceAllStringFunc(query, func(m string) string {
		cf = append(cf, m)
		return fmt.Sprintf("%c", 'A'+len(cf)-1)
	})
	fmt.Println(rep)

	termRe2 := regexp.MustCompile(`([^()∩∪≠≥≤=><\s])+\s*[≥≤=><]\s*([^()∩∪≠≥≤=><\s])+`)
	terms2 := termRe2.FindAllString(query, -1)
	fmt.Println(terms2)

	cf2 := make([]string, 0, 10)
	rep2 := termRe2.ReplaceAllStringFunc(query, func(m string) string {
		cf2 = append(cf2, m)
		return fmt.Sprintf("%c", 'A'+len(cf2)-1)
	})
	fmt.Println(cf2)
	fmt.Println(rep2) // (A ∪ (B C)) ∩ (D E)

	// rep3 := regexp.MustCompile(`([A-Z])(\s*)([A-Z])`).ReplaceAllString(rep2, "$1∩$3")
	rep3 := strings.ReplaceAll(rep2, " ", "")
	fmt.Println(rep3) // output: (A∪(BC))∩(DE)

	return nil
}

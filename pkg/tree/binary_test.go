package tree

import (
	"testing"
)

func TestNormal(t *testing.T) {
	root := &treeNode{
		element: "/usr",
		firstChild: &treeNode{
			element: "/share",
			firstChild: &treeNode{
				element: "/python",
				nextSibling: &treeNode{
					element: "/lib",
					firstChild: &treeNode{
						element: "/node",
						nextSibling: &treeNode{
							element: "/modules",
						},
					},
					nextSibling: &treeNode{
						element: "/bin",
						firstChild: &treeNode{
							element: "/log",
						},
						nextSibling: &treeNode{
							element: "/build",
						},
					},
				},
			},
			nextSibling: &treeNode{
				element: "/bill",
				firstChild: &treeNode{
					element: "/me",
				},
				nextSibling: &treeNode{
					element: "/fuck",
				},
			},
		},
	}
	preorderTraversal(root)
}

func TestParse(t *testing.T) {
	parse("(姓名 = 吴秉礼 or (性别 = 男 肌酐 > 0)) and (白细胞<=20 红细胞 >= 10)")
	parse("姓名 = 吴秉礼 性别 = 男 肌酐 > 0 白细胞<=20 红细胞 >= 10")
}

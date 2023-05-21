package intake

//
//type Node struct {
//	value    string
//	children map[rune]*Node
//	handlers map[string]http.HandlerFunc
//}
//
//func NewNode() *Node {
//	return &Node{
//		children: make(map[rune]*Node),
//		handlers: make(map[string]http.HandlerFunc),
//	}
//}
//
//func (n *Node) Insert(route string, method string, handler http.HandlerFunc) {
//	current := n
//	for _, r := range route {
//		child, ok := current.children[r]
//		if !ok {
//			child = NewNode()
//			current.children[r] = child
//		}
//		current = child
//	}
//	current.value = route
//	current.handlers[method] = handler
//}
//
//func (n *Node) Find(route string, method string) (http.HandlerFunc, bool) {
//	current := n
//	for _, r := range route {
//		child, ok := current.children[r]
//		if !ok {
//			return nil, false
//		}
//		current = child
//	}
//	handler, ok := current.handlers[method]
//	return handler, ok
//}
//
//// PrintTree Chat gippity generated tree viewer
//func (n *Node) PrintTree(prefix string, last bool) {
//	var nodePrefix string
//	var childPrefix string
//
//	if last {
//		nodePrefix = prefix + "└── "
//		childPrefix = prefix + "    "
//	} else {
//		nodePrefix = prefix + "├── "
//		childPrefix = prefix + "│   "
//	}
//
//	if n.value != "" {
//		fmt.Printf("%s%s\n", nodePrefix, n.value)
//	} else {
//		fmt.Printf("%s*\n", nodePrefix)
//	}
//
//	childValues := make([]rune, 0, len(n.children))
//	for r := range n.children {
//		childValues = append(childValues, r)
//	}
//	sort.Slice(childValues, func(i, j int) bool {
//		return childValues[i] < childValues[j]
//	})
//
//	for i, r := range childValues {
//		n.children[r].PrintTree(childPrefix, i == len(childValues)-1)
//	}
//}

package tfb

import "strings"

/*
	开发服务时，注册路由规则，映射handler；访问时，匹配路由规则，查找到对应的handler
	因此，Trie 树需要支持节点的插入与查询
*/

// trie 树(前缀树) 用于动态路由匹配
type node struct {
	pattern  string  // 待匹配路由， 例如 /p/:lang
	part     string  // 路由中的一部分， 例如 :lang
	children []*node // 子节点， 例如 [doc, xx, file]
	isWild   bool    // 是否精确匹配， part 含有 : 或者 * 时为true
}

// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// trie树插入节点
// 递归查找每一层的节点，如果没有匹配到当前part的节点，则新建一个
func (n *node) insert(pattern string, parts []string, height int) {
	// 匹配完所有需要匹配的节点，匹配结束
	if len(parts) == height {
		// 设置匹配路由 url
		n.pattern = pattern
		return
	}

	part := parts[height]
	// 查找第一个与之相匹配的节点
	child := n.matchChild(part)
	if child == nil {
		// 没匹配到节点， 新建一个节点
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		// 新建节点加入当前的trie树
		n.children = append(n.children, child)
	}
	// 递归遍历
	child.insert(pattern, parts, height+1)
}

// trie树中查找节点
// 递归查询每一层的节点，退出规则是，匹配到了*，匹配失败，或者匹配到了第len(parts)层节点
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 匹配路由url 为空说明匹配失败不存在
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	// 返回所有与之相匹配的结点
	children := n.matchChildren(part)
	// 遍历匹配的节点
	for _, child := range children {
		// 递归遍历匹配的节点，找到则返回
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}

	}
	return nil
}

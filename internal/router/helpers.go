package router

import (
	"fmt"
	"sort"
)

func (router *Router) PrintTree() {
	fmt.Println("Router Tree:")
	printNode(router.rootNode, "", true)
}

func printNode(node *RouterTreeNode, prefix string, isLast bool) {
	if node == nil {
		return
	}

	// Determine the current branch character
	connector := "├── "
	if isLast {
		connector = "└── "
	}

	// Print current node
	nodeName := node.name
	if nodeName == "" {
		nodeName = "(root)"
	}

	handlerInfo := ""
	if node.handler != nil {
		handlerInfo = fmt.Sprintf(" [Handler: %s %s]",
			node.handler.Method, node.handler.Pattern)
	}

	fmt.Printf("%s%s%s%s\n", prefix, connector, nodeName, handlerInfo)

	// Prepare prefix for children
	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	// Count total children
	totalChildren := len(node.staticChildren)
	if node.param != nil {
		totalChildren++
	}
	if node.wildcard != nil {
		totalChildren++
	}

	currentChild := 0

	// Print static children
	// Sort keys for consistent output
	staticKeys := make([]string, 0, len(node.staticChildren))
	for key := range node.staticChildren {
		staticKeys = append(staticKeys, key)
	}
	sort.Strings(staticKeys)

	for _, key := range staticKeys {
		child := node.staticChildren[key]
		currentChild++
		printNode(child, childPrefix, currentChild == totalChildren)
	}

	// Print param node
	if node.param != nil {
		currentChild++
		fmt.Printf("%s%s:%s%s\n",
			childPrefix,
			getConnector(currentChild == totalChildren),
			node.param.name,
			getHandlerInfo(node.param.handler))

		printParamChildren(node.param, childPrefix, currentChild == totalChildren)
	}

	// Print wildcard node
	if node.wildcard != nil {
		currentChild++
		fmt.Printf("%s%s*%s\n",
			childPrefix,
			getConnector(currentChild == totalChildren),
			getHandlerInfo(node.wildcard.handler))
	}
}

func printParamChildren(node *RouterTreeNode, prefix string, isLast bool) {
	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	totalChildren := len(node.staticChildren)
	if node.param != nil {
		totalChildren++
	}
	if node.wildcard != nil {
		totalChildren++
	}

	currentChild := 0

	// Print static children
	staticKeys := make([]string, 0, len(node.staticChildren))
	for key := range node.staticChildren {
		staticKeys = append(staticKeys, key)
	}
	sort.Strings(staticKeys)

	for _, key := range staticKeys {
		child := node.staticChildren[key]
		currentChild++
		printNode(child, childPrefix, currentChild == totalChildren)
	}

	// Print param node
	if node.param != nil {
		currentChild++
		fmt.Printf("%s%s:%s%s\n",
			childPrefix,
			getConnector(currentChild == totalChildren),
			node.param.name,
			getHandlerInfo(node.param.handler))
		printParamChildren(node.param, childPrefix, currentChild == totalChildren)
	}

	// Print wildcard node
	if node.wildcard != nil {
		currentChild++
		fmt.Printf("%s%s*%s\n",
			childPrefix,
			getConnector(currentChild == totalChildren),
			getHandlerInfo(node.wildcard.handler))
	}
}

func getConnector(isLast bool) string {
	if isLast {
		return "└── "
	}
	return "├── "
}

func getHandlerInfo(handler *Handler) string {
	if handler != nil {
		return fmt.Sprintf(" [Handler: %s %s]", handler.Method, handler.Pattern)
	}
	return ""
}

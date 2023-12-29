package utils

import "strings"

type HandlerPath struct {
	PathParts    []string
	HandlerParts []string
}

func SplitHandlerPath(handlerPath string) HandlerPath {
	// Create path parts by splitting on slashes
	pathParts := strings.Split(handlerPath, "/")

	// Create handler parts by splitting on dots from the last path part
	handlerParts := strings.Split(pathParts[len(pathParts)-1], ".")

	// Remove the last dot-separated part from the path parts
	pathParts = pathParts[:len(pathParts)-1]
	pathParts = append(pathParts, handlerParts[0])

	// Remove the first handler part from the handler parts
	handlerParts = handlerParts[1:]

	return HandlerPath{
		PathParts:    pathParts,
		HandlerParts: handlerParts,
	}
}

func (h HandlerPath) PackagePath() string {
	return strings.Join(h.PathParts, "/")
}

func (h HandlerPath) PackageName() string {
	return h.PathParts[len(h.PathParts)-1]
}

func (h HandlerPath) Handler() string {
	return strings.Join(h.HandlerParts, ".")
}

func (h HandlerPath) FuncName() string {
	return h.HandlerParts[0]
}

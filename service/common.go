package service

import "strings"

// ParentDir
//
//	@param path
//	@return string
func ParentDir(path string) string {
	parts := strings.Split(path, "/")
	return strings.Join(parts[0:len(parts)-1], "/")
}

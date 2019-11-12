package web

import (
	"fmt"
)

// IndexPath returns index path under resource root path.
func IndexPath() string {
	return ""
}

// EditPath returns edit path under resource root path.
func EditPath() string {
	return "/{id}/edit"
}

// NewPath returns new path under resource root path.
func NewPath() string {
	return "/new"
}

// ShowPath returns show path under resource root path.
func ShowPath() string {
	return "/{id}"
}

// CreatePath returns create path under resource root path.
func CreatePath() string {
	return ""
}

// UpdatePath returns update path under resource root path.
func UpdatePath() string {
	return "/{id}"
}

// InitDeletePath returns init delete path under resource root path.
func InitDeletePath() string {
	return "/{id}/init-delete"
}

// DeletePath returns delete path under resource root path.
func DeletePath() string {
	return "/{id}"
}

// SignupPath returns signup path.
func SignupPath() string {
	return "/signup"
}

// LoginPath returns login path.
func LoginPath() string {
	return "/login"
}

// ResPath

// ResPath
func ResPath(rootPath string) string {
	return "/" + rootPath + IndexPath()
}

// ResPathEdit
func ResPathEdit(rootPath string, r Identifiable) string {
	return fmt.Sprintf("/%s/%s/edit", rootPath, r.GetSlug())
}

// ResPathNew
func ResPathNew(rootPath string) string {
	return fmt.Sprintf("/%s/new", rootPath)
}

// ResPathInitDelete
func ResPathInitDelete(rootPath string, r Identifiable) string {
	return fmt.Sprintf("/%s/%s/init-delete", rootPath, r.GetSlug())
}

// ResPathID
func ResPathID(rootPath string, r Identifiable) string {
	return fmt.Sprintf("/%s/%s", rootPath, r.GetSlug())
}

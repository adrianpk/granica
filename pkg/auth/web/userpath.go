package web

import (
	"gitlab.com/mikrowezel/backend/web"
)

// UserRoot - User resource root path.
var UserRoot = "users"

// UserPath
func UserPath() string {
	return web.ResPath(UserRoot)
}

// UserPathEdit
func UserPathEdit(res web.Identifiable) string {
	return web.ResPathEdit(UserRoot, res)
}

// UserPathNew
func UserPathNew() string {
	return web.ResPathNew(UserRoot)
}

// UserPathInitDelete
func UserPathInitDelete(res web.Identifiable) string {
	return web.ResPathInitDelete(UserRoot, res)
}

// UserPathSlug
func UserPathSlug(res web.Identifiable) string {
	return web.ResPathSlug(UserRoot, res)
}

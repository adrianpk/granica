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
	// NOTE: As an exception for User model username
	// is used as main identifier and not Slug.
	// TODO: Analize if in a multi-tenant setup this could be
	// a problem.
	return web.ResPathEdit(UserRoot, res)
	//return fmt.Sprintf("/%s/%s/edit", UserRoot, res.U)
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

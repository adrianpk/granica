package web

import "gitlab.com/mikrowezel/backend/web"

// UserRoot - User resource root path.
var UserRoot = "users"

// UserPath
func UserPath() string {
	return ResPath(UserRoot)
}

// UserPathEdit
func UserPathEdit(res web.Identifiable) string {
	return ResPathEdit(UserRoot, res)
}

// UserPathNew
func UserPathNew() string {
	return ResPathNew(UserRoot)
}

// UserPathInitDelete
func UserPathInitDelete(res web.Identifiable) string {
	return ResPathInitDelete(UserRoot, res)
}

// UserPathSlug
func UserPathID(res web.Identifiable) string {
	return ResPathSlug(UserRoot, res)
}

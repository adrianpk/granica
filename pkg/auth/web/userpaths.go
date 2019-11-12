package web

// UserRoot - User resource root path.
var UserRoot = "users"

// UserPath
func UserPath() string {
	return ResPath(UserRoot)
}

// UserPathEdit
func UserPathEdit(model Identifiable) string {
	return ResPathEdit(UserRoot, model)
}

// UserPathNew
func UserPathNew() string {
	return ResPathNew(UserRoot)
}

// UserPathInitDelete
func UserPathInitDelete(model Identifiable) string {
	return ResPathInitDelete(UserRoot, model)
}

// UserPathID
func UserPathID(model Identifiable) string {
	return ResPathID(UserRoot, model)
}

package web

import "html/template"

var pathFxs = template.FuncMap{
	// User
	"userPath":           UserPath,
	"userPathEdit":       UserPathEdit,
	"userPathID":         UserPathID,
	"userPathInitDelete": UserPathInitDelete,
	"userPathNew":        UserPathNew,
}

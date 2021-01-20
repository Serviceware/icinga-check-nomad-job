package internal

import (
	"fmt"
)

const (
	OK       = 0
	WARNING  = 1
	CRITICAL = 2
	UNKNOWN  = 3
)

func createJobLink(address string, job string) string {
	link := fmt.Sprintf("%s/ui/jobs/%s", address, job)

	return "<a href=\"" + link + "\" target=\"_blank\">" + link + "</a>"
}

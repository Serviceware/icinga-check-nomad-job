package internal

import (
	"fmt"
)

type Status uint

const (
	OK       Status = 0
	WARNING  Status = 1
	CRITICAL Status = 2
	UNKNOWN  Status = 3
)

func (s Status) Max(other Status) Status {
	if s > other {
		return s
	} else {
		return other
	}
}

func createJobLink(address string, job string) string {
	link := fmt.Sprintf("%s/ui/jobs/%s", address, job)

	return "<a href=\"" + link + "\" target=\"_blank\">" + link + "</a>"
}

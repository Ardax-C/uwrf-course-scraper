package utils

import "strings"

func IsValidLink(link string) bool {
	return strings.Contains(link, "courseLightbox.cfm?subject=CIDS")
}

package utils

import "strings"

func isValidLink(link string) bool {
	return strings.Contains(link, "courseLightbox.cfm?subject=CIDS")
}

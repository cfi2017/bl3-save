package shared

import "strings"

func GuessIsProfileSav(name string) bool {
	parts := strings.Split(name, "/")
	name = parts[len(parts) - 1]
	if name == "profile.sav" {
		return true
	}
	return false
}

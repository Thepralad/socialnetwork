package security

import "strings"

func IsSalesianEmail(email string) bool {
	return strings.HasSuffix(email, "@salesiancollege.net")
}

package messaging

import "strconv"

const (
	maxReceivers = 255
	maxBodySize  = 1024 //KB (1 MB)
)

// ValidateRequest - validates an incoming Relay message sub count and body size in line with requirements
func ValidateRequest(subs []Subscription, body string) (okSubs bool, okBody bool, retMsg string) {
	if len(subs) <= maxReceivers {
		okSubs = true
		retMsg = ""
	} else {
		okSubs = false
		retMsg = "Too many clientIDs provided. MAX: " +
			strconv.FormatUint(maxReceivers, 10)
	}
	okBody = true //TODO
	return
}

package messaging

const (
	maxReceivers = 255
	maxBodySize  = 1024 //KB (1 MB)
)

// ValidateRequest - validates an incoming Relay message sub count and body size in line with requirements
func ValidateRequest(subs []Subscription, body string) (okSubs bool, okBody bool) {
	okSubs = len(subs) <= maxReceivers
	okBody = true //TODO
	return
}

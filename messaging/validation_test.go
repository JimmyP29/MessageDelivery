package messaging

import "testing"

type ValidateRequestResult struct {
	subs   []Subscription
	body   string
	okSubs bool
	okBody bool
}

var validateRequestResults = []ValidateRequestResult{
	{make([]Subscription, 100, 100), "foobar", true, true},
	{make([]Subscription, 256, 256), "foobar", false, true},
}

func TestValidateRequest(t *testing.T) {
	for _, test := range validateRequestResults {
		okSubs, okBody := ValidateRequest(test.subs, test.body)
		if okSubs != test.okSubs {
			t.Fatal("Expected result for subs not given")
		}
		if okBody != test.okBody {
			t.Fatal("Expected result for body not given")
		}
	}
}

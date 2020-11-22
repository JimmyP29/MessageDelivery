package messaging

import "testing"

type GetRequestedSubscriptionsResult struct {
	ids    []uint64
	subs    []Subscription
}

var getRequestedSubscriptionsResults = []GetRequestedSubscriptionsResult{
	{[123], []Subscription{topic = "testing", client{userID = 123}}},
	//{"", nil, false},
}

func TestGetRequestedSubscriptions(t *testing.T) {
	for _, test := range getRequestedSubscriptionsResults {
		subs := getRequestedSubscriptions(test.ids)
		//result := bytes.Compare(b, test.b)

		// if result != 0 {
		// 	t.Fatalf("Expected result: %v \n Actual result: %v\n", test.b, b)
		// }
		if subs[0].client.userID != test.subs[0].client.userID {
			t.Fatalf("Expected result: %v \n Actual result: %v\n", test.subs[0].client.userID, subs[0].client.userID)
		}
	}
}
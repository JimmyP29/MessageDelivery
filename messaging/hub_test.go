package messaging

import "testing"

var t = "test"
var c = CreateConnection()

var client1 = Client{
	userID:     123,
	connection: &c,
}

var client2 = Client{
	userID:     456,
	connection: &c,
}

var sub1 = Subscription{
	topic:  t,
	client: &client1,
}

var sub2 = Subscription{
	topic:  t,
	client: &client2,
}

var cs = make([]Client, 0)
var ss = make([]Subscription, 0)
var clients []Client = append(cs, client1, client2)
var subs []Subscription = append(ss, sub1, sub2)

var hub = Hub{
	clients:       clients,
	subscriptions: subs,
}

type NewHubResult struct {
	clients       []Client
	subscriptions []Subscription
	expected      *Hub
}

var newHubResults = []NewHubResult{
	{hub.clients, hub.subscriptions, &hub},
}

// func TestNewHub(t *testing.T) {
// 	for _, test := range newHubResults {
// 		hub := NewHub(test.clients, test.subscriptions)
// 		if hub != test.expected {
// 			t.Fatalf("Expected result: %v \n Actual result: %v\n", test.expected, hub)
// 		}
// 	}
// }

type GetRequestedSubscriptionsResult struct {
	ids      []uint64
	expected []Subscription
}

var getRequestedSubscriptionsResults = []GetRequestedSubscriptionsResult{
	{[]uint64{123}, hub.subscriptions},
}

func TestGetRequestedSubscriptions(t *testing.T) {
	for _, test := range getRequestedSubscriptionsResults {
		subs := hub.getRequestedSubscriptions(test.ids)

		if subs[0].client.userID != test.expected[0].client.userID {
			t.Fatalf("Expected result: %v \n Actual result: %v\n", test.expected[0].client.userID, subs[0].client.userID)
		}
	}
}

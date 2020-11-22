package messaging

import (
	"testing"
)

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

var client3 = Client{
	userID:     789,
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

type AddClientResult struct {
	client   Client
	expected *Hub
}

var addClientResults = []AddClientResult{
	{client3, &hub},
}

func TestAddClient(t *testing.T) {
	for _, test := range addClientResults {
		h := hub.AddClient(test.client)

		if h != test.expected {
			t.Fatalf("Expected result: %v \n Actual result: %v\n", test.expected, h)
		}
	}
}

type GetSubscriptionsResult struct {
	client   *Client
	expected []Subscription
}

var getSubscriptionsResults = []GetSubscriptionsResult{
	{&hub.clients[0], hub.subscriptions},
}

func TestGetSubscriptions(t *testing.T) {
	for _, test := range getSubscriptionsResults {
		subs := hub.getSubscriptions(test.client)

		if subs[0].client.userID != test.expected[0].client.userID {
			t.Fatalf("Expected result: %v \n Actual result: %v\n", test.expected[0].client.userID, subs[0].client.userID)
		}
	}
}

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

type PublishToSenderResult struct {
	message  []byte
	client   *Client
	expected bool
}

var publishToSenderResults = []PublishToSenderResult{
	{[]byte{34, 102, 111, 111, 98, 97, 114, 34}, &hub.clients[0], true},
}

// func TestPublishToSender(t *testing.T) {
// 	for _, test := range publishToSenderResults {
// 		isOK := hub.publishToSender(test.message, test.client)

// 		if isOK != test.expected {
// 			t.Fatalf("Expected result: %v \n Actual result: %v\n", test.expected, isOK)
// 		}
// 	}
// }

type PublishToReceiversResult struct {
	message  []byte
	subs     []Subscription
	expected bool
}

var publishToReceiverResults = []PublishToReceiversResult{
	{[]byte{34, 102, 111, 111, 98, 97, 114, 34}, subs, true},
}

// func TestPublishToReceivers(t *testing.T) {
// 	for _, test := range publishToReceiverResults {
// 		isOK := hub.publishToReceivers(test.message, test.subs)

// 		if isOK != test.expected {
// 			t.Fatalf("Expected result: %v \n Actual result: %v\n", test.expected, isOK)
// 		}
// 	}
// }

type HandleIdentityResult struct {
	client   *Client
	expected bool
}

var handleIdentityResults = []HandleIdentityResult{
	{&hub.clients[0], true},
}

// func TestHandleIdentity(t *testing.T) {
// 	for _, test := range handleIdentityResults {
// 		isOK := hub.handleIdentity(test.client)

// 		if isOK != test.expected {
// 			t.Fatalf("Expected result: %v \n Actual result: %v\n", test.expected, isOK)
// 		}
// 	}
// }

type HandleListResult struct {
	client   *Client
	expected bool
}

var handleListResults = []HandleListResult{
	{&hub.clients[0], true},
}

// func TestHandleList(t *testing.T) {
// 	for _, test := range handleListResults {
// 		isOK := hub.handleList(test.client)

// 		if isOK != test.expected {
// 			t.Fatalf("Expected result: %v \n Actual result: %v\n", test.expected, isOK)
// 		}
// 	}
// }

var message = Message{
	MsgType:   2,
	Body:      []byte{34, 102, 111, 111, 98, 97, 114, 34},
	ClientIDS: []uint64{456},
}

type HandleRelayResult struct {
	client   *Client
	message  *Message
	expected bool
}

var handleRelayResults = []HandleRelayResult{
	{&hub.clients[0], &message, true},
}

// func TestHandleRelay(t *testing.T) {
// 	for _, test := range handleRelayResults {
// 		isOK := hub.handleRelay(*test.client, *test.message)

// 		if isOK != test.expected {
// 			t.Fatalf("Expected result: %v \n Actual result: %v\n", test.expected, isOK)
// 		}
// 	}
// }

type HandleDefaultResult struct {
	client   *Client
	expected bool
}

var handleDefaultResults = []HandleDefaultResult{
	{&hub.clients[0], true},
}

// func TestHandleDefault(t *testing.T) {
// 	for _, test := range handleDefaultResults {
// 		isOK := hub.handleDefault(test.client)

// 		if isOK != test.expected {
// 			t.Fatalf("Expected result: %v \n Actual result: %v\n", test.expected, isOK)
// 		}
// 	}
// }

//var p = SerialiseMessage(&message)

type HandleReceiveMessageResult struct {
	client   *Client
	payload  []byte
	expected *Hub
}

var handleReceiveMessageResults = []HandleReceiveMessageResult{
	{&hub.clients[0], []byte{34, 102, 111, 111, 98, 97, 114, 34}, &hub},
}

func TestHandleReceiveMessage(t *testing.T) {
	for _, test := range handleReceiveMessageResults {
		h := hub.HandleReceiveMessage(*test.client, test.payload)

		if h != test.expected {
			t.Fatalf("Expected result: %v \n Actual result: %v\n", test.expected, h)
		}
	}
}

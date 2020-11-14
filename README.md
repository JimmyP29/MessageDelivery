# MessageDelivery
A simple message delivery system written in Go.

There will be a 
- `Hub` - Relays incoming message bodies to receivers based on user ID(s) defined in the message.

- `Clients` - Users who are connected to the hub. The client may send three types of messages.
- `Messages` 
    - ### Identity message
        The client can send an identity message which the hub will answer with the user_id of the
        connected user.
    - ### List message
        The client can send a list message which the hub will answer with the list of all connected client
        user_id's (excluding the requesting client).
    - ### Relay message    
        The client can send a relay message whose body is relayed to receivers marked in the message. The message body can be relayed to one or multiple
        receivers.



## Aim
---
To learn as well as demonstrate my capabilities at:
- Implementing the above idea as best I can.
- Following best practises when it comes to Go project/package structures as well as naming conventions and documentation etc as found at https://golang.org/doc/
- Using a Gitflow branching approach (albeit - on my own). This will take the form of `task/` branches off of `master`, I am going to do my best to try and break down the work in as vertical way as I can, I think this will prove to be challenging.
- `Testing` - I will be implementing tests within my packages for sure. I will *not* be taking a TDD approach to this however as I have never worked that way before anyway. I have enough to get used to as it is so I am not adding that to the pile! :sweat_smile:

*__NOTE: I am going to use this README to log my progress with this project. This exercise, to me; is as much about learning as it is demonstrating so I think this will be (alongside the git history) - valuable to the reader.__* 

## Initial Thoughts
---
I need to ensure that communication between the Hub and Clients is done on the network layer. My mind wandered to the following 3 places...

`http`
- At first I was thinking this could be done with http. A client could quite easily evoke a network request and receive a response, but I don't think this is suitable for the Relay message as it isn't a 1 to 1 request/response. The amount of messages that would have to be returned is dictated by the amount of user_id's. Simply not fit for purpose.

`Websockets`
- I  have briefly considered websockets to send messages over a stream, but am not keen on using them for this.

`Topics`
- I had thought of 'pubsub' before http because of the way the message behaviour reminded me of utilising service buses with Topics using .Net Core and Azure. I thought it would maybe not be appropriate for the scope of this project as I want this to be a self contained solution that doesn't require having to setup something so complicated. However I found [this Go package](https://godoc.org/cloud.google.com/go/pubsub), which seems to abstract all the cloud stuff away from what you're doing so I am going to pursue that. If this ends up being a dead end I am thinking of trying to replicate something similar anyway.

## Planning
--- 
__Pen and paper time.__

I'm going to draw out how I think the data structures and interfaces should be and architect what I am going to do - I am also going to create a spike branch for playing with the `pubsub` package above to see if it suits my needs.

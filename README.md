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

To learn as well as demonstrate my capabilities at:
- Implementing the above idea as best I can.
- Following best practises when it comes to Go project/package structures as well as naming conventions and documentation etc as found at https://golang.org/doc/
- Using a Gitflow branching approach (albeit - on my own). This will take the form of `task/` branches off of ~~`master`~~ `main`__`*`__, I am going to do my best to try and break down the work in as vertically way as I can, I think this will prove to be challenging.
- `Testing` - I will be implementing tests within my packages for sure. I will *not* be taking a TDD approach to this however as I have never worked that way before anyway. I have enough to get used to as it is so I am not adding that to the pile! :sweat_smile:

__`*`__ I have just noticed that my master branch is called main, this a mistake on my part as I would normally stick to the normal naming convention. I'm not sure how this happened and to keep things simple I am going to leave it as is.

*__NOTE: I am going to use this README to log my progress with this project. This exercise, to me; is as much about learning as it is demonstrating so I think this will be (alongside the git history) - valuable to the reader.__* 

## Initial Thoughts

I need to ensure that communication between the Hub and Clients is done on the network layer. My mind wandered to the following 3 places...

`http`
- At first I was thinking this could be done with http. A client could quite easily evoke a network request and receive a response, but I don't think this is suitable for the Relay message as it isn't a 1 to 1 request/response. The amount of messages that would have to be returned is dictated by the amount of user_id's. Simply not fit for purpose.

`Websockets`
- I  have briefly considered websockets to send messages over a stream, but am not keen on using them for this.

`Topics`
- I had thought of 'pubsub' before http because of the way the message behaviour reminded me of utilising service buses with Topics using .Net Core and Azure. I thought it would maybe not be appropriate for the scope of this project as I want this to be a self contained solution that doesn't require having to setup something so complicated. However I found [this Go package](https://godoc.org/cloud.google.com/go/pubsub), which seems to abstract all the cloud stuff away from what you're doing so I am going to pursue that. If this ends up being a dead end I am thinking of trying to replicate something similar anyway.

## Planning

__Pen and paper time.__

I'm going to draw out how I think the data structures and interfaces should be and architect what I am going to do - I am also going to create a spike branch for playing with the `pubsub` package above to see if it suits my needs.

__Edit:__ 

I spent a few hours spiking the `pubsub` package, I did need to setup Google Cloud which I hadn't done before. 

I used the example code and with a bit of tweaking was first able to get a message published - which made me very happy.  

![Image of Google Cloud subscribed message](/README_assets/topic_gc_output.jpg "Image of Google Cloud subscribed message")

I then got my code to both publish and subscibe to the Topic, pulling the message body back in the console.

![Image of topic output in the console](/README_assets/topic_output.PNG "Image of topic output in the console")

*However...*
I have 2 issues with this approach:
1. The onus of a client subscribing is on the Topic resource - not the Hub in my application. This makes it harder to control which client recieves what and also doesn't meet the requirement.
2. I don't think the reader will be able to run this all locally as I have my own private API key and subscription to Google Cloud. This isn't going to be put into any kind of config and this application won't be deployed anywhere, so although I learnt some new things - I don't think this is going to be appropriate for this project.

:thinking:

So I am going to revisit the idea of `websockets`. I think I can create a 'pubsub' like mechanism but I can put the control of subscribing within the `Hub`. It should also work within the contained project with no need to connect to cloud services, finally it is still operating on the network layer and as such - meets the requirement. 

I have worked a bit with websockets in the past and they have been awkward to work with and test, also I have never used them where the 'client' isn't a browser - so something new there. The package at [gorilla/websocket](https://github.com/gorilla/websocket) seems simple to implement, and I know that Gorilla is a respected third party. 

I am going to spike again to get a client/server `websocket` running within my project.

__Edit:__ OK so my spike worked but it following a tutorial that was using the browser as the client. I have found [another tutorial on YouTube](https://www.youtube.com/watch?v=Sphme0BqJiY) on how to create a chat application all within Go, creating your own clients. This tutorial is leveraging everything from the standard library and I'm hoping to use this to help guide me along my way when it comes to the details that I am very green at still. 

I know that the requirements said I can use any library I wish and I don't doubt that I could use the gorilla/websocket package to achieve the same thing. I have also read that having the minimum amount of dependencies in your package is a good thing, so I see it as a plus, especially whilst learning. :smile:



I'm breaking the workload as close as I can at first into the following tasks:

- ~~`Task 1` Create initial project structure, set up structs and work flow with sockets to get the `Hub` sending messages to a `Client`.~~ :heavy_check_mark:
- ~~`Task 2` Refine the messages to the 3 different types outlined in the requirements.~~ :heavy_check_mark:
- ~~`Task 3` I would like to squeeze in an additional task at this point for refactoring (also I forgot to add some checks in task 2 __`*`__ , which tbf if I was taking a TDD approach I would have realised sooner and wouldn't have merged to `main`).~~ :heavy_check_mark: Awesome, it works! :grin:.
- ~~`Task 4` Adding Unit tests.~~ :heavy_check_mark:

__`*`__ 
I still need to figure out the message body size problem. I am going to move onto testing for now.

## Finished.

I am proud of what I have achieved in a week. 2 weeks prior to this I had barely written any Go code, so this has been really rewarding for me. I said that vertically slicing the functionality would be difficult at the beginning of this README. I actually thought that once I had got the skelton structure finished in the first task, the details were quite nice to implement on top. I wish I hadn't underestimated the body size validation and the the testing though. I have logged down the detail of these things below however, so that the reader can at least appreciate the thought process behind these decisions.

__How does it work?__

The `Hub` is created upon the project first running.

Using `gorilla/websocket`, a socket is setup in `main.go` as soon as it's run. At the same time a static html file is servered which itself creates a new `WebSocket` object in JavaScript.

*Once the project is running...*

Go to `localhost:8888` in a browser, this automatically creates a new `Client`, assigning them a random __`uint64`__ `userID` and subscribing them to the `Hub` using the `Multiplay` topic. Hit f12 and open the __console__ to interact with the service. 

![Image of 'connected' in browser console](/README_assets/connected.PNG "Image of 'connected' in browser console")

(You're going to want to do this at least twice so there is more than one client)
![Image of (list) with one client in browser console](/README_assets/list_output_lonely.PNG "Nobody wants to be lonely :(")

---

From here fire the following commands to send messages.

- ### Identity
    ```js
    ws.send('{"type": 0 }');
    ```
    ![Image of identity output in browser console](/README_assets/identity_output.PNG "Image of identity output in browser console")
- ### List message
    ```js
    ws.send('{"type": 1 }');
    ```
    ![Image of list output in browser console](/README_assets/list_output.PNG "Image of list output in browser console")
- ### Relay message    
    ```js
    ws.send('{"type": 2, "body": "foobar", "clientIDS": [288, 458] }');
    ```
    Then looking at another client (458 in this example) yields:

    ![Image of relay output in browser console for client 458](/README_assets/relay_output_458.PNG "Image of relay output in browser console for client 458")

        
        
        
        
## Well... Almost.
The only thing that prevents this being 100% code complete - is the validation for the message body size. I wrongly assumed that this would be a case of using some method from the standard lib in order to get this data. 

I played with:
```go
binary.Size(v interface{} int)
```

as well as
```go
unsafe.Sizeof(x ArbitraryType) uintptr
```

Using either one of these with my `[]byte` returned a value that didn't change when increaing the message body and trying again. I believe it returns an expected value based upon the slice type, not what the slice actually will hold in memory. I hope though despite this omission, the reader will be able to see my intent with the `validation.go` file. I would have probably needed to create a second, private method that would be consumed within `ValidateRequest()`, then based on the return, the validation would behave accordingly.

## Testing

I have tried my best to implement the tests. I have done so successfully but there are some that have thrown me. They can be found in `hub_test.go` and have been left commented out to maintain a clean sweep for the tests that pass. 

- `TestNewHub` (line 60): This test seems to be adamant that the actual and expected values are different, even though the printed output reveals them to be identical.
- The other commented out tests seem to be complaining about the same thing, (taking `TestPublishToSender` (line 145) as an example):
  ![Image of comment copy of test output failure](/README_assets/test_output_failure.PNG "Comment copy of test output failure")

I think that they are both related in that they have something to do with reference vs value. I think the second point is because of the hub objects `[]Client` and `[]Subscription`. Somehow there's a derefence issue somewhere.

It has made me learn more about pointers though. I will continue to research these issues and see if I can fix them. I also need to allow myself more time for testing, especially using a framework I am not familiar with.

## Resources
Again, having not used Go much at all in the past, this was a learning experience as much as anything else. In order to do this I needed to find resources to help me achieve what I ended up with. As such I'd like to credit those sources here.
- The [Go docs](https://golang.org/doc/)
- This specific doc for [pubsub](https://godoc.org/cloud.google.com/go/pubsub) (pubsub spike).
- This [Google Cloud article](https://cloud.google.com/pubsub/docs/quickstart-client-libraries#create-topic-sub) (pubsub spike).
- This [YouTube video](https://www.youtube.com/watch?v=Sphme0BqJiY&t=324s) which demonstrated creating a small TCP chat in Go (Websocket spike)
- This [article](https://www.ribice.ba/golang-enums/) about using variables in a way similar to enums in other languages.
- This [YouTube video](https://www.youtube.com/watch?v=yyREnTgRTQ0&t=899s) which demonstrated using a websocket to create a local pubsub implementation. This helped me __a lot__ with the initial setup for the project.
- This [Youtube video](https://www.youtube.com/watch?v=sOeUf1YICSA&list=PLShDm2AZYnK2BEw4ltBF67U3qBamg1ts3) and this [Youtube video](https://www.youtube.com/watch?v=S1O0XI0scOM) with regards to unit testing.
- This [article](https://medium.com/what-i-talk-about-when-i-talk-about-technology/go-code-examples-httptest-newserver-f965fb349884) for creating mocks of websocket connections.

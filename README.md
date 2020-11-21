# dapr-distributed-calendar

This is a sample application built using Dapr as a proof-of-concept. I have experimented with the state store, pubsub and output bindings features available with Dapr.
I have used multiple languages for writing the different parts of this calendar app. This demonstrates the language-agnostic nature of Dapr and the flexibility that it bings to developing
applications. 

## Motivation

I am really enthusiastic about cool open source projects and I'm a fan of Azure. When I learnt about Dapr early this year, I knew I needed to get my hands dirty playing
around with what Dapr had to offer.
I wanted to explore Dapr an experience building a distributed application with it to understand what it brought to the table 
in comparison to conventional applications. 

I had built a SpringBoot app (LINK) on MVCS architecture before; it was a monolith application, all written in Java. 
Building a roughly similar architecture as a distributed applicaiton would intuitively require some additional work pertaining to service discovery, inter-pod communication
and network security. Things could get complicated if I needed additional checks, statestores or other controls which I would have to implement on my own.
This, in addition to the actual application itself. 

I wanted to find out how Dapr simplified this process and what additional work I would have to put in to get a distributed version of the same applciation using Dapr. 

## Architecture

I have tried to model this system on the Model View Controller Service (MVCS) architecture, as already mentioned. 

### Controller (written in *Javascript*): 

  * The controller supports creation of new events and deletion of existing events. 
    It forwards these requests to the **Go** code using service invocation.
  
    *Shown below is the add event flow*. 
  
    ```js
    app.post('/newevent', (req, res) => {
    const data = req.body.data;
    const eventId = data.id;
    console.log("New event registration! Event ID: " + eventId);


    console.log("Data passed as body to Go", JSON.stringify(data))
    fetch(invokeUrl+`/addEvent`, {
        method: "POST",
        body: JSON.stringify(data),
        headers: {
            "Content-Type": "application/json"
        }
    })
    ```
    where the invokeURL is defined as:
    ```js
    const invokeUrl = `http://localhost:${daprPort}/v1.0/invoke/${eventApp}/method`;
    ```
  
  
  * On creation of a new event, it publishes a message to a **pubsub** topic which is then picked up by the **Python** subscriber. 
  
    *Pubishing to the topic*
  
    ```js
    function send_notif(data) {
      var message = {
          "data": {
              "message": data,
          }
      };
      console.log("Message: ", message)
      request( { uri: publishUrl, method: 'POST', json: JSON.stringify(message) } );
    }
    ```
    where the publish URL is:
    ```js
    const publishUrl = `http://localhost:${daprPort}/v1.0/publish/${pubsub_name}/${topic}`;
    ```



# dapr-distributed-calendar

This is a sample application built using Dapr as a proof-of-concept. I have experimented with the state store, pubsub and output bindings features available with Dapr.
I have used multiple languages for writing the different parts of this calendar app. This demonstrates the language-agnostic nature of Dapr and the flexibility that it bings to developing
applications.

## Contents

* [**Motivation**](https://github.com/wjayesh/dapr-distributed-calendar#motivation)

* [**Architecture**](https://github.com/wjayesh/dapr-distributed-calendar#architecture)
  * [**Controller**](https://github.com/wjayesh/dapr-distributed-calendar#controller-written-in-javascript)
  * [**Services**](https://github.com/wjayesh/dapr-distributed-calendar#services)

* [**How to run**]() (*pending*)


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

### Controller (written in Javascript)

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
    
### Services

The services handle the requests forwarded by the controller. Each of the tasks listed with the controller is handled by a service written in 
a different language. I'll detail the implementation below.

* **Event Service** (written in Go):
  This service uses the statestore component Redis for storing and deleting events from memory. The code snippet shown below is from 
  `go_events.go` and demonstrates adding an event to the state store. 

  ```go
   var data = make([]map[string]string, 1)
   data[0] = map[string]string{
    "key":   event.ID,
    "value": event.Name + " " + event.Date,
   }
   state, _ := json.Marshal(data)
   log.Printf(string(state))


   resp, err := http.Post(stateURL, "application/json", bytes.NewBuffer(state))
  ```

  where the stateURL is defined as:


  ```go
  var stateURL = fmt.Sprintf(`http://localhost:%s/v1.0/state/%s`, daprPort, stateStoreName)
  ```

* **Messaging Service** (written in Python):

  This service subscribes to the topic that we post messages to, from the controller. It then uses the [SendGrid](https://docs.dapr.io/operations/components/setup-   bindings/supported-bindings/sendgrid/) output binding to 
  send an email about creation of a new event. 
  I have used the Dapr client for Python while writing this service. 

  The code below shows how the service registers as a **subscriber** with Dapr for a specific topic.
  

  ```python
  @app.route('/dapr/subscribe', methods=['GET'])
  def subscribe():
      subscriptions = [{'pubsubname': 'pubsub',
                        'topic': 'events-topic',
                        'route': 'getmsg'}]
      return jsonify(subscriptions)
  ```
  
  > The Dapr runtime calls the `/dapr/subscribe` endpoint to register new apps as subscribers. The other way to do this would be defining a configuration
  file, linked [here](https://github.com/dapr/docs/blob/3509967baa65ece9fb822e2948e4eb7ed8d34af5/daprdocs/content/en/developing-applications/building-blocks/pubsub/howto-publish-subscribe.md#declarative-subscriptions). 
  
  The following code receives the message posted to the topic and then calls the `send_email` function.
  
  ```py
  @app.route('/getmsg', methods=['POST'])
  def subscriber():
    print(request.json, flush=True)
    
    jsonRequest = request.json
    data = jsonRequest["data"]["message"]
    print(data, flush=True)
    
    send_email()
  ```

  The send_email functions calls the SendGrid binding with the message payload:
  
  ```py
  def send_email():
    with DaprClient() as d:
            
        req_data = {
            'metadata': {
                'emailTo': emailTo,
                'subject': subject
            },
            'data': data
        }


        print(req_data, flush=True)


        # Create a typed message with content type and body
        resp = d.invoke_binding(binding_name, 'create', json.dumps(req_data))
  ```
  
  where invoke_binding is a library function from the Dapr client. In the previous cases, we had called the endpoints directly; here we 
  use a function already implemented for us.


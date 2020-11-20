# dapr-distributed-calender

This is a sample application built using Dapr as a proof-of-concept. I have experimented with the state store, pubsub and output bindings features available with Dapr.
I have used multiple languages for writing the different parts of this calender app. This demonstrates the language-agnostic nature of Dapr and the flexibility that it bings to developing
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

* Controller, written in **Javascript**: 

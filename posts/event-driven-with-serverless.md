<!--id: 5-->
<!--title: Event Driven And Serverless -->
<!--author: Brian Jones-->
<!--visible: true-->

I recently gave a talk at the Chicago Gopher's meet up about Go, the Serverless Framework, and Event Driven architecture. Here are [The slides for that talk](https://docs.google.com/presentation/d/e/2PACX-1vRzEMfp0S8ZyNwkRYC_YSfSvAULUfoU2dHFYJ_kjDnkFtJ80CH5bpffFIs6CUlMnE1-pS5_4HfoAFex/pub?start=false&loop=false&delayms=3000&slide=id.g3fa45d265d_1_5) the the code examples can be found on github at [github.com/beeceej/sls-go-examples](https://github.com/beeceej/sls-go-examples). If you weren't able to see that talk, here's a quick synopsis of some of the highlights!

## The Buzz And the Hype

Serverless has a lot of hype and buzz around it now a days. A lot of times some of the more practical aspects get lost in the noise. One joke people usually make is, "Oh, Serverless. But there _HAS_ to be servers some where!". The people who dismiss Serverless as an ideology, purely based on that notion, are missing out on what Serverless actually brings to the table! Yes! there are servers somewhere, but they are abstracted away in such a way that we just don't need to care about them. Many abstractions work great, but there are definitely times where that abstraction leaks and it won't meet a use case. That is fine. But, just because it doesn't meet every use case doesn't mean that it isn't useful. So for the sake of the hype, let's look at some situations where the abstractions work really well!

## What is Serverless

![google-cloud-functions](https://cdn-images-1.medium.com/max/1600/1*TG4VTrSkg1egeFGgzihl9Q.png)

Google Cloud says it best:

_Small_

_Single Purpose_

_Functions_

_Cloud Events_

**Small**:

- Small in size (small binary, small bundles, etc...)
- Easy to Grok, bite-sized, reusable, functionality wrapped with a well defined entry point

**Single Purpose**:

- Easy to refactor, few dependencies between disparate code paths
- Do one thing, and do it really well

**Functions**:

- React to some event
  - HTTP Request
  - Event Notification
  - Pay by use, not for existence

**Cloud Events**:

- SNS Notifications
- SQS Queues ( Recently GA’d)
- Kinesis Streams
- Dynamo Streams
- Step Functions (State Machines)
- S3 Events

## The Serverless Framework

The Serverless Framework is a tool which reduces the friction for deploy _*Serverless*_ Style applications into the cloud. Many people consider the Serverless Framework to _Be_ Serverless, but that's a common misconception. What the Serverless Framework is though, is an easy to use, declarative syntax for deploying many functions and infrastructure to the cloud at once. It uses cloud formation under the hood, so in many ways the Serverless Framework is a different take on what [Terraform](https://github.com/hashicorp/terraform) offers. It is not necessary to use the Serverless Framework to have a serverless application, but it will definitely help out a lot!

## Event Driven

Event Driven Architectures are a natural fit for Serverless applications.

In a non-event-driven architecture you may have an application architecture which looks like this:

![non-event-driven](https://static.beeceej.com/images/non-event-driven.jpg)

In this type of architecture a user will upload an image of his/her puppy, and the Puppy Profile Service will then communicate with 3 external services. To do this though, the Puppy Profile Service needs to be intimately aware of the services which it speaks to. Does it really make sense for the Puppy Profile service to know the details of Instagram, or Tinder4Dogs, or PuppyFind.com? To fix this we can move to an architecture like:

![event-driven](https://static.beeceej.com/images/event-driven.jpg)

In this architecture, the Puppy Profile Service sends an event off. And it all it has to care about is it's event. This means we've successfully decoupled the Puppy Profile Service from all of it's consumers, and can easily add more services that listen to the event! Our code gets simpler, smaller and easier to grok.

This type of architecture isn't without downsides though. Moving to an event driven pattern we lose

- Observability, when an event is fired, who is listening, did they get the message, what happens when they fail, what about logs?
- Traceability, our control flow is no longer immediately apparent by looking at the code in Puppy Profile Service

But what we gain is flexibility, these services are decoupled from each other, and the only contract between them is the event which the consumers listen to. That's powerful, because now we have a system with few dependencies. Any changes can simply be added to the event notification itself, and with some care, we can preserve both forward and backwards compatibility (that's a tough problem sometimes). Also we pay less since we only pay per execution. As an added benefit, this style of application is easy to deploy! With the Serverless Framework we can set something like this up with a .Yaml file (lots of documentation [on the serverless framework website](https://serverless.com/framework/docs/), 4 functions and a little bit of service logic. If you go with the traditional approach you will be dealing with linux instances, load balancers, complex deployments, etc... This is the power of Serverless. Quick, easy, effective and promoting good application architecture.

Essentially, Serverless allows for quick iterations, easy to reason about code, and decoupled systems.

Again, For more information refer to the slides [you can access the presentation here](https://docs.google.com/presentation/d/e/2PACX-1vRzEMfp0S8ZyNwkRYC_YSfSvAULUfoU2dHFYJ_kjDnkFtJ80CH5bpffFIs6CUlMnE1-pS5_4HfoAFex/pub?start=false&loop=false&delayms=3000&slide=id.g3fa45d265d_1_5) or the code [on github](https://github.com/sls-go-examples) and, please, don't hesitate to reach out to me at <brian.jones@beeceej.com>

Until Next time!

# WebApi Loadbalancer
This is a very simple loadbalancer for web apis that is inspired by the pluralsight course "Scaling Go Applications Horizontally". 

It exposes two endpoints, /register and /unregister on port 2002 that you call on start of your web api in order to have requests routed to it. 

It then accepts requests on port 2000 and routes these to the available services. 

The code is simple and straightforward so I encourage you to go through the commits take a look for yourself:

In [Commit 1](https://github.com/smurariu/WebApiLoadbalancer/commit/888a383aaa463cbc5af5083a1b528907e12523df) we add the basic routing of requests to applications.

In [Commit 2](https://github.com/smurariu/WebApiLoadbalancer/commit/92baad9c64401e04c4fb47665deb04df06d89389) we add the registering and unregistering endpoints so that applications can register using urls like ```http://localhost:2002/register?port=5000``` (the port needs to be specified as a query string param).

In [Commit 3](https://github.com/smurariu/WebApiLoadbalancer/commit/20cc134c735752cc016c7a5585979e778f72d3f5) the heartbeat functionality is added to handle the case when web apis go offline without unregistering or become unreachable for any other reasons.

This loadbalancer is meant to be used in concert with the [WebApiMonitoring](https://github.com/smurariu/WebApiMonitoring) midleware for ASP.NET Core WebApis but of course it can be used with many other setups.

Happy coding!
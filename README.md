SPEC:

Implement service that can handle 10/tps split between reads and writes

Endpoints:

   POST /chat
   Create a new message for the given user
   Payload {"text": <string>, "username": <string>, "timeout": <int>|60}
   return 201, with id of the message in the payload
   
   GET /chat/:username
   Return a list of all messages for the given username that:
     1.) has not expired by way of timeout (age > timeout)
     2.) has not been previously read
    
    returns:
    [{"text": <string>, "id": <id>}, {}...]

Constraint: should achieve 10/tps reads and writes

# Running the project:

Dependencies:

  docker
  docker-compose
  python2.7 


from this directory run:

`docker-compose up`

or

`docker-compose up --build` to make sure to build local changes

to run tests:

`./test.py`



Part 2:

   Scaling this service to 1000/tps should not be a particularly difficult challenge.
   With almost no optimizations or performance considerations I was able to to achieve 3-400/tps with
   one service container.

   I'd probably go with something like ECS to use AWS's autoscaling abilities, I think it should
   be pretty straightforward to apply horizontal scaling techniques with these containers.

   The stateless/ephemeral nature of the docker nodes allow them to be added to or swapped out with impunity.
   I'd simply add load balanced nodes until the desired performance was reached. 

   As far as technology goes. I'd stick with the chosen tech stack unless there were significant
   changes to the use case (or if one was provided :D)

   I think the chosen tech would work well but I may reconsider some choices if 

   a.) Entity relations became important (users, friend lists, etc, I'd consider Postgresql and/or a proper ORM)
   b.) Any siginifcant business logic became necessary for the service to function. I'm not sure that go is
   the right choice for operating at a higher level. It lacks the necessary abstractions for more complex application architecture. (classes, interhitance etc)

   As for monitoring NewRelic is my go to lately.



Note to Reviewer:

Against probably better judgement I chose to use Go/Docker for this project. This is easily the biggest project I've ever completed with either. I am pretty sure I could have been finished in a fraction of the time using a language I am more comfortable with.
For example with django/python and sqlite and `./manage.py runserver` I could meet the requirements in half an hour or so.

However, I have been looking for an excuse to fiddle with both technolgies, and this seemed appropriate. For one, using django would bring
in a whole host of unecessary dependencies and code just for 2 endpoints. 

I hope you won't hold this decision against me, (I'd certainly advise anyone doing a code challenge to use their most comfortable tech)
But my curiosity got the best of me. I'd love to discuss the details of this further and what I learned doing this. (I had a lot of fun)

All in all I'd say I spent around 4-5 hours on this, including reading through documentation and thinking about how to approach it. Perhaps 2-3 actually writing code and debugging things.
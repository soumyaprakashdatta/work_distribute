
election:
--------
- in election at most one master will be selected
- for now we can save the master's identity in a key in redis (it wont be used for now, might find some use later) 
- detail of election process is same as we follow in internal codebase
- anyone participating election will register itself in a set in redis in every cycle


master server:
-------------
- this is the server which has won election
- participates election regularly, increases ttl of election key if elected consecutive times
- master will keep pinging nodes from the registration set
- if the registration set changes a redistribution will kick in
- if master does not get response from some worker for consecutive two times, it will remove node from registration list and do a redistribution
redistribution will be done by master node and will be persisted in a redis key


worker node:
-----------
- they will participate in election regularly
- they will listen to the distribution key for any update and based on that they will process work units


when worker dies:
----------------
master's ping check will fail, so master will remove the dead node from registration list and trigger a redistribution


when master dies:
----------------
- until the next election things won't change, 
- to avoid work computation loss, we can keep redundant computation meaning same computation can happen in more than one nodes
- in the next election cycle a new master will be selected and new work distribution will happen


if redis connection goes down:
-----------------------------
we have two choices, 
 - we can stop everything and crash
 - we can continue doing things that are already assigned, 
   once redis comes back up reestablish connections and subscriptions and continue as usual
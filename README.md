This is a prototype app, it tries to do the following -

-   It has two components, redis to store state and a pool of workers.
-   This app tries to distribute work from a central workpool in redis among available workers
-   In this app, we will try to do it in an automatic manner, meaning we won't have a fixed topology and when a new worker joins or an existing worker leaves, we will redistribute work automatically
-   our goal will be to distribute work as such all workers have about the same amount of work
-   we will try to design it without the need of any other external monitor/ service discovery mechanism

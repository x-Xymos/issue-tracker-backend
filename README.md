## Build Setup

``` bash
# clone
clone the repo into a gopath i.e. /home/user/go/src/

# install the mongodb golang driver

https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial

#install mongodb on your machine
https://docs.mongodb.com/manual/installation/

The backend connects to the database using the default port
# clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

# run the backend
make -B && make run
This runs the backend as a background process and automatically kills any already running backend processes when you rebuild.

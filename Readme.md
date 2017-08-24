# Lesson Overview

The purpose of this lesson is to build a docker app and deploy it to OpenShift.

Our app will be a simple web server that says "hello".
It will optionally connect to a database to know who to say hello to.

We'll start by building our app and showing it working locally.
Then we'll deploy it to OpenShift.

# Docker

Take a look at `go/web.go`.

## Build image

This will only work if you have a docker version 17.06.0 or newer.

```bash
docker --version

git submodule update

docker build -t ccouzens/web-go-with-redis .
docker push ccouzens/web-go-with-redis
docker run --rm -p 8080:8080 ccouzens/web-go-with-redis
```

### Demonstrate it

Redis is a no-SQL database.
We can connect our app to Redis running in a container.

```bash
docker network create --driver bridge go_app_nw
docker run -d --name my_redis --network=go_app_nw redis # start the database
docker run -it --rm --network=go_app_nw redis redis-cli -h my_redis # connect to the database
SET who "Chris"
exit
docker run --rm -e REDIS_HOST_PORT=my_redis:6379 -p 8080:8080 --network=go_app_nw ccouzens/web-go-with-redis
```

# OpenShift

# Conclusion

We have learnt about deploying networked containers in both Docker and
OpenShift.

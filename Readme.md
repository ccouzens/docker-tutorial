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
To get these same containers up and running in OpenShift follow the instructions below:
First login to your OpenShift cluster:
```bash
oc login <server url> -u <username>
```
Now create the project you'll be using to deploy these containers in, here we've used the name demo
```bash
oc new-project demo
```
We use the OpenShift catalog version of redis and pass it an environment variable that sets the password automatically. Using the 'latest' tag here indicates we want to use the most recent version of the redis image.
```bash
oc new-app redis:latest -e REDIS_PASSWORD=password123
```
Next we need to deploy the Docker image we built earlier in the project - in this case its hosted on dockerhub as ccouzens/web-go-with-redis.
```bash
oc new-app ccouzens/web-go-with-redis
```

Now to prove its working we expose externally the internal service that the deployment created. This creates a route, which is an external access point to a service within a project. Once this is done we can verify the application is working now by connecting to the url listed under HOST/PORT in the get routes output.
```bash
oc expose svc/web-go-with-redis

oc get routes
NAME                HOST/PORT                                                PATH      SERVICES            PORT       TERMINATION   WILDCARD
web-go-with-redis   web-go-with-redis-demo.beta-7.cor00005.cna.ukcloud.com             web-go-with-redis   8080-tcp                 None

curl web-go-with-redis-demo.beta-7.cor00005.cna.ukcloud.com

<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Hello World</title>
</head>
<body>
<h1>Hello World</h1>
</body>
</html>
```
Now we've proved it works we're going to pass in the environment variables that make it connect to the redis database and prove connectivity is working by making a change to redis. 

Firstly we'll make the database change by getting the name of the pod that redis is running and using the oc rsh command to open a remote shell on the redis container in that pod. Then we open the redis-cli tool and make the same change to the redis container as in the previous example.
```bash
oc get pods
NAME                        READY     STATUS    RESTARTS   AGE
redis-1-53f39               1/1       Running   0          2h
web-go-with-redis-1-l9zc4   1/1       Running   0          2h

oc rsh redis-1-53f39
sh-4.2$ redis-cli -a password123
127.0.0.1:6379> SET who Universe
OK
```

Now we'll pass in the environment variables to the web app container and wait for it to redeploy with the new config. Note that when this step is completed the web app pod will have changed to a new name indicating its been redeployed.

```bash
oc env dc/web-go-with-redis -e REDIS_HOST_PORT=redis.demo.svc:6379 -e REDIS_PASSWORD=password123

oc get pods
NAME                        READY     STATUS    RESTARTS   AGE
redis-1-53f39               1/1       Running   0          2h
web-go-with-redis-2-rrblz   1/1       Running   0          32s
```

Now if we again connect to the route we created earlier we should see that our changes have taken effect and we can see our updated database value.
```bash
curl web-go-with-redis-demo.beta-7.cor00005.cna.ukcloud.com

<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Hello Universe</title>
</head>
<body>
<h1>Hello Universe</h1>
</body>
</html>
```

Of course we could've created the webapp with the correct variables immediately since we knew what the service was going to be called, the command for that would have been as follows
```bash
oc new-app ccouzens/web-go-with-redis -e REDIS_HOST_PORT=redis.demo.svc:6379 -e REDIS_PASSWORD=password123
```

# Conclusion

We have learnt about deploying networked containers in both Docker and
OpenShift.

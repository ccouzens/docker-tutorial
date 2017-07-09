# Lesson Overview

The purpose of this lesson is to create your first docker container.
We'll do better and create several.

1. An instance of Ubuntu 17.04
1. An Apache web server
1. A colleague's web server
1. A better Apache web server
1. A go-lang compiler
1. A redis database
1. A go web server

## Definitions

The difference between images and containers can be confusing.
Basically a docker container is an instance of a docker image.

### Docker Image

[Official definition](https://docs.docker.com/glossary/?term=image):

> Docker images are the basis of containers.
> An Image is an ordered collection of root filesystem changes and the corresponding execution parameters for use within a container runtime.
> An image typically contains a union of layered filesystems stacked on top of each other.
> An image does not have state and it never changes.

### Docker Container

[Official definition](https://docs.docker.com/glossary/?term=container):

> A container is a runtime instance of a docker image.
> 
> A Docker container consists of
>
> * A Docker image
> * An execution environment
> * A standard set of instructions

## Ubuntu Container

```bash
docker create --name my_ubuntu -it ubuntu:17.04
# docker ps -a

docker start -ia my_ubuntu
apt-get update
apt-get install screenfetch
screenfetch # discuss up-time
exit

# go back in and look at the history
```

* [docker create](https://docs.docker.com/engine/reference/commandline/create/)
* [docker ps](https://docs.docker.com/engine/reference/commandline/ps/)
* [docker start](https://docs.docker.com/engine/reference/commandline/start/)

## Apache container 1

```html
<html>
<head><title>My cool website</title></head>
<body>
<h1>My cool website</h1>
By Chris
</body>
</html>
```

```bash
docker login

docker start -ia my_ubuntu
apt-get install vim nano
apt-get install apache2 apache2-utils
echo > /var/www/html/index.html
nano /var/www/html/index.html
exit

docker commit -c 'CMD [ "apache2ctl", "-DFOREGROUND" ]' my_ubuntu ccouzens/apache:v1

docker run -p 8080:80 -d --name ubuntu_apache ccouzens/apache:v1
# visit http://localhost:8080/
docker rm -f ubuntu_apache

docker push ccouzens/apache:v1
```

* [docker login](https://docs.docker.com/engine/reference/commandline/login/)
* [docker commit](https://docs.docker.com/engine/reference/commandline/commit/)
* [docker run](https://docs.docker.com/engine/reference/commandline/run/)
* [docker rm](https://docs.docker.com/engine/reference/commandline/rm/)
* [docker push](https://docs.docker.com/engine/reference/commandline/push/)

## Apache Container 2 (colleague)

```bash
docker run -p 8080:80 -d --name ubuntu_apache colleague/apache:v1
# visit http://localhost:8080/
docker rm -f ubuntu_apache
```

## Apache Container 3 (better image)

Our Apache docker image has a couple shortcomings:

1. The build steps are hard to automate
1. We're including Ubuntu

```html
<html>
<head><title>My even cooler website</title></head>
<body>
<h1>My even cooler website</h1>
By Chris
</body>
</html>
```

```dockerfile
FROM httpd:2.4
COPY *.html /usr/local/apache2/htdocs/
```

Save the html file as index.html.
Save the dockerfile as dockerfile.

```bash
docker build -t ccouzens/apache:v2 .
docker run --rm -p 8080:80 ccouzens/apache:v2
# visit http://localhost:8080/ then exit
^C

# optionally push and share
```

## Go Compiler

Take a look at `go/web.go`.
Most of you won't have go installed on your laptops.
So how do we compile it?

**Answer**: [use docker](https://hub.docker.com/_/golang/)!

```bash
# Linux host
docker run --rm -v "$PWD":/go/src/web -w /go/src/web golang:1.8 go build -v
./web # open http://localhost:8080/

# Mac host
docker run --rm -v "$PWD":/go/src/web -w /go/src/web -e GOOS=darwin golang:1.8 go build -v
./web # open http://localhost:8080/

# Windows host
docker run --rm -v "$PWD":/go/src/web -w /go/src/web -e GOOS=windows golang:1.8 go build -v
./web.exe # open http://localhost:8080/
```

Docker containers can be short lived tools, not just servers.

## Redis

Redis is a no-SQL database.

```bash
docker run -d --name my_redis -p 6379:6379 redis # start the database
docker run -it --rm --link my_redis redis redis-cli -h my_redis # connect to the database
SET who "Chris"
exit

REDIS_HOST_PORT=localhost:6379 ./web # open http://localhost:8080/
```
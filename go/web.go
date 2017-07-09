package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
)

var templ = template.Must(template.New("hello").Parse(templateStr))
var client *redis.Client

func who() string {
	val, err := client.Get("who").Result()
	if err != nil {
		log.Print(err)
		val = "World"
	}
	return val
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	templ.Execute(w, who())
}

func main() {
	client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST_PORT"),
		Password: "",
		DB:       0,
	})
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

const templateStr = `
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Hello {{.}}</title>
</head>
<body>
<h1>Hello {{.}}</h1>
</body>
</html>
`

package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/leon123858/terroform-sample/docker-http/pubsub"
)

var methodList = []string{
	"get-todo-list",
	"add-todo-item",
	"remove-todo-item",
	"update-todo-item",
}

func main() {
	e := echo.New()

	// init pubsub client
	pubsubInfo, err := pubsub.NewPubSub("tw-rd-ca-leon-lin")
	if err != nil {
		panic(err)
	}

	// Enable CORS
	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/call/:method", func(c echo.Context) error {
		method := c.Param("method")
		if !pubsub.Contains(methodList, method) {
			return c.JSON(http.StatusBadRequest, map[string]string{"status": "fail"})
		}

		req := new(pubsub.Request)
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"status": "fail"})
		}

		// print json
		fmt.Printf("%+v\n", req)

		err = pubsubInfo.Publish2Topic(method, *req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"status": "fail"})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "success"})
	})

	e.Logger.Fatal(e.Start(":8080"))
}

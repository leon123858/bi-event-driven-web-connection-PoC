package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// key is channel id, value is todo list
// sample data: {"1": ["todo1", "todo2"]}
var mockDb = map[string][]string{}

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/all", func(c echo.Context) error {
		return c.JSON(http.StatusOK, mockDb)
	})

	// add todo to channel,
	// query sample: http://localhost:8080/add?channelId=1&todoMsg=todo1
	e.GET("/add", func(c echo.Context) error {
		channelId := c.QueryParam("channelId")
		todoMsg := c.QueryParam("todoMsg")
		mockDb[channelId] = append(mockDb[channelId], todoMsg)
		return c.JSON(http.StatusOK, map[string]string{"status": "success"})
	})

	e.Logger.Fatal(e.Start(":8080"))
}

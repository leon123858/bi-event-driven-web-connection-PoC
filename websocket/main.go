package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/websocket"

	"github.com/leon123858/bi-event-driven-web-connection-PoC/websocket/pubsub"
)

type NoticeRequest struct {
	UserId string `json:"userId"`
}

type NoticeResponse struct {
	ChannelId int64  `json:"channelId"`
	Msg       string `json:"msg"`
}

var projectID = "tw-rd-ca-leon-lin"
var channelId int64
var noticesChannel = make(map[string](chan pubsub.Notice))

func connect(ws *websocket.Conn) {
	// ws get the first image after connect
	var in NoticeRequest
	if err := websocket.JSON.Receive(ws, &in); err != nil {
		println("Error First Receive", err.Error())
		return
	}
	if err := websocket.JSON.Send(ws, &NoticeResponse{ChannelId: channelId, Msg: ""}); err != nil {
		println("Error First Send: ", err.Error())
		return
	}
	// create a channel for this user
	(noticesChannel)[in.UserId] = make(chan pubsub.Notice)
	defer func() {
		close((noticesChannel)[in.UserId])
		delete(noticesChannel, in.UserId)
	}()
	for {
		// get message from channel
		msg := <-(noticesChannel)[in.UserId]
		// send message to client
		if err := websocket.JSON.Send(ws, NoticeResponse{ChannelId: channelId, Msg: msg.Msg}); err != nil {
			println("Error Send: ", err.Error())
			return
		}
	}
}

func main() {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "1234"
	}

	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	pubSubInfo := pubsub.PubSubInfo{}
	err := pubSubInfo.Init(projectID)
	if err != nil {
		panic(err)
	}
	channelId = pubSubInfo.ChannelId

	go func(p pubsub.PubSubInfo, c *map[string](chan pubsub.Notice)) {
		pubsub.PullMsgs(p, c) // pass the address of the c variable
	}(pubSubInfo, &noticesChannel)
	go func() {
		<-cancelSignal
		println("start to release resources")
		err := pubSubInfo.Release()
		if err != nil {
			panic(err)
		}
		println("end to release resources")
		// stop ListenAndServe when get cancelSignal
		os.Exit(0)
	}()
	// init websocket
	http.Handle("/", websocket.Handler(connect))

	println("Server started: on port " + PORT)
	if err := http.ListenAndServe(":"+PORT, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

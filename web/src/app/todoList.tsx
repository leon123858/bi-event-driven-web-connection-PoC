"use client";
import useWebSocket from "react-use-websocket";
import styles from "./page.module.css";
import { useState } from "react";
import { EventEmitter, EventSubscription } from "fbemitter";

const socketUrl = "wss://websocket-kt6w747drq-de.a.run.app";
const httpUrl = "https://http-kt6w747drq-de.a.run.app/call";

export default function TodoList() {
  let emitter = new EventEmitter();
  const [list, setList] = useState<string[]>([]);
  const [input, setInput] = useState<string>("");
  const [listName, setListName] = useState<string>("test");
  const [userId, setUserId] = useState<string>("default");
  const [channelId, setChannelId] = useState<string>("");
  const {
    sendMessage,
    sendJsonMessage,
    lastMessage,
    lastJsonMessage,
    readyState,
    getWebSocket,
  } = useWebSocket(socketUrl, {
    onOpen: () => {
      console.log("opened");
      sendJsonMessage({ userId: userId });
    },
    onClose: (e) => console.log("closed: " + e.reason),
    onMessage: (e) => {
      const data: {
        channelId: string;
        msg: string;
      } = JSON.parse(e.data);
      setChannelId(data.channelId);
      // console.log(data);
      if (data.msg === "") {
        const body = {
          name: listName,
          channelId: data.channelId,
          userId: userId,
        };
        // get list http request
        fetch(`${httpUrl}/get-todo-list`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(body),
        })
          .then((response) => response.json())
          .then((data) => {
            console.log("Success:", data);
          })
          .catch((error) => {
            console.error("Error:", error);
          });
        return;
      } else {
        const msg = JSON.parse(data.msg);
        console.log(msg.data);
        switch (msg.type) {
          case "getTodoList":
            setList(msg.data.map((v: any) => v.Name + ": " + v.Description));
            break;
          case "addTodoItem":
            if (msg.result) {
              emitter.emit("add-item", "");
            } else {
              alert("add item failed");
            }
            break;
        }
      }
    },
    //Will attempt to reconnect on all close events, such as server shutting down
    shouldReconnect: (closeEvent) => true,
  });

  return (
    <>
      <input
        value={listName}
        placeholder="list name"
        onChange={(e) => setListName(e.target.value)}
      ></input>
      <input
        value={userId}
        placeholder="user id"
        onChange={(e) => setUserId(e.target.value)}
      ></input>
      <div className={styles.description}>
        <input
          value={input}
          placeholder="item content"
          onChange={(e) => setInput(e.target.value)}
        ></input>
        <button
          onClick={async () => {
            // create item http request
            if (channelId !== "") {
              fetch(`${httpUrl}/add-todo-item`, {
                method: "POST",
                headers: {
                  "Content-Type": "application/json",
                },
                body: JSON.stringify({
                  name: listName,
                  channelId: channelId,
                  userId: userId,
                  description: input,
                  completed: false,
                }),
              })
                .then((response) => response.json())
                .then((data) => {
                  let token: EventSubscription;
                  const t = setTimeout(() => {
                    token.remove(); // 5 seconds later, remove the listener
                    alert("add item timeout");
                  }, 5000);
                  token = emitter.once("add-item", (data: any) => {
                    clearTimeout(t);
                    alert("add item success");
                    setList((list) => [`${listName}: ${input}`, ...list]);
                  });
                })
                .catch((error) => {
                  console.error("Error:", error);
                });
            } else {
              alert("should not connect to server");
            }
          }}
        >
          clink me to add item
        </button>
        <div>
          {list.map((v, i) => {
            return <p key={i}>{v}</p>;
          })}
        </div>
      </div>
    </>
  );
}

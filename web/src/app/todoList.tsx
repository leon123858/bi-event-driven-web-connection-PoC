"use client";
import useWebSocket from "react-use-websocket";
import styles from "./page.module.css";
import { useState } from "react";

const socketUrl = "ws://127.0.0.1:1234";

export default function TodoList() {
  const [list, setList] = useState<string[]>([]);
  const [input, setInput] = useState<string>("");
  const {
    sendMessage,
    sendJsonMessage,
    lastMessage,
    lastJsonMessage,
    readyState,
    getWebSocket,
  } = useWebSocket(socketUrl, {
    onOpen: () => console.log("opened"),
    onClose: (e) => console.log("closed: " + e.reason),
    onMessage: (e) => {
      console.log("message: " + e.data);
      setList([...list, e.data]);
    },
    //Will attempt to reconnect on all close events, such as server shutting down
    shouldReconnect: (closeEvent) => true,
  });

  return (
    <div className={styles.description}>
      <input value={input} onChange={(e) => setInput(e.target.value)}></input>
      <button
        onClick={async () => {
          sendMessage(input || "add item");
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
  );
}

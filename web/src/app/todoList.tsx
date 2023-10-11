"use client";
import useWebSocket from "react-use-websocket";
import { useState } from "react";
import Container from "@/components/MiddleContainer";
import ConfigForm from "@/components/ConfigForm";
import React from "react";
import {
  TodoItem,
  addTodoItem,
  editTodoItem,
  getTodoList,
  removeTodoItem,
} from "@/utils/httpRequester";
import useConfigHook from "@/utils/configHook";
import {
  dispatcher,
  listener,
  listenerWithTimeout,
} from "@/utils/socketDispatcher";
import { Button, Checkbox, Input, List, Typography } from "antd";

const defaultProp = {
  customBackendUrl: "default" as "default" | "custom",
  listName: "test",
  userName: "default",
  socketUrl: "wss://websocket-kt6w747drq-de.a.run.app",
  httpUrl: "https://http-kt6w747drq-de.a.run.app/call",
};

export default function TodoList() {
  const [list, setList] = useState<TodoItem[]>([]);
  const [input, setInput] = useState<string>("");
  const {
    socketUrl,
    basicData,
    channelId,
    setBasicData,
    setSocketUrl,
    setChannelId,
  } = useConfigHook(defaultProp);

  listener("get-item", (data: TodoItem[]) => {
    if (data)
      setList([
        ...data.map((v) => {
          return {
            Completed: v.Completed,
            Description: v.Description,
            ID: v.ID,
            Name: v.Name,
          };
        }),
      ]);
  });

  const { sendJsonMessage } = useWebSocket(socketUrl, {
    onOpen: () => {
      console.log("opened");
      sendJsonMessage({ userId: basicData.userName });
    },
    onClose: (e) => console.log("closed: " + e.reason),
    onMessage: async (e) => {
      const data: {
        channelId: string;
        msg: string;
      } = JSON.parse(e.data);
      setChannelId(data.channelId);
      if (data.msg === "") {
        // first round response protocol
        await getTodoList(data.channelId, basicData);
        return;
      } else {
        dispatcher(data.msg);
      }
    },
    //Will attempt to reconnect on all close events, such as server shutting down
    shouldReconnect: (closeEvent) => true,
  });

  return (
    <>
      <Container style={{ width: "100%", height: "10vh" }}>
        <h1>High Concurrency TodoList Demo</h1>
      </Container>

      <Container style={{ width: "100%", height: "50vh" }}>
        <ConfigForm
          defaultValue={defaultProp}
          connectCallback={(values: any) => {
            setBasicData({
              listName: values.listName,
              userName: values.userName,
              httpUrl: values.httpUrl,
            });
            setSocketUrl(values.socketUrl);
          }}
        ></ConfigForm>
      </Container>

      <Container style={{ width: "100%" }}>
        <List
          header={
            <>
              <Input
                style={{ backgroundColor: "gray" }}
                value={input}
                onChange={(e) => setInput(e.currentTarget.value)}
              ></Input>
              <Button
                style={{ width: "100%" }}
                onClick={async () => {
                  await addTodoItem(channelId, basicData, input, () => {
                    const tmpId = Math.random().toString(36).substr(2, 9);
                    setInput("");
                    setList([
                      {
                        ID: tmpId,
                        Name: basicData.listName,
                        Description: input,
                        Completed: false,
                      },
                      ...list,
                    ]);
                    listenerWithTimeout(
                      "add-item",
                      () => {
                        setList([...list.filter((item) => item.ID !== tmpId)]);
                        console.error("add item timeout");
                      },
                      (id: string) => {
                        setList([
                          {
                            ID: id,
                            Name: basicData.listName,
                            Description: input,
                            Completed: false,
                          },
                          ...list.filter((item) => item.ID !== tmpId),
                        ]);
                        console.log("add item success");
                      }
                    );
                  });
                }}
              >
                Add Item
              </Button>
            </>
          }
          bordered
          dataSource={list}
          renderItem={(item) => (
            <List.Item style={{ color: "wheat", border: "1px solid gray" }}>
              <List.Item.Meta
                avatar={
                  <Checkbox
                    checked={item.Completed}
                    onClick={async () => {
                      await editTodoItem(
                        item.ID,
                        channelId,
                        basicData,
                        item.Description,
                        !item.Completed,
                        async () => {
                          setList([
                            ...list.map((v) => {
                              if (v.ID === item.ID) {
                                v.Completed = !item.Completed;
                              }
                              return v;
                            }),
                          ]);
                          listenerWithTimeout(
                            "edit-item",
                            () => {
                              setList([
                                ...list.map((v) => {
                                  if (v.ID === item.ID) {
                                    v.Completed = !v.Completed;
                                  }
                                  return v;
                                }),
                              ]);
                              console.error("edit item timeout");
                            },
                            () => {
                              console.log("edit item success");
                            }
                          );
                        }
                      );
                    }}
                  ></Checkbox>
                }
                title={<p style={{ color: "wheat" }}>{item.Description}</p>}
                description={
                  <Button
                    onClick={async () => {
                      await removeTodoItem(
                        item.ID,
                        channelId,
                        basicData,
                        () => {
                          const originList = [...list];
                          setList([...list.filter((v) => v.ID !== item.ID)]);
                          listenerWithTimeout(
                            "remove-item",
                            () => {
                              setList(originList);
                              console.error("remove item timeout");
                            },
                            () => {
                              console.log("remove item success");
                            }
                          );
                        }
                      );
                    }}
                    type="text"
                    style={{ color: "red", border: "1px red solid" }}
                  >
                    Delete
                  </Button>
                }
              />
            </List.Item>
          )}
        />
      </Container>
    </>
  );
}

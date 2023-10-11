import { EventEmitter, EventSubscription } from "fbemitter";

const emitter = new EventEmitter();

const dispatcher = (message: string) => {
  const data = JSON.parse(message);
  switch (data.type) {
    case "getTodoList":
      emitter.emit("get-item", data.data);
      break;
    case "addTodoItem":
      emitter.emit("add-item", data.data);
      break;
    case "editTodoItem":
      emitter.emit("edit-item", data.data);
      break;
    default:
      alert("Error: unknown message type" + data.type);
      break;
  }
};

const listener = (event: string, callback: Function) => {
  emitter.addListener(event, callback);
};

const listenerWithTimeout = (
  event: string,
  ClearCallback: Function,
  SuccessCallback: Function,
  timeout = 10000
) => {
  let token: EventSubscription;
  const t = setTimeout(() => {
    token.remove(); // remove the listener
    // remove item in list where id is tmpId
    ClearCallback();
  }, timeout);
  token = emitter.once(event, (data: any) => {
    clearTimeout(t);
    // edit item in list where id is tmpId
    SuccessCallback(data);
  });
};

export { dispatcher, emitter, listener, listenerWithTimeout };

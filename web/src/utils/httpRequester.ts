export interface TodoItem {
  Completed: boolean;
  Description: string;
  ID: string;
  Name: string;
}

export interface ConfigProps {
  customBackendUrl: "default" | "custom";
  listName: string;
  userName: string;
  socketUrl: string;
  httpUrl: string;
}

export const getTodoList = async (
  channelId: string,
  basicData: { listName: string; userName: string; httpUrl: string }
) => {
  if (channelId === "") {
    alert("Error: Please get a channel first");
    return;
  }
  // get list http request
  await fetch(`${basicData.httpUrl}/get-todolist-topic`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      name: basicData.listName,
      channelId: channelId,
      userId: basicData.userName,
    }),
  })
    .then((response) => response.json())
    .then((_) => {})
    .catch((error) => {
      console.error("Error:", error);
      alert("Error: Please check your network connection, http request failed");
    });
};

export const addTodoItem = async (
  channelId: string,
  basicData: { listName: string; userName: string; httpUrl: string },
  input: string,
  callback: () => void
) => {
  const data = {
    name: basicData.listName,
    channelId: channelId,
    userId: basicData.userName,
    description: input,
    completed: false,
  };
  await fetch(`${basicData.httpUrl}/create-todolist-topic`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  })
    .then((response) => response.json())
    .then((_) => callback())
    .catch((error) => {
      console.error("Error:", error);
      alert("Error: Please check your network connection, http request failed");
    });
};

export const editTodoItem = async (
  targetId: string,
  channelId: string,
  basicData: { listName: string; userName: string; httpUrl: string },
  input: string,
  completed: boolean,
  callback: () => void
) => {
  const data = {
    id: targetId,
    channelId: channelId,
    userId: basicData.userName,
    description: input,
    completed: completed,
  };
  await fetch(`${basicData.httpUrl}/update-todolist-topic`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  })
    .then((response) => response.json())
    .then((_) => callback())
    .catch((error) => {
      console.error("Error:", error);
      alert("Error: Please check your network connection, http request failed");
    });
};

export const removeTodoItem = async (
  targetId: string,
  channelId: string,
  basicData: { listName: string; userName: string; httpUrl: string },
  callback: () => void
) => {
  const data = {
    id: targetId,
    channelId: channelId,
    userId: basicData.userName,
  };
  await fetch(`${basicData.httpUrl}/delete-todolist-topic`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  })
    .then((response) => response.json())
    .then((_) => callback())
    .catch((error) => {
      console.error("Error:", error);
      alert("Error: Please check your network connection, http request failed");
    });
};

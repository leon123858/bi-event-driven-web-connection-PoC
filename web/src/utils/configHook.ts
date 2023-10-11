import { useState } from "react";
import { ConfigProps } from "./httpRequester";

const useConfigHook = (defaultData: ConfigProps) => {
  const [socketUrl, setSocket] = useState(defaultData.socketUrl);
  const [basicData, setBasic] = useState({
    listName: defaultData.listName,
    userName: defaultData.userName,
    httpUrl: defaultData.httpUrl,
  });
  const [channelId, setChannelId] = useState("");
  const setSocketUrl = (url: string) => {
    setSocket(url || defaultData.socketUrl);
  };
  const setBasicData = (data: {
    listName: string;
    userName: string;
    httpUrl: string;
  }) => {
    if (!data.httpUrl) setBasic({ ...data, httpUrl: defaultData.httpUrl });
  };

  return {
    socketUrl,
    basicData,
    channelId,
    setSocketUrl,
    setBasicData,
    setChannelId,
  };
};

export default useConfigHook;

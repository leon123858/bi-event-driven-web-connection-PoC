import React from "react";
import { Flex } from "antd";

interface ContainerProps {
  style?: React.CSSProperties;
  children?: React.ReactNode;
}

const Container: React.FC<ContainerProps> = ({ style, children }) => {
  return (
    <Flex gap="middle" align="start" vertical>
      <Flex style={{ ...style }} justify={"center"} align={"center"}>
        {children}
      </Flex>
    </Flex>
  );
};

export default Container;

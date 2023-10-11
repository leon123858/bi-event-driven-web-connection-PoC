import React, { useEffect } from "react";
import { Button, Form, Input, Select } from "antd";
import type { FormInstance } from "antd/es/form";
import { ConfigProps } from "@/utils/httpRequester";

const { Option } = Select;

const layout = {
  labelCol: { span: 8 },
  wrapperCol: { span: 16 },
};

const tailLayout = {
  wrapperCol: {
    offset: 8,
    span: 16,
  },
};

interface ContainerProps {
  defaultValue: ConfigProps;
  connectCallback: (prop: any) => void;
}

const ConfigForm: React.FC<ContainerProps> = ({
  defaultValue,
  connectCallback,
}) => {
  const formRef = React.useRef<FormInstance>(null);
  useEffect(() => {
    formRef.current?.setFieldsValue(defaultValue);
  }, []);

  return (
    <Form
      {...layout}
      ref={formRef}
      onFinish={(values) => {
        connectCallback(values);
      }}
      style={{
        width: 600,
        backgroundColor: "white",
        padding: "15px",
        borderRadius: "5px",
      }}
    >
      <h1 style={{ paddingLeft: "5px" }}>Config</h1>
      <br></br>
      <Form.Item
        name={"customBackendUrl"}
        label={"CustomBackendUrl"}
        rules={[{ required: true }]}
      >
        <Select placeholder="Select custom when you want to use own backend">
          <Option value="custom">Custom</Option>
          <Option value="default">Default</Option>
        </Select>
      </Form.Item>
      <Form.Item
        noStyle
        shouldUpdate={(prevValues, currentValues) =>
          prevValues.customBackendUrl !== currentValues.customBackendUrl
        }
      >
        {({ getFieldValue }) =>
          getFieldValue("customBackendUrl") === "custom" ? (
            <Form.Item
              name="socketUrl"
              label="SocketUrl"
              rules={[{ required: true }]}
            >
              <Input />
            </Form.Item>
          ) : null
        }
      </Form.Item>
      <Form.Item
        noStyle
        shouldUpdate={(prevValues, currentValues) =>
          prevValues.customBackendUrl !== currentValues.customBackendUrl
        }
      >
        {({ getFieldValue }) =>
          getFieldValue("customBackendUrl") === "custom" ? (
            <Form.Item
              name="httpUrl"
              label="HttpUrl"
              rules={[{ required: true }]}
            >
              <Input />
            </Form.Item>
          ) : null
        }
      </Form.Item>
      <Form.Item name="listName" label="ListName" rules={[{ required: true }]}>
        <Input></Input>
      </Form.Item>
      <Form.Item name="userName" label="UserName" rules={[{ required: true }]}>
        <Input></Input>
      </Form.Item>
      <Form.Item {...tailLayout}>
        <Button type="primary" htmlType="submit">
          Connect
        </Button>
      </Form.Item>
    </Form>
  );
};

export default ConfigForm;

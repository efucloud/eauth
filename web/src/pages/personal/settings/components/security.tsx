import { List, Row, Col } from 'antd';
import React, { useState, useEffect, useRef } from 'react';
import { FormattedMessage, useIntl } from '@umijs/max';
import { LockOutlined } from '@ant-design/icons';
import { ModalForm, ProFormText } from '@ant-design/pro-components';
import { AuthedUserInfo } from "@/services/common.d";
import { systemGetUserinfo, updateSelfInfo } from '@/services/oauth.api';
import { setPassword } from '@/services/personal.api';
import { SetPassword } from '@/services/user';
import type { ProFormInstance } from '@ant-design/pro-components';
import { getColorPrimary } from '@/utils/global';



type Unpacked<T> = T extends (infer U)[] ? U : T;

const passwordStrength = (strength: string) => {
  if (strength === 'weak') {
    return <FormattedMessage id="model.user.password.weak" />;
  } else if (strength === 'medium') {
    return <FormattedMessage id="model.user.password.medium" />;
  } else if (strength === 'strong') {
    return <FormattedMessage id="model.user.password.strong" />;
  }
  return <FormattedMessage id="model.user.password.weak" />;
};



const SecurityView: React.FC = () => {
  const [currentUser, setCurrentUser] = useState<AuthedUserInfo>();
  const formRef = useRef<ProFormInstance>();
  const formRef1 = useRef<ProFormInstance>();
  const colorPrimary = getColorPrimary();
  const intl = useIntl();
  const [passwordModalVisible, setPasswordModalVisible] = React.useState(false);
  const [passwordSetModalVisible, setPasswordSetModalVisible] = React.useState(false);
  const getCurrentUser = async () => {
    const data = await systemGetUserinfo();
    setCurrentUser(data as AuthedUserInfo);
  }
  useEffect(() => { getCurrentUser(); }, []);
  const getData = () => {
    if (currentUser) {
      let list = [
        {
          title: <FormattedMessage id="model.user.phone" />,
          description: <>{currentUser.phone}</>,
          actions: [<a style={{color:colorPrimary}} key="modify"><FormattedMessage id='pages.operation.edit' /></a>],
        },
        {
          title: <FormattedMessage id="model.user.email" />,
          description:  <>{currentUser.email}</>,
          actions: [<a style={{color:colorPrimary}} key="modify"><FormattedMessage id='pages.operation.edit' /></a>],
        },
      ];
      if (currentUser.hasPassword) {
        list.unshift({
          title: <FormattedMessage id="model.user.password" />,
          description: <><FormattedMessage id='model.user.password.passwordStrength' />:  {passwordStrength(currentUser.passwordStrength || '')} </>,
          actions: [<a style={{color:colorPrimary}} key="modify" onClick={() => setPasswordSetModalVisible(true)}><FormattedMessage id='pages.operation.edit' /></a>],
        })
      } else {
        list.unshift({
          title: <FormattedMessage id="model.user.password" />,
          description: <FormattedMessage id='model.user.password.notSet' />,
          actions: [<a style={{color:colorPrimary}} key="modify" onClick={() => setPasswordModalVisible(true)}><FormattedMessage id='pages.operation.set' /></a>],
        })
      }
      return list;
    };
    return [];
  }
  const data = getData();
  const onSetPassword = async (values: Record<string, any>) => {
    const res = await setPassword(values as SetPassword);
    setPasswordModalVisible(false);
    setPasswordSetModalVisible(false);
  };
  const checkPassword = (rule: any, value: any, callback: any) => {
    const newPassword = formRef.current?.getFieldValue('newPassword');
    if (value && value !== newPassword) {
      callback(intl.formatMessage({ id: 'pages.login.resetPassword.checkPassword' }));
    } else {
      callback();
    }
  };
  const checkPassword1 = (rule: any, value: any, callback: any) => {
    const newPassword = formRef1.current?.getFieldValue('newPassword');
    if (value && value !== newPassword) {
      callback(intl.formatMessage({ id: 'pages.login.resetPassword.checkPassword' }));
    } else {
      callback();
    }
  };
  return (
    <>
      <Row gutter={{ xs: 16, sm: 32, md: 64, lg: 128 }}>
        <Col lg={16} md={16} sm={24} >
          <List<Unpacked<typeof data>>
            itemLayout="horizontal"
            dataSource={data}
            renderItem={(item) => (
              <List.Item actions={item.actions}>
                <List.Item.Meta title={item.title} description={item.description} />
              </List.Item>
            )}
          />
        </Col>
      </Row>
      <ModalForm
        formRef={formRef1}
        title={intl.formatMessage({ id: 'model.user.password.set' })}
        width="400px"
        open={passwordModalVisible}
        onOpenChange={setPasswordModalVisible}
        onFinish={onSetPassword}
      >
        <ProFormText.Password
          label={intl.formatMessage({ id: 'model.user.password.new' })}
          name="newPassword"
          fieldProps={{ size: 'large', prefix: <LockOutlined style={{color:colorPrimary}}/> }}
          rules={[
            {
              required: true,
              message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.user.password.new' }),
            },
          ]}
        />
        <ProFormText.Password
          label={intl.formatMessage({ id: 'model.user.password.confirm' })}
          name="confirmPassword"
          fieldProps={{ size: 'large', prefix: <LockOutlined style={{color:colorPrimary}}/> }}
          rules={[
            { validator: checkPassword1 },
            {
              required: true,
              message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.user.password.confirm' }),
            },
          ]}
        />
      </ModalForm>
      <ModalForm
        formRef={formRef}
        title={intl.formatMessage({ id: 'model.user.password.set' })}
        width="400px"
        open={passwordSetModalVisible}
        onOpenChange={setPasswordSetModalVisible}
        onFinish={onSetPassword}
      >
        <ProFormText.Password
          label={intl.formatMessage({ id: 'model.user.password.old' })}
          name="oldPassword"
          fieldProps={{ size: 'large', prefix: <LockOutlined style={{color:colorPrimary}}/> }}
          rules={[
            {
              required: true,
              message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.user.password.old' }),
            },
          ]}
        />
        <ProFormText.Password
          label={intl.formatMessage({ id: 'model.user.password.new' })}
          name="newPassword"
          fieldProps={{ size: 'large', prefix: <LockOutlined style={{color:colorPrimary}}/> }}
          rules={[
            {
              required: true,
              message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.user.password.new' }),
            },
          ]}
        />
        <ProFormText.Password
          label={intl.formatMessage({ id: 'model.user.password.confirm' })}
          name="confirmPassword"
          fieldProps={{ size: 'large', prefix: <LockOutlined style={{color:colorPrimary}}/> }}
          rules={[
            { validator: checkPassword },
            {
              required: true,
              message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.user.password.confirm' }),
            },
          ]}
        />
      </ModalForm>
    </>
  );
};

export default SecurityView;

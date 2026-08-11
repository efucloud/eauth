import type { ProFormInstance } from '@ant-design/pro-components';
import {
  PageContainer,
  ProForm,
  ProFormText,
  StepsForm,
} from '@ant-design/pro-components';
import { Card, Typography, Avatar, Row, Col,message } from 'antd';
import React, { useState, useEffect, useRef } from "react";
import { useIntl, useParams, FormattedMessage } from '@umijs/max';
import { LockOutlined } from '@ant-design/icons';
const { Text, Paragraph, Title } = Typography;
import { systemValidateUserInfo, systemUserResetPassword } from '@/services/security.api';
import { UserResetPassword } from '@/services/user';
import { ShortUser } from '@/services/user';

const ForgetPassword: React.FC = () => {
  const intl = useIntl();
  const formRef = useRef<ProFormInstance>();
  const code = useParams().code;
  const [info, setInfo] = useState<ShortUser>();
  const getUserInfo = async () => {
    const res = await systemValidateUserInfo({ code: code || '' });
    setInfo(res as ShortUser);
  };
  useEffect(() => {
    getUserInfo();
  }, [code]);
  const checkPassword = (rule: any, value: any, callback: any) => {
    const newPassword = formRef.current?.getFieldValue('password');
    if (value && value !== newPassword) {
      callback(intl.formatMessage({ id: 'pages.login.resetPassword.checkPassword' }));
    } else {
      callback();
    }
  };
  const onFinish = async (values: Record<string, any>) => {
    values['code'] = code;
    await systemUserResetPassword(values as UserResetPassword);
    message.success(intl.formatMessage({ id: 'pages.resetPassword.success' }));
    window.location.href = '/user/login';
  };
  return (
    <PageContainer title={false}
    ><div style={{ margin: 'center' }}>
        <Card
          style={{ backgroundColor: '#fff', width: '50%', margin: '0 auto' }}
        >
          <div style={{ textAlign: 'center', marginBottom: 10 }} >
            <Title level={2}> <Avatar size={64} src='/logo.png' />&nbsp;&nbsp;<FormattedMessage id='app.title' /></Title>
            <Title level={4}><FormattedMessage id='model.system.user' />&nbsp;&nbsp; <Avatar size={40} src={info?.avatar} />&nbsp;&nbsp;{info?.username}</Title>
          </div>
          <Row>
            <Col offset={6} lg={12} md={12} sm={24}>
              <ProForm
                formRef={formRef}
                onFinish={onFinish}
              >
                <ProFormText.Password
                  label={intl.formatMessage({ id: 'model.user.password.new' })}
                  name="password"
                  fieldProps={{ size: 'large', prefix: <LockOutlined /> }}
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
                  fieldProps={{ size: 'large', prefix: <LockOutlined /> }}
                  dependencies={['password']}
                  rules={[
                    { validator: checkPassword },
                    {
                      required: true,
                      message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.user.password.confirm' }),
                    },
                  ]}
                />
              </ProForm>
            </Col></Row>
        </Card></div>

    </PageContainer>
  );
};

export default ForgetPassword;
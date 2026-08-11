import type { ProFormInstance } from '@ant-design/pro-components';
import {
  PageContainer,
  ProFormText,
  ProForm,
} from '@ant-design/pro-components';
import { MailOutlined } from '@ant-design/icons';
import { message,  Card, Typography, Row, Col, Button, Avatar } from 'antd';
const { Text, Paragraph, Title } = Typography;
import React, { useState, useRef } from "react";
import { FormattedMessage, useIntl } from '@umijs/max';
import { systemEmailExist,sendResetPwdEmail } from '@/services/security.api';
import { ExistResponse } from '@/services/common';

const ForgetPassword: React.FC = () => {
  const intl = useIntl();
  const formRef = useRef<ProFormInstance>();
  const [checkSuccess, setCheckSuccess] = useState<boolean>(false);
  const checkEmail = async (email: string) => {
    const res = await systemEmailExist({ email }) as ExistResponse;
    if (res.exist) {
      return true;
    }
    return false;
  };
  const handleSubmit = async () => {
    const email = formRef.current?.getFieldValue('email') as string;
    if (!email) {
      return;
    }
    const exist = await checkEmail(email);
    if (!exist) {
      message.error(intl.formatMessage({ id: 'model.user.email.not.exist' }));
      return;
    }
    const res = await sendResetPwdEmail({ name: email } as ExistResponse);
    if (res === 'success') {
      setCheckSuccess(true);
    } else {
      message.error(intl.formatMessage({ id: 'model.user.email.send.fail' }));
    }
  };
  return (
    <PageContainer title={false}
    >
      <div style={{ margin: 'center' }}>
        <Card
          style={{ backgroundColor: '#fff', width: '50%', margin: '0 auto' }}
        >
          <div style={{ textAlign: 'center',marginBottom:10 }} >
          <Title level={2}> <Avatar size={64} src='/logo.png'/>&nbsp;&nbsp;<FormattedMessage id='app.title'/></Title></div>
          {checkSuccess === false && <><Paragraph style={{ textAlign: 'center' }}>
            <Title level={4}><FormattedMessage id='model.user.forget.password.reset' /></Title>
            <Text><FormattedMessage id='model.user.forget.password.reset.description' /></Text>
            <br />
            <br />
          </Paragraph>
            <Row>
              <Col offset={6} lg={12} md={12} sm={24}>
                <ProForm
                  layout="vertical"
                  submitter={false}
                  formRef={formRef}
                >
                  <ProFormText
                    name="email"
                    fieldProps={{ size: 'large', prefix: <MailOutlined /> }}
                    placeholder={intl.formatMessage({ id: 'model.user.email' })}
                    rules={[
                      {
                        required: true,
                        type: 'email',
                        message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.user.email' }),
                      },
                    ]}
                  />
                  <Button type="primary" htmlType="submit" size="large" style={{ width: '100%' }}
                    onClick={handleSubmit}
                  ><FormattedMessage id='pages.commit' /></Button>
                </ProForm>
              </Col>
            </Row>
          </>}
          {checkSuccess && <>
            <Paragraph style={{ textAlign: 'center' }}>
              <Title level={4}><FormattedMessage id='model.user.forget.password.email.send' /></Title>
              <Text><FormattedMessage id='model.user.forget.password.email.send.description' /></Text>
              <br />
              <br />
            </Paragraph>
            <Row>
              <Col offset={6} lg={12} md={12} sm={24}>
                <Button type="primary" htmlType="submit" size="large" style={{ width: '100%' }}
                  onClick={() => window.location.href = '/'}
                ><FormattedMessage id='pages.return' /></Button>
              </Col>
            </Row>
          </>}
        </Card>
      </div>
    </PageContainer>
  );
};

export default ForgetPassword;

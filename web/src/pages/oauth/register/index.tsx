import { PageContainer, ProFormText, LoginForm } from '@ant-design/pro-components';
import { Footer, ProviderInLine } from '@/components';
import { Alert } from 'antd';
import React from 'react';
import { useParams, useIntl, FormattedMessage,useModel } from '@umijs/max';
import { createStyles } from 'antd-style';
import {
  LockOutlined,
  UserOutlined,
  MobileOutlined,
  MailOutlined,
} from '@ant-design/icons';
import { flushSync } from 'react-dom';
import { registerByOIDC } from '@/services/oauth.api';
import { AccessTokenResponse, RegisterByOIDC } from '@/services/common.d';
import {  addToken } from '@/utils/global';

const useStyles = createStyles(({ token }) => {
  return {
    container: {
      display: 'flex',
      flexDirection: 'column',
      height: '100vh',
      overflow: 'auto',
      backgroundColor: token.colorBgContainer,
    },

  };
});
const OAuthRegister: React.FC = () => {
  const params = useParams();
  const intl = useIntl();
  const provider = params.provider || '';
  
  let initValues = {} as RegisterByOIDC;
  const urlParams = new URL(window.location.href).searchParams;
  initValues.username = urlParams.get('username') || '';
  initValues.email = urlParams.get('email') || '';
  const { styles } = useStyles();
  const defaultPath = '/index';
  const { initialState, setInitialState } = useModel('@@initialState');
  const getAuthedUserInfomation = async () => {
    if (initialState?.getCurrentUserInfo) {
      const userInfo = await initialState.getCurrentUserInfo();
      if (userInfo) {
        flushSync(() => {
          setInitialState((prevState) => ({
            ...prevState,  // 保留其他属性的值
            currentUser: userInfo, // 仅更新 councurrentUser 属性
          }));

        });

      }
    }
  };
  const handleSubmit = async (values: Record<string, any>) => {
    let data = values as RegisterByOIDC;
    data.nickname = data.username;
    data.code = urlParams.get('code') || '';
    
    const token = await registerByOIDC(data) as AccessTokenResponse;
    addToken( token);
    let redirect = urlParams.get('redirect') || defaultPath;
    if (redirect.includes('login') || redirect === '' || redirect === '/') {
      redirect = defaultPath;
    }
    await getAuthedUserInfomation();
    window.location.href = redirect;
  };

  return (
    <PageContainer  title={false}>
      <div className={styles.container}>
        <LoginForm
          contentStyle={{ minWidth: 380, maxWidth: '75vw'}}
          title={intl.formatMessage({ id: 'pages.login.register.title'})}
          subTitle={<div style={{ textAlign: 'center', alignItems: 'center', display: 'flex', justifyContent: 'center' }}><FormattedMessage id="pages.provider.oidc" defaultMessage="OIDC认证提供商" />:&nbsp;<ProviderInLine provider={provider} height={20} /></div>}
          initialValues={initValues}
          onFinish={async (values) => {
            await handleSubmit(values);
          }}
        >
          <>
            <Alert type="warning" showIcon message={intl.formatMessage({ id: 'pages.register.alert.description'})} />
            <br />
            <ProFormText
              name="username"
              fieldProps={{ size: 'large', prefix: <UserOutlined />}}
              placeholder={intl.formatMessage({ id: 'pages.login.username.placeholder'})}
              rules={[
                {
                  required: true,
                  message: (<FormattedMessage id="pages.login.username.required"  />),
                },
              ]}
            />
            <ProFormText.Password
              name="password"
              fieldProps={{ size: 'large', prefix: <LockOutlined />}}
              placeholder={intl.formatMessage({ id: 'pages.login.password.placeholder'})}
            />
            <ProFormText
              name="phone"
              fieldProps={{ size: 'large', prefix: <MobileOutlined />}}
              placeholder={intl.formatMessage({ id: 'model.user.phone'})}
              rules={[
                {
                  required: true,
                  message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.user.phone' }),
                },
                {
                  pattern: /^1\d{10}$/,
                  message: (
                    <FormattedMessage
                      id="pages.login.phoneNumber.invalid"
                      
                    />
                  ),
                },
              ]}
            />
            <ProFormText
              name="email"
              fieldProps={{ size: 'large', prefix: <MailOutlined />}}
              placeholder={intl.formatMessage({ id: 'model.user.email'})}
              rules={[
                {
                  required: true,
                  type: 'email',
                  message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.user.email' }),
                },
              ]}
            />
          </>
        </LoginForm>
      </div>
      <Footer />
    </PageContainer>
  );
};

export default OAuthRegister;

import { Footer, ProviderActionIcons } from '@/components';
import { ProDescriptions } from '@ant-design/pro-components';
import { systemLoginByUsername, thirdAuthMethods, mfaValidate } from '@/services/oauth.api';
import { LoginByFaceIdData, AccessTokenResponse, LoginByUsername, ThirdAuthMethod, MfaCode } from '@/services/common.d';
import {
  LockOutlined,
  UserOutlined,
  MobileOutlined,
  HomeOutlined,
} from '@ant-design/icons';
import {
  LoginFormPage,
  ProFormSelect,
  ProFormText,
  ProFormCaptcha,
} from '@ant-design/pro-components';
import { FormattedMessage, Helmet, SelectLang, useIntl, useModel } from '@umijs/max';
import { message, Tabs, Modal, Button, Image, Input } from 'antd';
import { createStyles } from 'antd-style';
import React, { useState, useEffect, Suspense, lazy } from 'react';
import { flushSync } from 'react-dom';
import Settings from '../../../config/defaultSettings';
import { addToken, getColorPrimary, getSystemOauthUrl, setSystemOauthUrl } from '@/utils/global';
import { checkUserHasFaceIdData, loginByFaceIdData } from '@/services/face_recognition.api';
const FaceRecognitionModal = lazy(() => import('@/components/FaceRecognition'));

const useStyles = createStyles(({ token }) => {
  return {
    action: {
      marginLeft: '8px',
      color: 'rgba(0, 0, 0, 0.2)',
      fontSize: '24px',
      verticalAlign: 'middle',
      cursor: 'pointer',
      transition: 'color 0.3s',
      '&:hover': {
        color: token.colorPrimaryActive,
      },
    },
    lang: {
      width: 42,
      height: 42,
      lineHeight: '42px',
      position: 'fixed',
      right: 16,
      borderRadius: token.borderRadius,
      ':hover': {
        backgroundColor: token.colorBgTextHover,
      },
    },
    home: {
      width: 42,
      height: 42,
      lineHeight: '42px',
      position: 'fixed',
      right: 68,
      textAlign: 'center',
      borderRadius: token.borderRadius,
      ':hover': {
        backgroundColor: token.colorBgTextHover,
      },
    },
    container: {
      display: 'flex',
      flexDirection: 'column',
      height: '100vh',
      overflow: 'auto',
      backgroundImage:
        "url('/login-bg.png')",
      backgroundSize: '100% 100%',
    },

  };
});
const Home = () => {
  const { styles } = useStyles();

  return (
    <div className={styles.home} onClick={() => window.location.href = '/'} >
      <HomeOutlined />
    </div>
  );
}
const Lang = () => {
  const { styles } = useStyles();

  return (
    <div className={styles.lang} data-lang>
      {SelectLang && <SelectLang />}
    </div>
  );
};

const Login: React.FC = () => {
  const colorPrimary = getColorPrimary();
  const [authMethod, setAuthMethod] = useState<ThirdAuthMethod>({} as ThirdAuthMethod);
  const [type, setType] = useState<string>('account');
  const [username, setUsername] = useState<string>('');
  const { initialState, setInitialState } = useModel('@@initialState');
  const { styles } = useStyles();
  const intl = useIntl();
  const [mfaVisible, setMfaVisible] = useState(false);
  const [token, setToken] = useState<AccessTokenResponse>();
  const [mfaCodeValue, setMfaCodeValue] = useState('');
  const [faceRecognitionModalVisibe, setFaceRecognitionModalVisibe] = useState(false);

  const defaultRedirect = (): string => {
    let oauthUrl = getSystemOauthUrl();
    if (oauthUrl !== '' && oauthUrl.indexOf('/user/login', -1) < 0) {
      if (oauthUrl === `${window.location.origin}/index`) {
        return '/index';
      }
      return oauthUrl;
    }
    let re = window.location.href.split("?redirect=")
    if (re.length === 2) {
      return re[1];
    }
    return '/index';
  };
  const getLoginItems = () => {
    let loginItems = [
      {
        key: 'account',
        label: intl.formatMessage({
          id: 'pages.login.accountLogin.tab',
          defaultMessage: '账户密码',
        }),
      },
      {
        key: 'mobile',
        label: intl.formatMessage({
          id: 'pages.login.phoneLogin.tab',
          defaultMessage: '手机验证码',
        }),
      },
    ];
    if (authMethod?.faceRecognition) {
      loginItems.push({
        key: 'faceRecognition',
        label: intl.formatMessage({
          id: 'pages.login.faceRecognition.tab',
          defaultMessage: '人脸识别',
        }),
      });
    };
    return loginItems;

  };

  const getAuthMethods = async () => {
    setAuthMethod(await thirdAuthMethods());
  }
  useEffect(() => {
    if (initialState?.currentUser) {
      const currentUser = initialState.currentUser;
      if (currentUser && currentUser?.id > 0) {
        window.location.href = defaultRedirect();
      }
    }
    getAuthMethods();
  }, []);

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
  const defaultLoginSuccessMessage = intl.formatMessage({
    id: 'pages.login.success',
    defaultMessage: '登录成功！',
  });
  const handleSubmit = async (values: Record<string, any>) => {
    setMfaCodeValue('');
    let token = {} as AccessTokenResponse;
    try {
      if (type === 'account') {
        token = await systemLoginByUsername(values as LoginByUsername);
        console.log(token);
      }   else if (type === 'faceRecognition') {
        token = await loginByFaceIdData(values as LoginByFaceIdData);
      }
      if (token.mfa === true) {
        setMfaVisible(true);
        setToken(token);
      } else {
        addToken(  token);
        message.success(defaultLoginSuccessMessage);
        await getAuthedUserInfomation();
        const redirect = defaultRedirect();
        setSystemOauthUrl("/index");
        window.location.href = redirect
      }

    } catch (error) {
      console.log(error);
      const defaultLoginFailureMessage = intl.formatMessage({
        id: 'pages.login.failure',
      });
      message.error(defaultLoginFailureMessage);
    }
  };
  const mfaValidator = async () => {
    const values = { userId: token?.id, code: mfaCodeValue } as MfaCode
    const t = await mfaValidate(values) as AccessTokenResponse;
    setMfaVisible(false);
    addToken( t);
    message.success(defaultLoginSuccessMessage);
    window.location.href = defaultRedirect();
  };

  return (
    <div className={styles.container}>
      <Helmet>
        <title>
          {intl.formatMessage({
            id: 'menu.login',
            defaultMessage: '登录页',
          })}
          - {Settings.title}
        </title>
      </Helmet>
      <Home />
      <Lang />
      <div
        style={{
          flex: '1',
          padding: '32px 0',
        }}
      >
        <LoginFormPage
          contentstyle={{
            minWidth: 380,
            maxWidth: '75vw',
          }}
          mainStyle={{ minWidth: 380 }}
          logo={<img alt="logo" src="/logo.png" />}
          title={<span onClick={() => { window.location.href = '/'; }}><FormattedMessage id='pages.layouts.userLayout.title' /></span>}
          subTitle={intl.formatMessage({ id: 'pages.layouts.userLayout.subTitle' })}
          initialValues={{
            rememberMe: '-',
          }}
          actions={[
            <ProviderActionIcons key="icons" oidcs={authMethod.oidcs ? authMethod.oidcs : []} allowRegistry={true} />,
          ]}
          onFinish={async (values) => {
            if (type === 'faceRecognition') {
              setUsername(values.username);
              const exist = await checkUserHasFaceIdData({ username: values.username });
              if (exist === 'exist') {
                setFaceRecognitionModalVisibe(true);
              } else {
                message.error(intl.formatMessage({ id: 'pages.login.faceRecognition.notExist' }));
                return
              }
            } else {
              await handleSubmit(values);
            }

          }}
        >
          <Tabs activeKey={type} onChange={setType} centered items={getLoginItems()} />
          {type === 'account' && (
            <>
              <ProFormText
                name="username"
                fieldProps={{
                  size: 'large',
                  prefix: <UserOutlined />,
                }}
                placeholder={intl.formatMessage({
                  id: 'pages.login.username.email.phone',
                })}
                rules={[
                  {
                    required: true,
                    message: (<FormattedMessage id="pages.login.username.email.phone" />),
                  },
                ]}
              />
              <ProFormText.Password
                name="password"
                fieldProps={{
                  size: 'large',
                  prefix: <LockOutlined />,
                }}
                placeholder={intl.formatMessage({
                  id: 'pages.login.password.placeholder',
                })}
                rules={[
                  {
                    required: true,
                    message: (
                      <FormattedMessage
                        id="pages.login.password.required"

                      />
                    ),
                  },
                ]}
              />
            </>
          )}
          {type === 'mobile' && (
            <>
              <ProFormText
                fieldProps={{
                  size: 'large',
                  prefix: <MobileOutlined />,
                }}
                name="mobile"
                placeholder={intl.formatMessage({
                  id: 'pages.login.phoneNumber.placeholder',
                })}
                rules={[
                  {
                    required: true,
                    message: (
                      <FormattedMessage
                        id="pages.login.phoneNumber.required"

                      />
                    ),
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
              <ProFormCaptcha
                fieldProps={{
                  size: 'large',
                  prefix: <LockOutlined />,
                }}
                captchaProps={{
                  size: 'large',
                }}
                placeholder={intl.formatMessage({
                  id: 'pages.login.captcha.placeholder',
                })}
                captchaTextRender={(timing, count) => {
                  if (timing) {
                    return `${count} ${intl.formatMessage({
                      id: 'pages.getCaptchaSecondText',
                    })}`;
                  }
                  return intl.formatMessage({
                    id: 'pages.login.phoneLogin.getVerificationCode',
                  });
                }}
                name="captcha"
                rules={[
                  {
                    required: true,
                    message: (
                      <FormattedMessage
                        id="pages.login.captcha.required"
                      />
                    ),
                  },
                ]}
                onGetCaptcha={async (phone) => {
                  console.log(phone);
                  message.success('获取验证码成功！验证码为：1234');
                }}
              />
            </>
          )}
          {type === 'webauthn' && (
            <>
              <ProFormText
                name="username"
                fieldProps={{
                  size: 'large',
                  prefix: <UserOutlined />,
                }}
                placeholder={intl.formatMessage({
                  id: 'pages.login.username.placeholder',
                })}
                rules={[
                  {
                    required: true,
                    message: (
                      <FormattedMessage
                        id="pages.login.username.required"
                      />
                    ),
                  },
                ]}
              />
            </>
          )}
          {type === 'faceRecognition' && (
            <>
              <ProFormText
                name="username"
                fieldProps={{
                  size: 'large',
                  prefix: <UserOutlined />,
                }}
                placeholder={intl.formatMessage({
                  id: 'pages.login.username.placeholder',

                })}
                rules={[
                  {
                    required: true,
                    message: (
                      <FormattedMessage
                        id="pages.login.username.required"
                      />
                    ),
                  },
                ]}
              />
            </>
          )}
          <div style={{ marginBottom: 24 }} >
            {type === 'account' && <a style={{ color: colorPrimary }}  onClick={() => window.location.href = "/forget/password"} >
              <FormattedMessage id="pages.login.forgotPassword" />
            </a>}
          </div>
        </LoginFormPage>
      </div>
      {type === 'faceRecognition' && <Suspense>
        <FaceRecognitionModal visible={faceRecognitionModalVisibe} onSuccess={async (faceIdData: number[]) => {
          setFaceRecognitionModalVisibe(false);
          let values = {} as LoginByFaceIdData;
          values.faceIdData = faceIdData;
          values.username = username;
          await handleSubmit(values);
        }} onCancel={() => { setFaceRecognitionModalVisibe(false) }} />


      </Suspense>}
      <Modal
        open={mfaVisible}
        title={intl.formatMessage({ id: 'pages.validate.code.title' })}
        closable
        onCancel={() => { setMfaVisible(false) }}
        footer={[
          <Button key="confirm" type='primary' onClick={mfaValidator}>
            <FormattedMessage id='pages.operation.confirm' />
          </Button>
        ]}
        style={{ textAlign: 'center' }}
      >
        <br />
        {token?.image && <Image
          src={token.image}
          preview={false}
          width={200}
        />}
        <br />
        <br />
        {token?.secret && <>
          <ProDescriptions column={1}>
            <ProDescriptions.Item copyable >{token.secret}</ProDescriptions.Item>
          </ProDescriptions>
        </>}
        <Input placeholder={intl.formatMessage({ id: 'pages.validate.code' })} allowClear onChange={(e) => {
          setMfaCodeValue(e.target.value);
        }} />
      </Modal>
      <Footer />
    </div>
  );
};

export default Login;

import { PageContainer, ProDescriptions } from '@ant-design/pro-components';
import { Footer, Provider } from '@/components';
import { Modal, Button, Image,  Input } from "antd";
import React, { useEffect, useState } from 'react';
import { FormattedMessage, useParams,useIntl } from '@umijs/max';
import { createStyles } from 'antd-style';
import { systemLoginByOidc,mfaValidate } from '@/services/oauth.api';
import { LoginByOIDC, AccessTokenResponse ,MfaCode} from '@/services/common.d';
import { getBind, deleteBind, getSearchParams,  addToken, getSystemOauthUrl } from '@/utils/global';

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

const OAuthCallback: React.FC = () => {
  const intl = useIntl();
  const params = useParams();
  const provider = params.provider || '';
  const defaultPath = getSystemOauthUrl() || '/';
  
  const [mfaVisible, setMfaVisible] = useState(false);
  const [token, setToken] = useState<AccessTokenResponse>();
  const [value, setValue] = useState('');
  const [defaultRedirect,setDefaultRedirect] =useState('/');
  const { styles } = useStyles();
  const login = async () => {
    let loginParam = {} as LoginByOIDC;
    loginParam.code = getSearchParams().get('code') as string;
    loginParam.provider = provider;
    loginParam.redirectUri = window.location.origin + window.location.pathname;
    const bindId = getBind();
    deleteBind();
    if (bindId) loginParam.bindId = bindId;
    const urlParams = new URL(window.location.href).searchParams;

    let redirect = urlParams.get('redirect') || defaultPath;
    if (redirect.includes('login') || redirect === '' || redirect === '/') {
      redirect = defaultPath;
    }
    setDefaultRedirect(redirect);
    const token = await systemLoginByOidc(loginParam) as AccessTokenResponse;
    if (token.need === true) {
      window.location.href = `/oauth/register/${provider}?code=${token.code}&username=${token.username||''}&nickname=${token.nickname||''}&email=${token.email||''}&phone=${token.phone||''}&redirect=${redirect}`;
      return;
    } else if (token.mfa === true) {
      setMfaVisible(true);
      setToken(token);
    } else {
      addToken(  token);
      window.location.href = redirect;
    }
  };
  const mfaValidator= async ()=>{
    const values= {userId:token?.id,code:value}as MfaCode
    const t =await  mfaValidate(values) as AccessTokenResponse;
    setMfaVisible(false);
    addToken( t);
    window.location.href = defaultRedirect;
  };
  
  useEffect(() => { login(); }, []);
  return (
    <PageContainer title={false}>
      <div className={styles.container}>
       <div onClick={()=>{ window.location.href = '/user/login';}}> <Provider provider={provider} height={100} /></div>
      </div>
      <Modal
        open={mfaVisible}
        title={intl.formatMessage({id:'pages.validate.code.title',defaultMessage:'请使用身份验证器扫码或者输入下面的密钥'})} 
        closable
        onCancel={()=>{setMfaVisible(false)}}
        footer={[
          <Button key="confirm" type='primary'  onClick={mfaValidator}>
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
        <Input placeholder= {intl.formatMessage({id:'pages.validate.code',defaultMessage:'请输入验证码'})} allowClear onChange={(e) => {
          setValue(e.target.value);
        }} />
      </Modal>
      <Footer />
    </PageContainer>
  );
};

export default OAuthCallback;

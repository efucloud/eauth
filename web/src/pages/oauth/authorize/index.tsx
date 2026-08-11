import { PageContainer } from '@ant-design/pro-components';
import { Footer } from '@/components';
import React, { useEffect } from 'react';
import { FormattedMessage } from '@umijs/max';
import { getOidcCode } from '@/services/oauth.api';
import { OidcCodeResponse } from '@/services/common.d';
import {  getSearchParams, setSystemOauthUrl, deleteSystemOauthUrl } from '@/utils/global';


const Authorize: React.FC = () => {
  const clientId = getSearchParams().get('client_id') as string;
  const redirectUri = getSearchParams().get('redirect_uri') as string;
  const state = getSearchParams().get('state') as string;
  const responseType = getSearchParams().get('response_type') as string;
  
  setSystemOauthUrl(window.location.href);
  const getCode = async () => {
    const res = await getOidcCode({  clientId, redirectUri, state, responseType }) as OidcCodeResponse;
    if (res.redirectUri) {
      deleteSystemOauthUrl();
      window.location.href = res.redirectUri;
    }
  };
  useEffect(() => {
    getCode();
  });
  return (
    <PageContainer title={false} >
      <div style={{ textAlign: 'center', marginTop: '100px' }}><FormattedMessage id='app.description' /></div>
      <Footer />
    </PageContainer>
  );

};

export default Authorize;

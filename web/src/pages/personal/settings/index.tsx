import {   PageContainer } from '@ant-design/pro-components';
import React, {   useEffect, useRef, useState } from 'react';
import BaseView from './components/base';
import BindingView from './components/binding';
import SecurityView from './components/security';
import { useIntl, useModel } from '@umijs/max';
import Mfa from './components/mfa';
import { ThirdAuthMethod } from '@/services/common.d';
import { thirdAuthMethods } from '@/services/oauth.api';
import FaceRecognition from './components/face_recognition';
 
const Settings: React.FC = () => {
  const intl = useIntl();
  const [tabList, setTabList] = useState<[]>([]);
  const getTablist = async () => {
    const res = await thirdAuthMethods() as ThirdAuthMethod;
    let tabs = [{
      key: 'base',
      tab: intl.formatMessage({ id: 'pages.personal.settings.base' }),
      children: <BaseView />,
    },
    {
      key: 'security',
      tab: intl.formatMessage({ id: 'pages.personal.settings.security' }),
      children:<SecurityView   />,
    }
    ];
    if (res.oidcs?.length > 0) {
      tabs.push({
        key: 'binding',
        tab: intl.formatMessage({ id: 'pages.personal.settings.bind' }),
        children:<BindingView />,
      });
    }
    if (res.mfa) {
      tabs.push({
        key: 'mfa',
        tab: intl.formatMessage({ id: 'pages.personal.settings.mfa' }),
        children:<Mfa />,
      });
    };
    if (res?.faceRecognition) {
      tabs.push(
        {
          key: 'faceRecognition',
          tab: intl.formatMessage({ id: 'pages.personal.settings.faceRecognition' }),
          children:<FaceRecognition />,
        });
    }
    setTabList(tabs);
  };

  useEffect(() => {
    getTablist();
  }, []);

  return (
    <PageContainer
      breadcrumb={false}
      ghost
      title={false}
      subTitle={false}
      tabList={tabList}
      
    >
    </PageContainer>
  );
};
export default Settings;

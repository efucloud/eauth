import { PageContainer, ProDescriptions } from '@ant-design/pro-components';
import { useParams, useIntl, useNavigate, FormattedMessage, Access, useAccess } from '@umijs/max';
import { Card, Space, Button, Tabs } from 'antd';
import { EditOutlined } from '@ant-design/icons';
import React, { useEffect, useState } from 'react';
import type { ApplicationDetail } from '@/services/application.d';
import { getApplication } from '@/services/application.api';
import { getColorPrimary } from '@/utils/global';
const BaseAddress = "/app";

const AppDetail: React.FC = () => {
  const colorPrimary = getColorPrimary();
  const access = useAccess();
  const navigate = useNavigate();
  const intl = useIntl();
  const params = useParams();
  const [appInfo, setAppInfo] = useState<ApplicationDetail>();
  const getInfo = async () => {
    const data = await getApplication({ id: params.id! });
    setAppInfo(data as ApplicationDetail);
  }
  useEffect(() => { getInfo(); }, [params.id]);
  const selectItem = {
    false: {
      text: intl.formatMessage({ id: 'model.enable.bool.false' }),
      status: 'Error',
    },
    true: {
      text: intl.formatMessage({ id: 'model.enable.bool.true' }),
      status: 'Success',
    },
  };
  return (
    <PageContainer title={intl.formatMessage({ id: 'menu.application' })} header={{ breadcrumb: {}, onBack: () => navigate(`${BaseAddress}`) }} >
      <Card bordered={false}
      >
        <ProDescriptions style={{ marginBottom: 32 }} column={3} title={intl.formatMessage({ id: 'pages.detail.baseinfo' })}
          extra={
            <Space>
              <Access accessible={access.systemAdminAccess === true}>
                <Button type={'primary'} icon={<EditOutlined />} block
                  onClick={() => {
                    navigate(
                      {
                        pathname: `${BaseAddress}/update/${params.id}`,
                      },
                      { replace: true },
                    );
                  }}
                ><FormattedMessage id="pages.operation.edit" /></Button>
              </Access>
            </Space>
          }>
          <ProDescriptions.Item label={intl.formatMessage({ id: 'model.application.name' })}>{appInfo?.name}</ProDescriptions.Item>
          <ProDescriptions.Item label={intl.formatMessage({ id: 'model.application.code' })}>{appInfo?.code}</ProDescriptions.Item>
          <ProDescriptions.Item tooltip={{ color: colorPrimary, title: intl.formatMessage({ id: 'model.application.enable.description' }) }} label={intl.formatMessage({ id: 'model.application.enable' })}
            valueEnum={{
              false: {
                text: intl.formatMessage({ id: 'model.application.enable.disable' }),
                status: 'Error',
              },
              true: {
                text: intl.formatMessage({ id: 'model.application.enable.enable' }),
                status: 'Success',
              },
            }}
          >{appInfo?.enable}</ProDescriptions.Item>
          <ProDescriptions.Item label={intl.formatMessage({ id: 'model.application.home' })}><a style={{ color: colorPrimary }} href={appInfo?.home} target='_blank' rel="noreferrer">{appInfo?.home}</a></ProDescriptions.Item>
          <ProDescriptions.Item label={intl.formatMessage({ id: 'model.application.redirectUri' })}><a style={{ color: colorPrimary }} href={appInfo?.home} target='_blank' rel="noreferrer">{appInfo?.redirectUri}</a></ProDescriptions.Item>
          <ProDescriptions.Item label={intl.formatMessage({ id: 'model.application.redirectUriMatchType' })}
            valueEnum={{
              'equal': {
                text: intl.formatMessage({ id: 'model.application.redirect_uri.macth.equal' }),
              },
              'regex': {
                text: intl.formatMessage({ id: 'model.application.redirect_uri.macth.regex' }),
              },
              'prefix': {
                text: intl.formatMessage({ id: 'model.application.redirect_uri.macth.prefix' }),
              },
              'contain': {
                text: intl.formatMessage({ id: 'model.application.redirect_uri.macth.contain' }),
              },
            }}
          >{appInfo?.redirectUriMatchType}</ProDescriptions.Item>
          <ProDescriptions.Item copyable ellipsis tooltip={{ color: colorPrimary, title: intl.formatMessage({ id: 'model.application.clientId.tooltip' }) }} label={intl.formatMessage({ id: 'model.application.clientId' })}>{appInfo?.clientId}</ProDescriptions.Item>
          <ProDescriptions.Item copyable ellipsis tooltip={{ color: colorPrimary, title: intl.formatMessage({ id: 'model.application.clientSecret.tooltip' }) }} label={intl.formatMessage({ id: 'model.application.clientSecret' })}>{appInfo?.clientSecret}</ProDescriptions.Item>

          <ProDescriptions.Item valueType='textarea' label={intl.formatMessage({ id: 'model.application.description' })}>{appInfo?.description}</ProDescriptions.Item>
        </ProDescriptions>
      </Card>

    </PageContainer>
  );
};
export default AppDetail;


import { PageContainer, ProDescriptions } from '@ant-design/pro-components';
import { useParams, useIntl, useNavigate, FormattedMessage, Access, useAccess } from '@umijs/max';
import { Card, Space, Button } from 'antd';
import { Provider, ProviderInLine } from '@/components';
import { EditOutlined } from '@ant-design/icons';
import React, { useEffect, useState } from 'react';
import type { ProviderOidcDetail } from '@/services/system_provider_oidc.d';
import { getProviderOidc } from '@/services/provider_oidc.api';
const BaseAddress = '/provider/oidc';


const ProviderOidcDetailIndex: React.FC = () => {
  const access = useAccess();
  const navigate = useNavigate();
  const intl = useIntl();
  const params = useParams();
  const [instanceInfo, setInstanceInfo] = useState<ProviderOidcDetail>();
  const getInstanceInfo = async () => {
    setInstanceInfo(await getProviderOidc({ id: params.id! }) as ProviderOidcDetail);
  }
  useEffect(() => { getInstanceInfo(); }, [params.id]);

  return (
    <PageContainer title={intl.formatMessage({ id: 'menu.provider.oidc' })} header={{ breadcrumb: {}, onBack: () => navigate(`${BaseAddress}`) }} >
      <Card bordered={false}
      >
        <ProDescriptions style={{ marginBottom: 32 }} column={2} title={(<Space><ProviderInLine provider={instanceInfo?.category ? instanceInfo.category : ''} height={25} /></Space>)}
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
          <ProDescriptions.Item label={intl.formatMessage({ id: 'model.provider.oidc.name' })}>{instanceInfo?.name}</ProDescriptions.Item>

          <ProDescriptions.Item label={intl.formatMessage({ id: 'model.provider.oidc.enable' })}
            valueEnum={{
              false: {
                text: intl.formatMessage({ id: 'model.provider.oidc.enable.disable' }),
                status: 'Error',
              },
              true: {
                text: intl.formatMessage({ id: 'model.provider.oidc.enable.enable' }),
                status: 'Success',
              },
            }}
          >{instanceInfo?.enable}</ProDescriptions.Item>
          <ProDescriptions.Item copyable ellipsis label={intl.formatMessage({ id: 'model.provider.oidc.clientId' })}>{instanceInfo?.clientId}</ProDescriptions.Item>
          <ProDescriptions.Item copyable ellipsis label={intl.formatMessage({ id: 'model.provider.oidc.clientSecret' })}>{instanceInfo?.clientSecret}</ProDescriptions.Item>
          <ProDescriptions.Item copyable label={intl.formatMessage({ id: 'model.provider.oidc.authorizationEndpoint' })}>{instanceInfo?.authorizationEndpoint}</ProDescriptions.Item>
          <ProDescriptions.Item copyable label={intl.formatMessage({ id: 'model.provider.oidc.tokenEndpoint' })}>{instanceInfo?.tokenEndpoint}</ProDescriptions.Item>
          <ProDescriptions.Item copyable label={intl.formatMessage({ id: 'model.provider.oidc.userinfoEndpoint' })}>{instanceInfo?.userinfoEndpoint}</ProDescriptions.Item>

        </ProDescriptions>
      </Card>
      <br />


    </PageContainer>
  );
};
export default ProviderOidcDetailIndex;



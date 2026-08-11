import { PageContainer } from '@ant-design/pro-components';
import { useNavigate } from '@umijs/max';
import React, { useState, useEffect } from 'react';
import { FormattedMessage, useIntl, useAccess } from '@umijs/max';
import { listProviderOidc, changeProviderOidcStatus } from '@/services/provider_oidc.api';
import { ProviderOidcDetail, ProviderOidcDetailList } from '@/services/system_provider_oidc.d';
import { Provider, ProviderTwoLine } from '@/components';
import { Card, List, Popconfirm, message, Tooltip } from 'antd';
const BaseAddress = '/provider/oidc';
import { IntlShape } from "react-intl";
import { getColorPrimary } from '@/utils/global';

const ProvicerOidcCardList: React.FC = () => {
  const intl = useIntl();
  const access = useAccess();
  const navigate = useNavigate();
  const [providers, setProviders] = useState<ProviderOidcDetail[]>([]);
  const allProviders = ['gitlab', 'github', 'wechat', 'alipay', 'feiShu', 'weibo', 'wechatWork', 'qq', 'dingTalk', 'gitee', 'baidu', 'tiktok', 'bilibili', 'custom'];
  const colorPrimary = getColorPrimary();

  const hasProviderConfig = (item: ProviderOidcDetail) => {
    const id = String(item.id ?? '').trim();
    return id !== '' && id !== '0';
  };

  const getProviders = async () => {
    const data = await listProviderOidc({}) as ProviderOidcDetailList;
    let providers = [] as ProviderOidcDetail[];
    for (let i = 0; i < allProviders.length; i++) {
      providers.push({
        id: 0,
        createdAt:'',
        updatedAt:'',
        name: intl.formatMessage({ id: `model.provider.oidc.category.${allProviders[i]}`, defaultMessage: allProviders[i] }),
        category: allProviders[i],
        issuer: '',
        enable: false,
        authorizationEndpoint: '',
        tokenEndpoint: '',
        userinfoEndpoint: '',
        clientId: '',
        clientSecret: '',
        scopes: [],
      })
    }
    if (data && data?.total > 0) {
      data.data?.forEach((item) => {
        providers.find((items) => {
          if (items.category === item.category) {
            items.id = item.id;
            items.name = item.name;
            items.category = item.category;
            items.issuer = item.issuer;
            items.enable = item.enable;
            items.authorizationEndpoint = item.authorizationEndpoint;
            items.tokenEndpoint = item.tokenEndpoint;
            items.userinfoEndpoint = item.userinfoEndpoint;
            items.clientId = item.clientId;
            items.clientSecret = item.clientSecret;
            items.scopes = item.scopes;
            return true;
          }
          return false;
        });
      });
    }
    setProviders(providers)
  };
  useEffect(() => { getProviders(); }, []);
  const content = (
    <div  >
      <p>
        <FormattedMessage id='model.provider.oidc.third.support'  />
      </p>

    </div>
  );
  const handleEnable = async (intl: IntlShape, enable: boolean, selectedRows: ProviderOidcDetail[]) => {
    const hide = message.loading(intl.formatMessage({ id: 'pages.operation.updating' }));
    if (!selectedRows) return true;

    try {
      await changeProviderOidcStatus({
        ids: selectedRows.map((row) => row.id),
        enable: enable,
      });
      hide();
      message.success(intl.formatMessage({ id: 'pages.operation.update.success' }));
      getProviders();
      return true;
    } catch (error) {
      hide();
      message.error(intl.formatMessage({ id: 'pages.operation.update.failed' }));
      return false;
    }
  };
  const cardExtra = (item: ProviderOidcDetail): React.ReactNode[] => {
    const nodes = [] as React.ReactNode[];
    const configured = hasProviderConfig(item);
    if (configured) {
      if (item.enable === true) {
        if (access.systemEditAccess === true) {
          nodes.push(
            <Popconfirm
              key={item.category + '-disable'}
              title={intl.formatMessage({ id: 'pages.operation.disable.description'}) +
                intl.formatMessage({ id: 'provider.oidc.model.name'}) + '【' + item.name + '】'
              }
              description={intl.formatMessage({ id: 'model.provider.oidc.enable.disable.description'})}
              onConfirm={() => {
                handleEnable(intl, false, [item]);
              }}
              okText={intl.formatMessage({ id: 'pages.operation.confirm'})}
              cancelText={intl.formatMessage({ id: 'pages.operation.cancel'})}
            >
              <a style={{color:colorPrimary}} className='operation-disable'><FormattedMessage id="pages.operation.disable" defaultMessage="禁用" />&nbsp;&nbsp;</a>
            </Popconfirm>
          );
        }

      } else {
        if (access.systemEditAccess === true) {
          nodes.push(
            <Popconfirm
              key={item.category + '-enable'}
              title={intl.formatMessage({ id: 'pages.operation.enable.description'}) +
                intl.formatMessage({ id: 'provider.oidc.model.name'}) + '【' + item.name + '】'
              }
              description={intl.formatMessage({ id: 'model.provider.oidc.enable.enable.description'})}
              onConfirm={() => {
                handleEnable(intl, true, [item]);
              }}
              okText={intl.formatMessage({ id: 'pages.operation.confirm'})}
              cancelText={intl.formatMessage({ id: 'pages.operation.cancel'})}
            >
              <a style={{color:colorPrimary}} className='operation-enable'><FormattedMessage id="pages.operation.enable" defaultMessage="启用" />&nbsp;&nbsp;</a>
            </Popconfirm>
          );
        }
      }
    }
    if (configured) {
      if (access.systemEditAccess === true) {
        nodes.push(
          <a
          style={{color:colorPrimary}}
            key="edit"
            onClick={() => {
              navigate(
                {
                  pathname: `${BaseAddress}/update/${item.id}`,
                },
                { replace: true },
              );
            }}
              >
            <FormattedMessage id="pages.operation.edit" defaultMessage="修改" />&nbsp;&nbsp;
          </a>

        );
      }
      if (access.systemViewAccess === true) {
        nodes.push(
          <a
          style={{color:colorPrimary}}
            key={item.category + '-detail'}
            onClick={() => {
              navigate(
                {
                  pathname: `${BaseAddress}/detail/${item.id}`,
                },
                { replace: true },
              );
            }}
          >
            <FormattedMessage id="pages.operation.detail" defaultMessage="详情" />&nbsp;&nbsp;
          </a>
        );
      }
    } else {
      if (access.systemAdminAccess === true) {
        nodes.push(
          <Tooltip key={item.category + '-edit'} placement="right" color={colorPrimary} title={intl.formatMessage({ id: 'model.provider.oidc.unconfig' })}>
            <a
             style={{color:colorPrimary}}
              key={item.category + '-config'}
              onClick={() => {
                navigate(
                  {
                    pathname: `${BaseAddress}/create?category=${item.category}`,
                  },
                  { replace: true },
                );
              }}
            >
              <FormattedMessage id="pages.operation.wait.config" defaultMessage="待配置" />&nbsp;&nbsp;
            </a>
          </Tooltip>
        );
      }
    }
    return nodes;
  };
  return (
    <PageContainer title={intl.formatMessage({ id: 'menu.provider.oidc' })} header={{ breadcrumb: {} }} content={content}  >
      <div  >
        <List<ProviderOidcDetail>

          rowKey="id"
          grid={{
            gutter: 16,
            xs: 1,
            sm: 2,
            md: 4,
            lg: 4,
            xl: 4,
            xxl: 4,
          }}
          dataSource={providers}
          renderItem={(item) => {
            return (
              <List.Item key={item.category} >
                <Card
                  hoverable
                  extra={cardExtra(item as ProviderOidcDetail)}
                  key={item.category}
                >
                  <Card.Meta
                    title={<><p>{item.name} </p>
                    </>}
                    avatar={<ProviderTwoLine provider={item?.category}  />}
                  />

                </Card>
              </List.Item>
            );
          }}
        />
      </div>
    </PageContainer>
  );
};
export default ProvicerOidcCardList;

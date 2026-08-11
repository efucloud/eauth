import {
  FooterToolbar,
  PageContainer,
  ProForm,
  ProFormText,
  ProFormList,
  ProFormSwitch,
} from '@ant-design/pro-components';
import { Card, Col, Row, Space } from 'antd';
import type { FC } from 'react';

import { useEffect, useState } from 'react';
import { useParams, useIntl, useNavigate, FormattedMessage } from '@umijs/max';
import { getProviderOidc, updateProviderOidc, createProviderOidc } from '@/services/provider_oidc.api';
import { ProviderOidcDetail, ProviderOidcCreate, ProviderOidcUpdate } from '@/services/system_provider_oidc.d';
import { getColorPrimary, getSearchParams } from '@/utils/global';
import { ProviderInLine } from '@/components';
const Update = 'update';
const Create = 'create';
const BaseAddress = '/provider/oidc';

const ProviderOidcForm: FC<Record<string, any>> = () => {
  const [showMoreFields, setShowMoreFields] = useState<boolean>(false);
  const [category, setCategory] = useState<string>('');
  const navigate = useNavigate();
  const intl = useIntl();
  const params = useParams();
  let initialValues = {} as ProviderOidcDetail;
  const colorPrimary = getColorPrimary();
  useEffect(() => {
    const searchParams = getSearchParams();
    let categ = searchParams.get('category');
    if (categ) {
      setCategory(categ);
      initialValues.category = categ;
      initialValues.autoGetAuthAddress = true;
    }
  }, []);
  const mode = params.id ? Update : Create;

  const onFinish = async (values: Record<string, any>) => {
    values.category = category;
    if (values.scopes) {
      const { scopes } = values;
      values.scopes = scopes.map((item: any) => item.value);
    }
    if (mode === Create) {
      const detail = await createProviderOidc(values as ProviderOidcCreate) as ProviderOidcDetail;
      if (detail.id) {
        // 跳转到详情页
        navigate(`${BaseAddress}/detail/${detail.id}`);
      }
    } else {
      values['id'] = params.id;
      const detail = await updateProviderOidc(values as ProviderOidcUpdate) as ProviderOidcDetail;
      if (detail.id) {
        // 跳转到详情页
        navigate(`${BaseAddress}/detail/${detail.id}`);
      }
    }

  };
  //          

  const onInitData = async () => {
    let initData = {} as ProviderOidcDetail;
    if (mode === Update) {
      const res = await getProviderOidc({ id: params.id });
      initData = res as ProviderOidcDetail;
      setCategory(initData.category);
      const { scopes } = initData;
      if (scopes) {
        initData.scopes = scopes.map((item: string) => ({ value: item }));
      }
      initData.autoGetAuthAddress = true;
      return initData;
    } else {
      initData.enable = true;
      initData.autoGetAuthAddress = true;
      return initData;
    }
  };

  return (
    <ProForm
      layout="horizontal"
      submitter={{
        render: (props, dom) => {
          return (
            <FooterToolbar>
              {dom}
            </FooterToolbar>
          );
        },
      }}
      onValuesChange={(changeValues) => {
        if (changeValues.autoGetAuthAddress === false) {
          setShowMoreFields(true);
        } else if (changeValues.autoGetAuthAddress) {
          setShowMoreFields(false);
        }
      }}
      initialValues={initialValues}
      request={onInitData}
      onFinish={onFinish}
    >
      <PageContainer
        title={intl.formatMessage({ id: 'menu.provider.oidc' })}
        subTitle={<><FormattedMessage id='pages.oauth2.callback' />: <a style={{ color: colorPrimary }} href=''>{window.location.origin}/oauth/callback/{category}</a></>}
        header={{ breadcrumb: {}, onBack: () => navigate(`${BaseAddress}`) }}>
        <Card title={(<Space><ProviderInLine provider={category} /></Space>)}
          bordered={false}>
          <Row gutter={{ xs: 16, sm: 32, md: 64, lg: 128 }}>
            <Col lg={12} md={12} sm={24} >
              <ProFormText
                label={intl.formatMessage({ id: 'model.provider.oidc.name' })}
                name="name"
                rules={[
                  {
                    required: true,
                    message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.provider.oidc.name' }),
                  },
                ]}
              />
            </Col>

            <Col lg={12} md={12} sm={24}>
              <ProFormSwitch
                label={intl.formatMessage({ id: 'model.provider.oidc.enable' })}
                name="enable"
              />
            </Col>
          </Row>
          <Row gutter={{ xs: 16, sm: 32, md: 64, lg: 128 }}>
            <Col lg={12} md={12} sm={24} >
              <ProFormText
                label={intl.formatMessage({ id: 'model.provider.oidc.clientId' })}
                name="clientId"
                rules={[
                  {
                    required: true,
                    message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.provider.oidc.clientId' }),
                  },
                ]}
              />
            </Col>
            <Col lg={12} md={12} sm={24} >
              <ProFormText.Password
                label={intl.formatMessage({ id: 'model.provider.oidc.clientSecret' })}
                name="clientSecret"
                rules={[
                  {
                    required: true,
                    message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.provider.oidc.clientSecret' }),
                  },
                ]}
              />
            </Col>
            <Col lg={12} md={12} sm={24}>
              <ProFormSwitch
                label={intl.formatMessage({ id: 'model.provider.oidc.autoGetAuthAddress' })}
                name="autoGetAuthAddress"
              />
            </Col>
          </Row>

          {showMoreFields && <Row gutter={{ xs: 16, sm: 32, md: 64, lg: 128 }}>
            <Col lg={12} md={12} sm={24} >
              <ProFormText
                label={intl.formatMessage({ id: 'model.provider.oidc.authorizationEndpoint' })}
                name="authorizationEndpoint"
                rules={[
                  {
                    required: true,
                    message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.provider.oidc.authorizationEndpoint' }),
                  },
                  { type: 'url', message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.url.right' }) },

                ]}
              />
            </Col>
            <Col lg={12} md={12} sm={24} >
              <ProFormText
                label={intl.formatMessage({ id: 'model.provider.oidc.tokenEndpoint' })}
                name="tokenEndpoint"
                rules={[
                  {
                    required: true,
                    message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.provider.oidc.tokenEndpoint' }),
                  },
                  { type: 'url', message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.url.right' }) },

                ]}
              />
            </Col>
          </Row>}
          {showMoreFields && <Row gutter={{ xs: 16, sm: 32, md: 64, lg: 128 }}>
            <Col lg={12} md={12} sm={24} >
              <ProFormText
                tooltip={{ color: colorPrimary, title: intl.formatMessage({ id: 'model.provider.oidc.issuer.description' }) }}
                label={intl.formatMessage({ id: 'model.provider.oidc.issuer' })}
                name="issuer"
                rules={[
                  { type: 'url', message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.url.right' }) },
                ]}
              />
            </Col>
            <Col lg={12} md={12} sm={24} >
              <ProFormText
                tooltip={{ color: colorPrimary, title: intl.formatMessage({ id: 'model.provider.oidc.userinfoEndpoint.description' }) }}

                label={intl.formatMessage({ id: 'model.provider.oidc.userinfoEndpoint' })}
                name="userinfoEndpoint"
                rules={[
                  { type: 'url', message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.url.right' }) },
                ]}
              />
            </Col>
          </Row>}
          <Row gutter={{ xs: 16, sm: 32, md: 64, lg: 128 }}>
            <Col lg={24} md={24} sm={24} >
              <ProForm.Item isListField style={{ marginBlockEnd: 0 }}
                label={intl.formatMessage({ id: 'model.provider.oidc.scopes' })}
              >
                <ProFormList
                  name="scopes"
                  creatorButtonProps={{
                    creatorButtonText: intl.formatMessage({ id: 'pages.operation.add' }),
                    icon: false,
                    type: 'link',
                    style: { width: 'unset' },
                  }}
                  // isValidateList
                  copyIconProps={false}
                  deleteIconProps={{ tooltipText: intl.formatMessage({ id: 'pages.operation.delete' }) }}
                  itemRender={({ listDom, action }) => (
                    <div
                      style={{
                        display: 'inline-flex',
                        marginInlineEnd: 25,
                      }}
                    >
                      {listDom}
                      {action}
                    </div>
                  )}
                >
                  <ProFormText allowClear={false} width="xs" name='value'
                  />
                </ProFormList>
              </ProForm.Item>
            </Col>
          </Row>
        </Card>
      </PageContainer>
    </ProForm>
  );
};
export default ProviderOidcForm;

import {
  FooterToolbar,
  PageContainer,
  ProForm,
  ProFormText,
  ProFormTextArea,
  ProFormRadio,
  ProFormSelect,
  ProFormList,
  ProFormSwitch,
  ProFormCheckbox,
} from '@ant-design/pro-components';
import { Card, Col, Row, Tabs } from 'antd';
import type { TabsProps } from 'antd';
import type { FC } from 'react';
import { useState } from 'react';
import { useParams, useIntl, useNavigate } from '@umijs/max';
import { getApplication, updateApplication, createApplication } from '@/services/application.api';
import { listUser } from '@/services/user.api';
import type { ApplicationDetail, ApplicationCreate, ApplicationUpdate } from '@/services/application.d';
import type { UserDetailList, UserDetail } from '@/services/user.d';
import { getColorPrimary } from '@/utils/global';

const Update = 'update';
const Create = 'create';
const BaseAddress = "/app";



const AdvancedForm: FC<Record<string, any>> = () => {
  const colorPrimary = getColorPrimary();
  const navigate = useNavigate();
  const intl = useIntl();
  const params = useParams();
  const mode = params.id ? Update : Create;
  const onFinish = async (values: Record<string, any>) => {
    if (values?.redirectUris) {
      values.redirectUris = values.redirectUris.map((item: any) => item.value);
    }
    let detail = {};
    if (mode === Create) {
      detail = await createApplication(values as ApplicationCreate) as ApplicationDetail;
    } else {
      values['id'] = params.id;
      detail = await updateApplication(values as ApplicationUpdate) as ApplicationDetail;
    }
    let user = detail as ApplicationDetail;
    if (user.id) {
      // 跳转到详情页
      navigate(`${BaseAddress}/detail/${user.id}`);
    }
  };


  const onInitData = async () => {
    let initData = { enable: true, } as ApplicationDetail;
    if (mode === Update) {
      const res = await getApplication({ id: params.id });
      initData = res as ApplicationDetail;



    }
    return initData;
  };

  const matchTypes = [
    {
      label: intl.formatMessage({ id: 'model.application.redirect_uri.macth.equal' }),
      value: 'equal',
    },
    {
      label: intl.formatMessage({ id: 'model.application.redirect_uri.macth.regex' }),
      value: 'regex',
    },

    {
      label: intl.formatMessage({ id: 'model.application.redirect_uri.macth.prefix' }),
      value: 'prefix',
    },
    {
      label: intl.formatMessage({ id: 'model.application.redirect_uri.macth.contain' }),
      value: 'contain',
    },
  ];
  const supports = [
    {
      label: 'OIDC',
      value: 'oidc',
    },
    {
      label: 'SAML 2.0',
      value: 'saml2',
    },
    {
      label: 'CAS 3.0',
      value: 'cas3',
    },
  ];

  return (
    <ProForm
      layout="vertical"
      submitter={{
        render: (props, dom) => {
          return (
            <FooterToolbar>
              {dom}
            </FooterToolbar>
          );
        },
      }}
      request={onInitData}
      onFinish={onFinish}
    >
      <PageContainer title={intl.formatMessage({ id: 'menu.application' })} header={{ breadcrumb: {}, onBack: () => navigate(`${BaseAddress}`) }}>
        <Card title={mode === Create ? intl.formatMessage({ id: 'model.create' }) : intl.formatMessage({ id: 'model.update' })}
          bordered={false}>
          <Row gutter={{ xs: 16, sm: 32, md: 64, lg: 128 }}>
            <Col lg={8} md={12} sm={24} >
              <ProFormText
                label={intl.formatMessage({ id: 'model.application.name' })}
                name="name"
                rules={[
                  {
                    required: true,
                    message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.application.name' }),
                  },
                ]}
              />
            </Col>
            <Col lg={8} md={12} sm={24} >
              <ProFormText
                label={intl.formatMessage({ id: 'model.application.code' })}
                name="code"
                rules={[
                  {
                    required: true,
                    message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.application.code' }),
                  },
                ]}
              />
            </Col>
            <Col lg={8} md={8} sm={24}>
              <ProFormSwitch
                tooltip={{ color: colorPrimary, title: intl.formatMessage({ id: 'model.application.enable.description' }) }}
                label={intl.formatMessage({ id: 'model.application.enable' })}
                name="enable"
              />
            </Col>



          </Row>
          <Row gutter={{ xs: 16, sm: 32, md: 64, lg: 128 }}>
            <Col lg={12} md={12} sm={24}>
              <ProFormText
                label={intl.formatMessage({ id: 'model.application.redirectUri' })}
                name="redirectUri"
                rules={[
                  {
                    required: true,
                    message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.application.redirectUri' }),
                  },
                ]}
              />
            </Col>
            <Col lg={8} md={12} sm={24}   >
              <ProFormSelect
                label={intl.formatMessage({ id: 'model.application.redirectUriMatchType' })}
                name="redirectUriMatchType"
                rules={[
                  {
                    required: true,
                    message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.application.redirectUriMatchType' }),
                  },
                ]}
                options={matchTypes}
              />
            </Col>
            <Col lg={12} md={12} sm={24}>
              <ProFormText
                label={intl.formatMessage({ id: 'model.application.home' })}
                name="home"
                rules={[
                  {
                    required: true,
                    message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.application.home' }),
                  },
                  { type: 'url', message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.url.right' }) },
                ]}
              />
            </Col>
          </Row>
          <Row gutter={{ xs: 16, sm: 32, md: 64, lg: 128 }}>
            <Col lg={24} md={24} sm={24}>
              <ProFormTextArea
                label={intl.formatMessage({ id: 'model.application.description' })}
                name="description"
                rules={[
                  {
                    required: true,
                    message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.application.description' }),
                  },
                ]}
              />
            </Col>
          </Row>
        </Card>

      </PageContainer>
    </ProForm>
  );
};
export default AdvancedForm;

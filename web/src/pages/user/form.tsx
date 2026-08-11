import {
  FooterToolbar,
  PageContainer,
  ProForm,
  ProFormSelect,
  ProFormText,
} from '@ant-design/pro-components';
import { Card, Col, Row } from 'antd';
import type { FC } from 'react';
import { useParams, useIntl, useNavigate } from '@umijs/max';
import { getUser, updateUser, createUser } from '@/services/user.api';
import { UserDetail, UserUpdate, UserCreate } from '@/services/user.d';
const Update = 'update';
const Create = 'create';
const BaseAddress = "/user";

const AdvancedForm: FC<Record<string, any>> = () => {
  const navigate = useNavigate();
  const intl = useIntl();
  const params = useParams();
  const mode = params.id ? Update : Create;

  const onFinish = async (values: Record<string, any>) => {
    let detail = {};
    if (mode === Create) {
      detail = await createUser(values as UserCreate);
    } else {
      values['id'] = params.id;
      detail = await updateUser(values as UserUpdate);
    }
    let user = detail as UserDetail;
    if (user.id) {
      // 跳转到详情页
      navigate(`${BaseAddress}/detail/${user.id}`);
    }
  };
  const roles = [
    {
      label: intl.formatMessage({ id: 'model.user.role.admin' }),
      value: 'admin',
    },
    {
      label: intl.formatMessage({ id: 'model.user.role.edit' }),
      value: 'edit',
    },
    {
      label: intl.formatMessage({ id: 'model.user.role.view' }),
      value: 'view',
    },
    {
      label: intl.formatMessage({ id: 'model.user.role.none' }),
      value: 'none',
    },
  ];
  const languages = [
    {
      label: intl.formatMessage({ id: 'model.user.language.zh' }),
      value: 'zh-CN',
    },
    {
      label: intl.formatMessage({ id: 'model.user.language.en' }),
      value: 'en-US',
    },
  ]
  const onInitData = async () => {
    let initData = {} as UserDetail;
    if (mode !== Create) {
      const res = await getUser({ id: params.id });
      initData = res as UserDetail;
    }
    return initData;
  };
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
      <PageContainer title={intl.formatMessage({ id: 'menu.user' })} header={{ breadcrumb: {}, onBack: () => navigate(`${BaseAddress}`) }}>
        <Card title={mode === Create ? intl.formatMessage({ id: 'model.user.create' }) : intl.formatMessage({ id: 'model.user.update' })}
          bordered={false}>
          <Row gutter={{ xs: 16, sm: 32, md: 64, lg: 128 }}>
            <Col lg={8} md={12} sm={24} >
              <ProFormText
                label={intl.formatMessage({ id: 'model.user.username' })}
                name="username"
                rules={[
                  {
                    required: true,
                    message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.user.username' }),
                  },
                ]}
              />
            </Col>
            <Col lg={8} md={12} sm={24} >
              <ProFormText
                label={intl.formatMessage({ id: 'model.user.nickname' })}
                name="nickname"
                rules={[
                  {
                    required: true,
                    message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.user.nickname' }),
                  },
                ]}
              />
            </Col>
            <Col lg={8} md={12} sm={24} >
              <ProFormSelect
                label={intl.formatMessage({ id: 'model.user.role' })}
                name="role"
                rules={[
                  {
                    required: true,
                    message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.user.username' }),
                  },
                ]}
                options={roles}
              />
            </Col>
          </Row>
          <Row gutter={{ xs: 16, sm: 32, md: 64, lg: 128 }}>
            <Col lg={8} md={12} sm={24} >
              <ProFormText
                label={intl.formatMessage({ id: 'model.user.email' })}
                name="email"
                rules={[
                  {
                    type: 'email',
                    message: intl.formatMessage({ id: 'pages.user.email.invaid' }),
                  },
                  {
                    required: true,
                    message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.user.email' }),
                  },
                ]}
              />
            </Col>
            <Col lg={8} md={12} sm={24} >
              <ProFormText
                label={intl.formatMessage({ id: 'model.user.phone' })}
                name="phone"
                rules={[
                  {
                    pattern: /^1\d{10}$/,
                    message: intl.formatMessage({ id: 'pages.user.phone.invaid' }),
                  },
                  {
                    required: true,
                    message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.user.phone' }),
                  },
                ]}
              />
            </Col>
            <Col lg={8} md={12} sm={24} >
              <ProFormText
                label={intl.formatMessage({ id: 'model.user.jobNumber' })}
                name="jobNumber"
              />
            </Col>
          </Row>
          <Row gutter={{ xs: 16, sm: 32, md: 64, lg: 128 }}>
            <Col lg={8} md={12} sm={24} >
              <ProFormSelect
                label={intl.formatMessage({ id: 'model.user.language' })}
                name="language"
                options={languages}
              />

            </Col>
            {mode === Create && <Col lg={8} md={12} sm={24} >
              <ProFormText
                label={intl.formatMessage({ id: 'model.user.password' })}
                name="password"
              />
            </Col>
            }
          </Row>
        </Card>

      </PageContainer>
    </ProForm>
  );
};
export default AdvancedForm;

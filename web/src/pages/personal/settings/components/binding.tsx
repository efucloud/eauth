import { Divider, Row, Card, Col, List, Popconfirm, message, } from 'antd';
import React, { Fragment, useEffect, useState } from 'react';
import { UserAuthProfileDetail, UserAuthProfileDetailList } from '@/services/user_auth_profile.d'
import { getUserPersonalAuthProfile, changeUserAuthProfileStatus, deleteUserAuthProfile } from '@/services/auth_profile.api'
import { ProviderTwoLine, Provider } from '@/components';
import { FormattedMessage, useIntl, useModel } from '@umijs/max';
import { IntlShape } from "react-intl";
import { CloseCircleOutlined, CheckCircleOutlined } from '@ant-design/icons'
import { ThirdAuthMethod } from '@/services/common.d';
import { thirdAuthMethods } from '@/services/oauth.api';
import { addBind, getColorPrimary } from '@/utils/global';
import { Typography } from 'antd';
const { Text, Paragraph, Link } = Typography;
import dayjs from 'dayjs';

const BindingView: React.FC = () => {

  const intl = useIntl();
  const [authedProfiles, setAuthedProfiles] = useState<UserAuthProfileDetail[]>([]);
  const [authMethod, setAuthMethod] = useState<ThirdAuthMethod>();
  const colorPrimary = getColorPrimary();
  const getCurrentUserAuthedProfiles = async () => {
    const data = await getUserPersonalAuthProfile();
    const list = data as UserAuthProfileDetailList;
    if (list.data) {
      setAuthedProfiles(list.data);
    }
  }
  const getAuthMethods = async () => {
    setAuthMethod(await thirdAuthMethods());
  }


  useEffect(() => {
    getCurrentUserAuthedProfiles();
    getAuthMethods();
  },
    [])

  const handleEnable = async (intl: IntlShape, enable: boolean, selectedRows: UserAuthProfileDetail[]) => {
    const hide = message.loading(intl.formatMessage({ id: 'pages.operation.updating' }));
    if (!selectedRows) return true;

    try {
      await changeUserAuthProfileStatus({
        ids: selectedRows.map((row) => row.id),
        enable: enable,
      });
      hide();
      message.success(intl.formatMessage({ id: 'pages.operation.update.success' }));
      getCurrentUserAuthedProfiles();
      return true;
    } catch (error) {
      hide();
      message.error(intl.formatMessage({ id: 'pages.operation.update.failed' }));
      return false;
    }
  };
  const handleDelete = async (intl: IntlShape, selectedRows: UserAuthProfileDetail[]) => {
    const hide = message.loading(intl.formatMessage({ id: 'pages.operation.updating' }));
    if (!selectedRows) return true;

    try {
      await deleteUserAuthProfile({
        ids: selectedRows.map((row) => row.id),
      });
      hide();
      message.success(intl.formatMessage({ id: 'pages.operation.update.success' }));
      getCurrentUserAuthedProfiles();
      return true;
    } catch (error) {
      hide();
      message.error(intl.formatMessage({ id: 'pages.operation.update.failed' }));
      return false;
    }
  };

  const BindingViewItem = (): React.ReactNode => {
    const { initialState } = useModel('@@initialState');

    return <Row >
      {authMethod?.oidcs?.map((provider) => <Col key={provider.category} lg={2} md={4} sm={8} xs={12} onClick={() => {
        addBind(initialState?.currentUser?.id)
        window.location.href = provider.address;
      }}><ProviderTwoLine provider={provider.category} /></Col>)}
    </Row>;
  }

  return (
    <Fragment>
      {BindingViewItem()}
      <Divider style={{ color: colorPrimary }}></Divider>
      <List
        itemLayout="horizontal"
        grid={{ gutter: 16, column: 3 }}
        dataSource={authedProfiles}
        renderItem={(item) => (
          <List.Item>
            <Card
              title={<div style={{ display: 'flex', alignItems: 'center', justifyContent: 'flex-start' }}>
                < Provider provider={item.provider ? item.provider : ''} height={40} />
                &nbsp;&nbsp; {item.enable ? <CheckCircleOutlined style={{ color: 'green' }} /> : <CloseCircleOutlined style={{ color: 'red' }} />}
              </div>}
              extra={<>
                <Popconfirm
                  key={item.id + '-unbind'}
                  title=''
                  description={intl.formatMessage({ id: 'pages.personal.settings.unbind.description' })}

                  okText={intl.formatMessage({ id: 'pages.operation.confirm' })}
                  cancelText={intl.formatMessage({ id: 'pages.operation.cancel' })}
                  onConfirm={() => {
                    handleDelete(intl, [item]);
                  }}
                >
                  <a><FormattedMessage id="pages.personal.settings.unbind" defaultMessage="解绑" /></a>
                </Popconfirm>&nbsp;&nbsp;

                {item.enable === true && <Popconfirm
                  key={item.id + '-disable'}
                  title=''
                  description={intl.formatMessage({ id: 'pages.personal.settings.bind.disable.description' })}
                  okText={intl.formatMessage({ id: 'pages.operation.confirm' })}
                  cancelText={intl.formatMessage({ id: 'pages.operation.cancel' })}
                  onConfirm={() => {
                    handleEnable(intl, false, [item]);
                  }}
                >
                  <a style={{ color: colorPrimary }} className='delete'><FormattedMessage id="pages.personal.settings.bind.disable" defaultMessage="禁用" /></a>
                </Popconfirm>}
                {item.enable === false && <Popconfirm
                  key={item.id + '-enable'}
                  title=''
                  description={intl.formatMessage({ id: 'pages.personal.settings.bind.enable.description' })}
                  okText={intl.formatMessage({ id: 'pages.operation.confirm' })}
                  cancelText={intl.formatMessage({ id: 'pages.operation.cancel' })}
                  onConfirm={() => {
                    handleEnable(intl, true, [item]);
                  }}
                >
                  <a style={{ color: colorPrimary }} ><FormattedMessage id="pages.personal.settings.bind.enable" defaultMessage="启用" /></a>
                </Popconfirm>}
              </>}
            >
              <Card.Meta
                description={
                  <Paragraph>
                    <Text strong><FormattedMessage id='pages.personal.settings.oidc.user.name' />:&nbsp;&nbsp;</Text><Text>{item.loginName}</Text>
                    <br />
                    <Text strong><FormattedMessage id='pages.personal.settings.oidc.user.home' />:&nbsp;&nbsp;</Text><Link href={item.home}>{item.home}</Link>
                    <br />
                    <Text strong><FormattedMessage id='pages.personal.settings.oidc.bind.time' />:&nbsp;&nbsp;</Text><Text>{dayjs(item.createdAt).format('YYYY-MM-DD HH:mm:ss')}</Text>
                  </Paragraph>}
              />
            </Card>

          </List.Item>
        )}
      />
    </Fragment>
  );
};

export default BindingView;

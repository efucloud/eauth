import { PageContainer, ProDescriptions, ProTable, ActionType } from '@ant-design/pro-components';
import { useParams, useIntl, useNavigate, FormattedMessage, Access, useAccess } from '@umijs/max';
import { Card, Space, Button, Col, Row, Popconfirm, message, Avatar, Divider } from 'antd';
import { EditOutlined } from '@ant-design/icons';
import type { TableListPagination } from '@/services/common.d';
import type { ProColumns } from '@ant-design/pro-components';
import React, { useEffect, useState, useRef } from 'react';
import type { UserDetail } from '@/services/user.d';
import { getUser, } from '@/services/user.api';
import { getUserAuthProfileByUserId, changeUserAuthProfileStatus, deleteUserAuthProfile } from '@/services/auth_profile.api';
import { UserAuthProfileDetail } from '@/services/user_auth_profile.d';
import { ProviderInLine } from '@/components';
import { IntlShape } from "react-intl";
import { listUserFaceRecognition, deleteUserFaceRecognition, changeUserFaceRecognitionStatus } from '@/services/face_recognition.api';
import { FaceRecognitionDetail } from '@/services/face_recognition.d';
import { getColorPrimary } from '@/utils/global';
const BaseAddress = "/user";

const UserDetail: React.FC = () => {
  const colorPrimary = getColorPrimary();
  const access = useAccess();
  const navigate = useNavigate();
  const intl = useIntl();
  const params = useParams();
  const authProfileActionRef = useRef<ActionType>();
  const faceIdDataActionRef = useRef<ActionType>();
  const [userInfo, setUserInfo] = useState<UserDetail>();
  const getUserInfo = async () => {
    const data = await getUser({ id: params.id });
    setUserInfo(data as UserDetail);
  }
  useEffect(() => {
    getUserInfo();
  }, [params.id]);


  const authProfileHandleEnable = async (intl: IntlShape, enable: boolean, selectedRows: UserAuthProfileDetail[]) => {
    const hide = message.loading(intl.formatMessage({ id: 'pages.operation.updating' }));
    if (!selectedRows) return true;

    try {
      await changeUserAuthProfileStatus({
        ids: selectedRows.map((row) => row.id),
        enable: enable,
      });
      hide();
      message.success(intl.formatMessage({ id: 'pages.operation.update.success' }));
      authProfileActionRef.current?.reloadAndRest?.();
      return true;
    } catch (error) {
      hide();
      message.error(intl.formatMessage({ id: 'pages.operation.update.failed' }));
      return false;
    }
  };
  const authProfileHandleDelete = async (intl: IntlShape, selectedRows: UserAuthProfileDetail[]) => {
    const hide = message.loading(intl.formatMessage({ id: 'pages.operation.updating' }));
    if (!selectedRows) return true;

    try {
      await deleteUserAuthProfile({
        ids: selectedRows.map((row) => row.id),
      });
      hide();
      message.success(intl.formatMessage({ id: 'pages.operation.update.success' }));
      authProfileActionRef.current?.reloadAndRest?.();
      return true;
    } catch (error) {
      hide();
      message.error(intl.formatMessage({ id: 'pages.operation.update.failed' }));
      return false;
    }
  };
  const userAuthProfileColumns: ProColumns<UserAuthProfileDetail>[] = [
    {
      title: intl.formatMessage({ id: 'model.user_auth_profile.provider' }),
      dataIndex: 'provider',
      render: (dom, entity) => { return (<ProviderInLine provider={entity.provider} />); },
    },
    {
      title: intl.formatMessage({ id: 'model.user_auth_profile.createdAt' }),
      dataIndex: 'createdAt',
      search: false,
      valueType: 'dateTime',
    },
    {
      title: intl.formatMessage({ id: 'model.latestUsedTime' }),
      dataIndex: 'latestUsedTime',
      search: false,
      sorter: true,
      valueType: 'dateTime',
    },
    {
      title: intl.formatMessage({ id: 'model.user_auth_profile.enable' }),
      dataIndex: 'enable',

      search: false,
      valueEnum: {
        false: {
          text: intl.formatMessage({ id: 'model.user_auth_profile.enable.disable' }),
          status: 'Error',
        },
        true: {
          text: intl.formatMessage({ id: 'model.user_auth_profile.enable.enable' }),
          status: 'Success',
        },
      },
    },
    {
      title: intl.formatMessage({ id: 'pages.operation' }),
      dataIndex: 'option',
      valueType: 'option',
      hideInTable: !access.systemEditAccess,
      render: (_, item) => {
        let nodes = [] as React.ReactNode[];
        nodes.push(
          <Popconfirm
            key={item.id + '-unbind'}
            title=''
            description={intl.formatMessage({ id: 'pages.personal.settings.unbind.description' })}

            okText={intl.formatMessage({ id: 'pages.operation.confirm' })}
            cancelText={intl.formatMessage({ id: 'pages.operation.cancel' })}
            onConfirm={() => {
              authProfileHandleDelete(intl, [item]);
            }}
          >
            <a><FormattedMessage id="pages.personal.settings.unbind" defaultMessage="解绑" /></a>
          </Popconfirm>
        );
        if (item.enable === true) {
          nodes.push(
            <Popconfirm
              key={item.id + '-disable'}
              title=''
              description={intl.formatMessage({ id: 'pages.personal.settings.bind.disable.description' })}
              okText={intl.formatMessage({ id: 'pages.operation.confirm' })}
              cancelText={intl.formatMessage({ id: 'pages.operation.cancel' })}
              onConfirm={() => {
                authProfileHandleEnable(intl, false, [item]);
              }}
            >
              <a style={{ color: colorPrimary }} className='delete'><FormattedMessage id="pages.personal.settings.bind.disable" defaultMessage="禁用" /></a>
            </Popconfirm>
          );
        }
        if (item.enable === false) {
          nodes.push(<Popconfirm
            key={item.id + '-enable'}
            title=''
            description={intl.formatMessage({ id: 'pages.personal.settings.bind.enable.description' })}
            okText={intl.formatMessage({ id: 'pages.operation.confirm' })}
            cancelText={intl.formatMessage({ id: 'pages.operation.cancel' })}
            onConfirm={() => {
              authProfileHandleEnable(intl, true, [item]);
            }}
          >
            <a style={{ color: colorPrimary }} ><FormattedMessage id="pages.personal.settings.bind.enable" defaultMessage="启用" /></a>
          </Popconfirm>);
        }
        return nodes;
      },
    },

  ];

  const UserAuthProfiles = (userId: number) => {
    return (
      <ProTable<UserAuthProfileDetail, TableListPagination>
        rowKey="id"
        actionRef={authProfileActionRef}
        search={false}
        params={{ userId: userId }}
        toolBarRender={false}
        pagination={false}
        request={getUserAuthProfileByUserId}
        columns={userAuthProfileColumns}
      >

      </ProTable>
    );
  };
  const faceIdDataHandleEnable = async (intl: IntlShape, enable: boolean, selectedRows: FaceRecognitionDetail[]) => {
    const hide = message.loading(intl.formatMessage({ id: 'pages.operation.updating' }));
    if (!selectedRows) return true;
    try {
      await changeUserFaceRecognitionStatus({
        ids: selectedRows.map((row) => row.id),
        enable: enable,
      });
      hide();
      message.success(intl.formatMessage({ id: 'pages.operation.update.success' }));
      faceIdDataActionRef.current?.reloadAndRest?.();
      return true;
    } catch (error) {
      hide();
      message.error(intl.formatMessage({ id: 'pages.operation.update.failed' }));
      return false;
    }
  };
  const faceIdDataHandleDelete = async (intl: IntlShape, selectedRows: FaceRecognitionDetail[]) => {
    const hide = message.loading(intl.formatMessage({ id: 'pages.operation.updating' }));
    if (!selectedRows) return true;

    try {
      await deleteUserFaceRecognition({
        ids: selectedRows.map((row) => row.id),
      });
      hide();
      message.success(intl.formatMessage({ id: 'pages.operation.update.success' }));
      faceIdDataActionRef.current?.reloadAndRest?.();
      return true;
    } catch (error) {
      hide();
      message.error(intl.formatMessage({ id: 'pages.operation.update.failed' }));
      return false;
    }
  };
  const userFaceIdColumns: ProColumns<FaceRecognitionDetail>[] = [
    {
      title: intl.formatMessage({ id: 'model.face_recognition.name' }),
      dataIndex: 'name',
      search: false,
    },
    {
      title: intl.formatMessage({ id: 'model.face_recognition.createdAt' }),
      dataIndex: 'createdAt',
      search: false,
      sorter: true,
      valueType: 'dateTime',
    },
    {
      title: intl.formatMessage({ id: 'model.latestUsedTime' }),
      dataIndex: 'latestUsedTime',
      search: false,
      sorter: true,
      valueType: 'dateTime',
    },
    {
      title: intl.formatMessage({ id: 'pages.operation' }),
      dataIndex: 'option',
      valueType: 'option',
      hideInTable: !access.systemEditAccess,
      render: (_, item) => {
        let nodes = [] as React.ReactNode[];
        nodes.push(
          <Popconfirm
            key={item.id + '-unbind'}
            title=''
            description={intl.formatMessage({ id: 'pages.personal.settings.unbind.description' })}

            okText={intl.formatMessage({ id: 'pages.operation.confirm' })}
            cancelText={intl.formatMessage({ id: 'pages.operation.cancel' })}
            onConfirm={() => {
              faceIdDataHandleDelete(intl, [item]);
            }}
          >
            <a><FormattedMessage id="pages.personal.settings.unbind" defaultMessage="解绑" /></a>
          </Popconfirm>
        );
        if (item.enable === true) {
          nodes.push(
            <Popconfirm
              key={item.id + '-disable'}
              title=''
              description={intl.formatMessage({ id: 'pages.personal.settings.bind.disable.description' })}
              okText={intl.formatMessage({ id: 'pages.operation.confirm' })}
              cancelText={intl.formatMessage({ id: 'pages.operation.cancel' })}
              onConfirm={() => {
                faceIdDataHandleEnable(intl, false, [item]);
              }}
            >
              <a style={{ color: colorPrimary }} className='delete'><FormattedMessage id="pages.personal.settings.bind.disable" defaultMessage="禁用" /></a>
            </Popconfirm>
          );
        }
        if (item.enable === false) {
          nodes.push(<Popconfirm
            key={item.id + '-enable'}
            title=''
            description={intl.formatMessage({ id: 'pages.personal.settings.bind.enable.description' })}
            okText={intl.formatMessage({ id: 'pages.operation.confirm' })}
            cancelText={intl.formatMessage({ id: 'pages.operation.cancel' })}
            onConfirm={() => {
              faceIdDataHandleEnable(intl, true, [item]);
            }}
          >
            <a style={{ color: colorPrimary }} ><FormattedMessage id="pages.personal.settings.bind.enable" defaultMessage="启用" /></a>
          </Popconfirm>);
        }

        return nodes;
      },
    },
  ]
  const UserFaceIdDatas = (userId: number) => {
    return (
      <ProTable<FaceRecognitionDetail, TableListPagination>
        rowKey="id"
        actionRef={faceIdDataActionRef}
        search={false}
        params={{ userId: userId }}
        toolBarRender={false}
        pagination={false}
        request={listUserFaceRecognition}
        columns={userFaceIdColumns}
      >
      </ProTable>
    );
  };

  return (
    <PageContainer title={intl.formatMessage({ id: 'menu.user' })} header={{ breadcrumb: {}, onBack: () => navigate(`${BaseAddress}`) }}
      subTitle={userInfo?.avatar && <Avatar src={userInfo?.avatar} />} >
      <Card bordered={false}
      >
        <ProDescriptions style={{ marginBottom: 32 }} column={3} title={intl.formatMessage({ id: 'pages.detail.baseinfo' })}
          extra={
            <Space>
              <Access accessible={access.systemEditAccess === true}>
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
          <ProDescriptions.Item label={intl.formatMessage({ id: 'model.user.username' })}>{userInfo?.username}</ProDescriptions.Item>
          <ProDescriptions.Item label={intl.formatMessage({ id: 'model.user.nickname' })}>{userInfo?.nickname}</ProDescriptions.Item>
          <ProDescriptions.Item label={intl.formatMessage({ id: 'model.user.role' })}
            valueEnum={{
              'admin': {
                text: intl.formatMessage({ id: 'model.user.role.admin' }),
              },
              'edit': {
                text: intl.formatMessage({ id: 'model.user.role.edit' }),
              },
              'view': {
                text: intl.formatMessage({ id: 'model.user.role.view' }),
              },
              'none': {
                text: intl.formatMessage({ id: 'model.user.role.none' }),
              },
            }}
          >{userInfo?.role}</ProDescriptions.Item>

          <ProDescriptions.Item label={intl.formatMessage({ id: 'model.user.email' })}>{userInfo?.email}</ProDescriptions.Item>
          <ProDescriptions.Item label={intl.formatMessage({ id: 'model.user.phone' })}>{userInfo?.phone}</ProDescriptions.Item>
          <ProDescriptions.Item label={intl.formatMessage({ id: 'model.user.jobNumber' })}>{userInfo?.jobNumber}</ProDescriptions.Item>
          <ProDescriptions.Item label={intl.formatMessage({ id: 'model.user.enable' })}
            valueEnum={{
              false: {
                text: intl.formatMessage({ id: 'model.user.enable.disable' }),
                status: 'Error',
              },
              true: {
                text: intl.formatMessage({ id: 'model.user.enable.enable' }),
                status: 'Success',
              },
            }}
          >{userInfo?.enable}</ProDescriptions.Item>
        </ProDescriptions>
      </Card>
      <br />
      <Row gutter={24}>

        <Col xl={12} lg={24} md={24} sm={24} xs={24}>
          <Card bordered={false} title={intl.formatMessage({ id: 'model.user.authed.third.method' })}>
            <Divider style={{ marginTop: 10, borderColor: '#7cb305' }} ><FormattedMessage id='model.user.authed.third.method.oidc' /></Divider>
            {userInfo?.id && UserAuthProfiles(userInfo?.id)}
            <Divider style={{ marginTop: 10, borderColor: '#7cb305' }} ><FormattedMessage id='model.user.authed.third.method.face' /></Divider>
            {userInfo?.id && UserFaceIdDatas(userInfo?.id)}
          </Card>
        </Col>
      </Row>

    </PageContainer>
  );
};
export default UserDetail;




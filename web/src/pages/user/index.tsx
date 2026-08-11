import { PlusOutlined, DownOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { RequestOptionsType } from '@ant-design/pro-utils';
import {
  FooterToolbar,
  PageContainer,
  ProTable,
  ModalForm,
  ProFormText,
  ProFormSelect
} from '@ant-design/pro-components';
import { Button, message, Popconfirm, Dropdown, Space, Alert, Avatar, Modal, Row, Col } from 'antd';
import React, { useRef, useState } from 'react';
import type { UserDetail, UserResetPassword,  UserRole } from '@/services/user.d';
import type { TableListPagination } from '@/services/common.d';
import { deleteUser, changeUserStatus, listUser, resetUserPassword, setUserRole } from '@/services/user.api';
import { FormattedMessage, useIntl, Access, useAccess } from '@umijs/max';
import { useNavigate } from '@umijs/max';
import { IntlShape } from "react-intl";
import { getColorPrimary } from '@/utils/global';
const BaseAddress = "/user";


const GLobalUserTableList: React.FC = () => {
  const colorPrimary=getColorPrimary();
  const access = useAccess();
  /** 新建窗口的弹窗 */
  const actionRef = useRef<ActionType>();
  const [selectedRowsState, setSelectedRows] = useState<UserDetail[]>([]);
  const [userInfo, setUserinfo] = useState<UserDetail>();
  const [resetPasswordModalVisible, setResetPasswordModalVisible] = useState<boolean>(false);
  const [roleModalVisible, setRoleModalVisible] = useState<boolean>(false);
 
  /** 国际化配置 */
  const intl = useIntl();
  const navigate = useNavigate();

  const handleRemove = async (intl: IntlShape, selectedRows: UserDetail[]) => {
    const hide = message.loading(intl.formatMessage({ id: 'pages.operation.deleting' }));
    if (!selectedRows) return true;

    try {
      await deleteUser({
        ids: selectedRows.map((row) => row.id),
      });
      hide();
      message.success(intl.formatMessage({ id: 'pages.operation.delete.success' }));
      actionRef.current?.reload();
      return true;
    } catch (error) {
      hide();
      message.error(intl.formatMessage({ id: 'pages.operation.delete.failed' }));
      return false;
    }
  };

  const handleEnable = async (intl: IntlShape, enable: boolean, selectedRows: UserDetail[]) => {
    const hide = message.loading(intl.formatMessage({ id: 'pages.operation.updating' }));
    if (!selectedRows) return true;

    try {
      await changeUserStatus({
        ids: selectedRows.map((row) => row.id),
        enable: enable,
      });
      hide();
      message.success(intl.formatMessage({ id: 'pages.operation.update.success' }));
      actionRef.current?.reload();
      return true;
    } catch (error) {
      hide();
      message.error(intl.formatMessage({ id: 'pages.operation.update.failed' }));
      return false;
    }
  };
 
  const moreItems = (record: UserDetail) => {
    let nodes = [];

    if (record.enable === false) {
      nodes.push({
        key: record.id + '-enable',
        label: (<Popconfirm
          key={record.id + '-enable'}
          description={intl.formatMessage({ id: 'model.user.enable.enable.description' })}
          title={intl.formatMessage({ id: 'pages.operation.enable.description' }) +
            intl.formatMessage({ id: 'user.model.name' }) + '【' + record.username + '】'}
          onConfirm={() => {
            handleEnable(intl, true, [record]);
          }}
          okText={intl.formatMessage({ id: 'pages.operation.confirm' })}
          cancelText={intl.formatMessage({ id: 'pages.operation.cancel' })}
        >
          <a style={{color:colorPrimary}} className='operation-enable'><FormattedMessage id="pages.operation.enable" defaultMessage="启用" /></a>
        </Popconfirm>)
      });
    } else {
      nodes.push({
        key: record.id + '-disable',
        label: (<Popconfirm
          key={record.id + '-disable'}
          description={intl.formatMessage({ id: 'model.user.enable.disable.description' })}
          title={intl.formatMessage({ id: 'pages.operation.disable.description' }) +
            intl.formatMessage({ id: 'user.model.name' }) + '【' + record.username + '】'}
          onConfirm={() => {
            handleEnable(intl, false, [record]);
          }}
          okText={intl.formatMessage({ id: 'pages.operation.confirm' })}
          cancelText={intl.formatMessage({ id: 'pages.operation.cancel' })}
        >
          <a style={{color:colorPrimary}} className='operation-disable'><FormattedMessage id="pages.operation.disable" defaultMessage="禁用" /></a>
        </Popconfirm>)
      });
    };
   
    nodes.push({
      key: record.id + '-delete',
      label: (<Popconfirm
        key={record.id + '-delete'}
        description={intl.formatMessage({ id: 'model.user.delete.description' })}
        title={intl.formatMessage({ id: 'pages.operation.delete.description' }) +
          intl.formatMessage({ id: 'user.model.name' }) + '【' + record.username + '】'}
        onConfirm={() => {
          handleRemove(intl, [record])
        }}
        okText={intl.formatMessage({ id: 'pages.operation.confirm' })}
        cancelText={intl.formatMessage({ id: 'pages.operation.cancel' })}
      >
        <a style={{color:colorPrimary}} className='operation-delete'><FormattedMessage id="pages.operation.delete"  /></a>
      </Popconfirm>)
    });
    nodes.push({
      key: 'set-role',
      label: (<a
        key="set-role"
        onClick={() => {
          setRoleModalVisible(true);
          setUserinfo(record);
        }}
      >
        <FormattedMessage id="model.user.set.role" defaultMessage="设置角色" />
      </a>)
    });
    nodes.push({
      key: 'reset-password',
      label: (<a
        key="reset-password"
        onClick={() => {
          setUserinfo(record);
          setResetPasswordModalVisible(true);
        }}
      >
        <FormattedMessage id="model.user.reset.password" defaultMessage="重置密码" />
      </a>)
    });
    return nodes;
  };
  const columns: ProColumns<UserDetail>[] = [
    {
      title: intl.formatMessage({ id: 'model.user.username' }),
      dataIndex: 'username',
      render: (dom, entity) => {
        return (
          <>
            {entity.avatar && <Avatar src={entity.avatar} size={24} />}&nbsp;&nbsp;
            <a
              onClick={() => {
                navigate(`${BaseAddress}/detail/${entity.id}`)
              }}
            >
              {dom}
            </a>
          </>
        );
      },
    },
    {
      title: intl.formatMessage({ id: 'model.user.email' }),
      dataIndex: 'email',
      valueType: 'text',
    },
    {
      title: intl.formatMessage({ id: 'model.user.phone' }),
      dataIndex: 'phone',
      valueType: 'text',
    },
    {
      title: intl.formatMessage({ id: 'model.user.role' }),
      dataIndex: 'role',
      
      valueEnum: {
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
      },
    },
    {
      title: intl.formatMessage({ id: 'model.user.jobNumber' }),
      dataIndex: 'job_number',
      valueType: 'text',
    },


    {
      title: intl.formatMessage({ id: 'model.user.enable' }),
      dataIndex: 'enable',
      
      search: false,
      valueEnum: {
        false: {
          text: intl.formatMessage({ id: 'model.user.enable.disable' }),
          status: 'Error',
        },
        true: {
          text: intl.formatMessage({ id: 'model.user.enable.enable' }),
          status: 'Success',
        },
      },
    },
    {
      title: intl.formatMessage({ id: 'model.createdAt' }),
      dataIndex: 'createdAt',
      search: false,
      sorter: true,
      valueType: 'dateTime',

    },
    {
      title: intl.formatMessage({ id: 'pages.operation' }),
      dataIndex: 'option',
      valueType: 'option',
      hideInTable: !access.systemEditAccess,
      render: (_, record) => {
        const nodes = [
          <a
            key="edit"
            onClick={() => {
              navigate(
                {
                  pathname: `${BaseAddress}/update/${record.id}`,
                },
                { replace: true },
              );
            }}
          >
            <FormattedMessage id="pages.operation.edit" />
          </a>,

        ];

        nodes.push(
          <Dropdown menu={{ items: moreItems(record) }} key='more' >
            <a style={{color:colorPrimary}} onClick={(e) => {
              setUserinfo(record);
              e.preventDefault();
            }}>
              <Space>
                <FormattedMessage id="pages.operation.more" defaultMessage="更多" />
                <DownOutlined />
              </Space>
            </a>
          </Dropdown>

        );

        return nodes;
      },
    },
  ];
  
 
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
  return (

    <PageContainer
      title={intl.formatMessage({ id: 'pages.user.list.title' })}
      subTitle={intl.formatMessage({ id: 'pages.user.list.subTitle' })}
    >
      <ProTable<UserDetail, TableListPagination>
        actionRef={actionRef}
        rowKey="id"
        search={{
           showHiddenNum:true,
           span:{xs: 24, sm: 12, md: 8, lg: 6, xl: 6, xxl: 6}
         }}

        toolBarRender={() => [
          <Access accessible={access.systemAdminAccess} key={'create'}>
            <Button
              type="primary"
              key="primary"
              onClick={() => {
                navigate(`${BaseAddress}/create`)
              }}
            >
              <PlusOutlined /> <FormattedMessage id="pages.operation.create" defaultMessage="新建" />
            </Button>
          </Access>,
        ]}
        request={listUser}
        columns={columns}
        rowSelection={{
          onChange: (_, selectedRows) => {
            setSelectedRows(selectedRows);
          },
        }}
      />
      {selectedRowsState?.length > 0 && (
        <FooterToolbar
          extra={
            <div>
              <FormattedMessage id="pages.operation.selected" defaultMessage="已选择" />{' '}
              <a
                style={{
                  fontWeight: 600,
                }}
              >
                {selectedRowsState.length}
              </a>{' '}
              <FormattedMessage id="pages.operation.selected.term" defaultMessage="项" /> &nbsp;&nbsp;
            </div>
          }
        >
          <Button
            danger
            onClick={() => {
              Modal.confirm({
                title: intl.formatMessage({ id: 'pages.operation.delete.confirm.title' }),
                content: intl.formatMessage({ id: 'pages.operation.delete.confirm.content' }),
                onOk: async () => {
                  await handleRemove(intl, selectedRowsState);
                  setSelectedRows([]);
                  actionRef.current?.reloadAndRest?.();
                }
              });
            }}
          >
            <FormattedMessage id="pages.operation.batch.delete" />
          </Button>
        </FooterToolbar>
      )}
      <ModalForm
        title={intl.formatMessage({ id: 'model.user.reset.password' })}

        width="400px"
        open={resetPasswordModalVisible}
        onOpenChange={setResetPasswordModalVisible}
        onFinish={async (value) => {
          value['id'] = userInfo?.id;
          const res = await resetUserPassword(value as UserResetPassword);
          if (res) {

            message.success(intl.formatMessage({ id: 'model.user.reset.password.success' }));
            setResetPasswordModalVisible(false);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }
        }}
      >
        <Alert message={intl.formatMessage({ id: 'model.user.username' }) + ': ' + userInfo?.username} type="error" />
        <br />
        <ProFormText.Password
          width="md"
          label={intl.formatMessage({ id: 'model.user.newPassword' })}
          name="newPassword"
          rules={[
            {
              required: true,
              message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.user.newPassword' }),
            },
          ]}
        />
        <ProFormText.Password
          width="md"
          label={intl.formatMessage({ id: 'model.user.confirmPassword' })}
          name="confirmPassword"
          rules={[
            {
              required: true,
              message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.user.confirmPassword' }),
            },
            ({ getFieldValue }) => ({
              validator(rule, value) {
                if (!value || getFieldValue('newPassword') === value) {
                  return Promise.resolve();
                }
                return Promise.reject(new Error(intl.formatMessage({ id: 'model.user.password.not.match' })));
              },
            }),
          ]}
          dependencies={['newPassword']}
        />
      </ModalForm>
     
      <ModalForm
        title={intl.formatMessage({ id: 'model.user.set.role' })}

        width="400px"
        open={roleModalVisible}
        onOpenChange={setRoleModalVisible}
        onFinish={async (value) => {
          value['ids'] = [userInfo?.id];
          const res = await setUserRole(value as UserRole);
          if (res) {
            message.success(intl.formatMessage({ id: 'operation.success' }));
            setRoleModalVisible(false);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }
        }}
      >
        <Alert message={intl.formatMessage({ id: 'model.user.username' }) + ': ' + userInfo?.username} type="error" />
        <br />

        <ProFormSelect
          label={intl.formatMessage({ id: 'model.user.role' })}
          name="role"
          rules={[
            {
              required: true,
              message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.user.role' }),
            },
          ]}
          options={roles}
        />
      </ModalForm>
    </PageContainer>

  );

};

export default GLobalUserTableList;

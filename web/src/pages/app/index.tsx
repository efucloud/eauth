import { PlusOutlined, DownOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import {
  FooterToolbar,
  PageContainer,
  ProTable,
} from '@ant-design/pro-components';
import { Button, message, Popconfirm, Dropdown, Space,Modal } from 'antd';
import { Typography } from 'antd';
const { Text, Paragraph, Link } = Typography;
import React, { useRef, useState } from 'react';
import type { ApplicationDetail} from '@/services/application.d';
import type { TableListPagination } from '@/services/common.d';
import { deleteApplication, listApplication, changeApplicationStatus } from '@/services/application.api';
import { FormattedMessage, useIntl, Access, useAccess } from '@umijs/max';
import { useNavigate } from '@umijs/max';
import { IntlShape } from "react-intl";
import { User } from '@/components';
import { getColorPrimary } from '@/utils/global';
const BaseAddress = "/app";



const ApplicationTableList: React.FC = () => {
  const colorPrimary=getColorPrimary();
  /** 新建窗口的弹窗 */
  const actionRef = useRef<ActionType>();
  const access = useAccess();
  const [selectedRowsState, setSelectedRows] = useState<ApplicationDetail[]>([]);
  /** 国际化配置 */
  const intl = useIntl();
  const navigate = useNavigate();

  const handleRemove = async (intl: IntlShape, selectedRows: ApplicationDetail[]) => {
    const hide = message.loading(intl.formatMessage({ id: 'pages.operation.deleting' }));
    if (!selectedRows) return true;

    try {
      await deleteApplication({
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

  const handleEnable = async (intl: IntlShape, enable: boolean, selectedRows: ApplicationDetail[]) => {
    const hide = message.loading(intl.formatMessage({ id: 'pages.operation.updating' }));
    if (!selectedRows) return true;

    try {
      await changeApplicationStatus({
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
  const moreItems = (record: ApplicationDetail) => {
    let nodes = [];

    if (record.enable === false) {
      nodes.push({
        key: record.id + '-enable',
        label: (<Popconfirm
          key={record.id + '-enable'}
          title=''
          description={intl.formatMessage({ id: 'pages.operation.enable.description'}) +
            intl.formatMessage({ id: 'application.model.name'}) + '【' + record.name + '】'}
          onConfirm={() => {
            handleEnable(intl, true, [record]);
          }}
          okText={intl.formatMessage({ id: 'pages.operation.confirm'})}
          cancelText={intl.formatMessage({ id: 'pages.operation.cancel'})}
        >
          <a style={{color:colorPrimary}} ><FormattedMessage id="pages.operation.enable" defaultMessage="启用" /></a>
        </Popconfirm>)
      });
    } else {
      nodes.push({
        key: record.id + '-disable',
        label: (<Popconfirm
          key={record.id + '-disable'}
          title=''
          description={intl.formatMessage({ id: 'pages.operation.disable.description'}) +
            intl.formatMessage({ id: 'application.model.name'}) + '【' + record.name + '】'}
          onConfirm={() => {
            handleEnable(intl, false, [record]);
          }}
          okText={intl.formatMessage({ id: 'pages.operation.confirm'})}
          cancelText={intl.formatMessage({ id: 'pages.operation.cancel'})}
        >
          <a style={{color:colorPrimary}} ><FormattedMessage id="pages.operation.disable" defaultMessage="禁用" /></a>
        </Popconfirm>)
      });
    };
    nodes.push({
      key: record.id + '-delete',
      label: (<Popconfirm
        key={record.id + '-delete'}
        title=''
        description={intl.formatMessage({ id: 'pages.operation.delete.description'}) +
          intl.formatMessage({ id: 'application.model.name'}) + '【' + record.name + '】'}
        onConfirm={() => {
          handleRemove(intl, [record])
        }}
        okText={intl.formatMessage({ id: 'pages.operation.confirm'})}
        cancelText={intl.formatMessage({ id: 'pages.operation.cancel'})}
      >
        <a style={{color:colorPrimary}} className='delete'><FormattedMessage id="pages.operation.delete"  /></a>
      </Popconfirm>)
    });
    return nodes;
  };
  const columns: ProColumns<ApplicationDetail>[] = [
    {
      title: intl.formatMessage({ id: 'model.application.name'}),
      dataIndex: 'name',
      render: (dom, entity) => {
        return (
          <a
            onClick={() => {
              navigate(`${BaseAddress}/detail/${entity.id}`)
            }}
          >
            {dom}
          </a>
        );
      },
    },
    {
      title: intl.formatMessage({ id: 'model.application.code' }),
      dataIndex: 'code',
      search: false,
    },
    
    {
      title: intl.formatMessage({ id: 'model.application.enable'}),
      tooltip: {color:colorPrimary,title:intl.formatMessage({ id: 'model.application.enable.description' })},
      dataIndex: 'enable',
      
      search: false,
      valueEnum: {
        false: {
          text: intl.formatMessage({ id: 'model.application.enable.disable'}),
          status: 'Error',
        },
        true: {
          text: intl.formatMessage({ id: 'model.application.enable.enable'}),
          status: 'Success',
        },
      },
    },
    
    {
      title: intl.formatMessage({ id: 'model.application.home'}),
      dataIndex: 'home',
      search: false,
      render: (dom, entity) => {
        return (
          <a style={{color:colorPrimary}} href={entity.home} target='_blank' rel='noreferrer'>{entity.home}</a>
        );
      },
    },
    {
      title: intl.formatMessage({ id: 'pages.operation'}),
      dataIndex: 'option',
      valueType: 'option',
      hideInTable: !access.systemEditAccess,
      render: (_, record) => {
        const nodes = [
          <a
            key="config"
            onClick={() => {
              navigate(
                {
                  pathname: `${BaseAddress}/update/${record.id}`,
                },
                { replace: true },
              );
            }}
          >
            <FormattedMessage id="pages.operation.edit"  />
          </a>,

        ];

        nodes.push(
          <Dropdown menu={{ items: moreItems(record) }} key='more' >
            <a style={{color:colorPrimary}} onClick={(e) => e.preventDefault()}>
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
  return (
    <PageContainer
    >
      <Paragraph>
        <Text>OpenID Connect</Text>
        <ul>
          <li>
            <Link target='_blank' href={window.location.origin}>Issuer: {window.location.origin}</Link>
          </li>
          <li>
            <Link target='_blank' href={window.location.origin + '/.well-known/openid-configuration'}>Discovery Endpoint</Link>
          </li>
        </ul>
      </Paragraph>
      <ProTable< ApplicationDetail, TableListPagination>
        actionRef={actionRef}
        rowKey="id"
        search={{
           showHiddenNum:true,
           span:{xs: 24, sm: 12, md: 8, lg: 6, xl: 6, xxl: 6}
         }}

        toolBarRender={() => [
          <Access accessible={access.systemAdminAccess === true} key={'create'}>
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
        request={listApplication}
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
    </PageContainer>

  );

};

export default ApplicationTableList;

import { LogoutOutlined, UserOutlined ,DashboardOutlined, HomeOutlined,TableOutlined} from '@ant-design/icons';
import { FormattedMessage, useModel } from '@umijs/max';
import { Spin } from 'antd';
import { createStyles } from 'antd-style';
import React, { useCallback } from 'react';
import { flushSync } from 'react-dom';
import HeaderDropdown from '../HeaderDropdown';
import {    deleteToken } from '@/utils/global';

import access from '@/access';

export type GlobalHeaderRightProps = {
  menu?: boolean;
  children?: React.ReactNode;
};

export const AvatarName = () => {
  const { initialState } = useModel('@@initialState');
  const { currentUser } = initialState || {};
  return <span className="anticon">{currentUser?.username}</span>;
};

const useStyles = createStyles(({ token }) => {
  return {
    action: {
      display: 'flex',
      height: '48px',
      marginLeft: 'auto',
      overflow: 'hidden',
      alignItems: 'center',
      padding: '0 8px',
      cursor: 'pointer',
      borderRadius: token.borderRadius,
      '&:hover': {
        backgroundColor: token.colorBgTextHover,
      },
    },
  };
});

export const AvatarDropdown: React.FC<GlobalHeaderRightProps> = ({ menu, children }) => {
   
  /**
   * 退出登录，并且将当前的 url 保存
   */
  const loginOut = async () => {
    deleteToken();
    window.location.pathname='/index';
  };
  const { styles } = useStyles();

  const { initialState, setInitialState } = useModel('@@initialState');

  const onMenuClick = useCallback(

    (event: any) => {
      const { key } = event;
      if (key === 'logout') {
        flushSync(() => {
          setInitialState((s) => ({ ...s, currentUser: undefined }));
        });
        loginOut();
        return;
      }
      
        if (key === 'console') {
          window.location.href = '/dashboard';
        }else{
          window.location.href = `/personal/${key}`;
        }
      
    },
    [setInitialState],
  );

  const loading = (
    <span className={styles.action}>
      <Spin
        size="small"
        style={{
          marginLeft: 8,
          marginRight: 8,
        }}
      />
    </span>
  );

  if (!initialState) {
    return loading;
  }

  const { currentUser } = initialState;

  if (!currentUser || !currentUser.username) {
    return loading;
  }

  let menuItems = [
    ...(menu
      ? [
        {
          key: 'settings',
          icon: <UserOutlined />,
          label:  <FormattedMessage id='pages.personal.center' />,
        },
        {
          type: 'divider' as const,
        },
      ]
      : []),
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: <FormattedMessage id='pages.personal.logout' />,
    },
  ];
    if (menu&&currentUser?.role!=='none') {
      menuItems.unshift({
        key: 'console',
        icon: <DashboardOutlined />,
        label: <FormattedMessage id='pages.console' />,
      })
    }
   
  return (
    <HeaderDropdown
      menu={{
        selectedKeys: [],
        onClick: onMenuClick,
        items: menuItems,
      }}
    >
      {children}
    </HeaderDropdown>
  );
};

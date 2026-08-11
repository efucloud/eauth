import { AvatarDropdown, AvatarName, Footer, SelectLang, About } from '@/components';
import type { Settings as LayoutSettings, AppListProps } from '@ant-design/pro-components';
import { LinkOutlined } from '@ant-design/icons';

import type { RequestConfig, RunTimeLayoutConfig } from '@umijs/max';
import { FormattedMessage } from '@umijs/max';
import { history, Link } from '@umijs/max';
import defaultSettings from '../config/defaultSettings';
import { errorConfig } from './requestErrorConfig';
import { systemGetUserinfoNoError } from '@/services/oauth.api';
import type { AuthedUserInfo } from '@/services/common.d';
const isDev = process.env.NODE_ENV === 'development';
const loginPath = '/user/login';
const oidcLoginPrefix = '/oauth/callback';
const oidcRegisterPrefix = '/oauth/register';

/**
 * @see  https://umijs.org/zh-CN/plugins/plugin-initial-state
 * */
// export async function getInitialState(): Promise<{
//   settings?: Partial<LayoutSettings>;
//   currentUser?: API.CurrentUser;
//   appList?: AppListProps;
//   loading?: boolean;
//   fetchUserInfo?: () => Promise<API.CurrentUser | undefined>;
// }> {
export async function getInitialState(): Promise<{
  settings?: Partial<LayoutSettings>;
  currentUser?: AuthedUserInfo;
  appList?: AppListProps;
  loading?: boolean;
  getCurrentUserInfo?: () => Promise<AuthedUserInfo | undefined>;
}> {

  const fetchUserInfo = async (): Promise<AuthedUserInfo | undefined> => {
    try {
      const data = await systemGetUserinfoNoError();
      return data as AuthedUserInfo;
    } catch (error) {
      history.push(loginPath);
    }
    return undefined;
  };
  let settings = defaultSettings as Partial<LayoutSettings>;

  const { location } = history;
  if (location.pathname === '/' || location.pathname.startsWith('/oauth/authorize') ||
    location.pathname.indexOf('/forget/password') >= 0 ||
    location.pathname.indexOf('/reset/password') >= 0) {
    return {
      currentUser: {} as AuthedUserInfo,
      getCurrentUserInfo: fetchUserInfo,
      settings: settings,
    };
  } else if (!(location.pathname.startsWith(loginPath) ||
    location.pathname.startsWith(oidcLoginPrefix) ||
    location.pathname.startsWith(oidcRegisterPrefix) ||
    location.pathname.indexOf('/forget/password') >= 0 ||
    location.pathname.indexOf('/reset/password') >= 0
  )) {
    const currentUser = await fetchUserInfo();
    if (currentUser) {
      return {
        getCurrentUserInfo: fetchUserInfo,
        currentUser,
        settings: settings,
      };
    }
    return {
      getCurrentUserInfo: fetchUserInfo,
      currentUser,
      settings: settings,
    };
  }
  return {

    currentUser: {} as AuthedUserInfo,
    getCurrentUserInfo: fetchUserInfo,
    settings: settings,
  };
}

// ProLayout 支持的api https://procomponents.ant.design/components/layout
export const layout: RunTimeLayoutConfig = ({ initialState, setInitialState }) => {
  const waterMark = function () {
    if (initialState?.currentUser?.username !== undefined) {
      return [initialState?.currentUser?.username, initialState?.currentUser?.role];
    } else {
      return undefined;
    }
  };
  const settings = initialState?.settings;

  return {
    siderWidth: 220,
    actionsRender: (): React.ReactNode[] => {
      let nodes = [] as React.ReactNode[];

      nodes.push(<About key="about" />);
      nodes.push(<SelectLang key="SelectLang" />);

      return nodes;
    },

    appList: initialState?.appList,
    avatarProps: {
      src: initialState?.currentUser?.avatar,
      title: <AvatarName />,
      render: (_, avatarChildren) => {
        console.log(initialState?.currentUser);
        if (initialState?.currentUser?.username) {
          return <AvatarDropdown menu={true}>{avatarChildren}</AvatarDropdown>;
        } else {
          return <a href='/user/login'><FormattedMessage id='pages.click.login' /></a>;
        }
      },
    },
    waterMarkProps: {
      // todo 根据路径判断获取用户
      content: waterMark(),
    },
    footerRender: () => <Footer />,
    onPageChange: () => {
      const { location } = history;
      // 如果没有登录，重定向到 login
      if (!initialState?.currentUser && location.pathname !== loginPath) {
        history.push(loginPath);
      }
    },
    bgLayoutImgList: [
      {
        src: '/bg2.webp',
        left: 85,
        bottom: 100,
        height: '303px',
      },
      {
        src: '/bg2.webp',
        bottom: -68,
        right: -45,
        height: '303px',
      },
      {
        src: '/default.webp',
        bottom: 0,
        left: 0,
        width: '331px',
      },
    ],
    links: isDev
      ? [
        <Link key="openapi" to="/umi/plugin/openapi" target="_blank">
          <LinkOutlined />
          <span>OpenAPI 文档</span>
        </Link>,
      ]
      : [],
    menuHeaderRender: undefined,
    // 自定义 403 页面
    // unAccessible: <div>unAccessible</div>,
    // 增加一个 loading 的状态
    childrenRender: (children) => {
      // if (initialState?.loading) return <PageLoading />;
      setInitialState((preInitialState) => ({
        ...preInitialState,
        settings,
      }));
      return (
        <>
          {children}
        </>
      );
    },
    ...settings,

  };
};


/**
 * @name request 配置，可以配置错误处理
 * 它基于 axios 和 ahooks 的 useRequest 提供了一套统一的网络请求和错误处理方案。
 * @doc https://umijs.org/docs/max/request#配置
 */
export const request: RequestConfig = {
  ...errorConfig,
};


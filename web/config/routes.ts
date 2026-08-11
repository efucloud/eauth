/**
 * @name umi 的路由配置
 * @description 只支持 path,component,routes,redirect,wrappers,name,icon 的配置
 * @param path  path 只支持两种占位符配置，第一种是动态参数 :id 的形式，第二种是 * 通配符，通配符只能出现路由字符串的最后。
 * @param component 配置 location 和 path 匹配后用于渲染的 React 组件路径。可以是绝对路径，也可以是相对路径，如果是相对路径，会从 src/pages 开始找起。
 * @param routes 配置子路由，通常在需要为多个路径增加 layout 组件时使用。
 * @param redirect 配置路由跳转
 * @param wrappers 配置路由组件的包装组件，通过包装组件可以为当前的路由组件组合进更多的功能。 比如，可以用于路由级别的权限校验
 * @param name 配置路由的标题，默认读取国际化文件 menu.ts 中 menu.xxxx 的值，如配置 name 为 login，则读取 menu.ts 中 menu.login 的取值作为标题
 * @param icon 配置路由的图标，取值参考 https://ant.design/components/icon-cn， 注意去除风格后缀和大小写，如想要配置图标为 <StepBackwardOutlined /> 则取值应为 stepBackward 或 StepBackward，如想要配置图标为 <UserOutlined /> 则取值应为 user 或者 User
 * @doc https://umijs.org/docs/guides/routes
 */
export default [
  {
    name: 'forget-password',
    menuRender: false,
    hideInMenu: true,
    path: '/forget/password',
    component: './login/forget-pasword',
  },
  {
    name: 'reset-password',
    menuRender: false,
    hideInMenu: true,
    path: '/reset/password/:code',
    component: './login/reset-password',
  },
  {
    name: 'oauth-callback',
    layout: false,
    path: '/oauth/callback/:provider',
    component: './oauth/callback',
  },
  {
    name: 'oauth-register',
    layout: false,
    path: '/oauth/register/:provider',
    component: './oauth/register',
  },
  {
    name: 'authroize',
    layout: false,
    path: '/oauth/authorize',
    component: './oauth/authorize',
  },
  {
    path: '/',
    redirect: '/index',
  },
  {
    name: 'index',
    menuRender: false,
    hideInMenu: true,
    path: '/index',
    component: './index',
  },
   
  {
    icon: 'dashboard',
    name: 'dashboard',
    path: 'dashboard',
    access: 'systemViewAccess',
    component: './dashboard',
  },
  {
    path: '/user',
    layout: false,
    routes: [
      {
        path: '/user/login',
        layout: false,
        name: 'login',
        component: './login',
      },
      {
        path: '/user',
        redirect: '/user/login',
      },
      {
        path: '/user/logout',
        layout: false,
        name: 'logout',
        component: './login/logout',
      },
      {
        component: '404',
        path: '/user/*',
      },
    ],
  },
  {
    path: '/user',
    redirect: '/user/list',
  },
  {
    name: 'user',
    icon: 'team',
    access: 'systemViewAccess',
    path: '/user/list',
    component: './user',
  },
  {
    name: 'user',
    icon: 'team',
    path: '/user/detail/:id',
    component: './user/detail',
    hideInMenu: true,
    access: 'systemViewAccess',
  },
  {
    name: 'user',
    icon: 'team',
    path: '/user/create',
    component: './user/form',
    hideInMenu: true,
    access: 'systemEditAccess',
  },
  {
    name: 'user',
    icon: 'team',
    path: '/user/update/:id',
    component: './user/form',
    hideInMenu: true,
    access: 'systemEditAccess',
  },
  {
    name: 'user',
    path: '/user/detail/:id',
    component: './user/detail',
    hideInMenu: true,
    access: 'systemViewAccess',
  },
   
  {
    name: 'app',
    icon: 'apartment',
    path: '/application/index',
    component: './index',
    hideInMenu: true,
    menuRender: false,
  },
  {
    path: '/app',
    redirect: '/app/list',
  },
  {
    name: 'app',
    icon: 'gold',
    path: '/app/list',
    component: './app',
    access: 'systemViewAccess',
  },
  {
    name: 'app',
    icon: 'gold',
    path: '/app/detail/:id',
    component: './app/detail',
    hideInMenu: true,
    access: 'systemViewAccess',
  },
  {
    name: 'app',
    icon: 'gold',
    path: '/app/create',
    component: './app/form',
    hideInMenu: true,
    access: 'systemAdminAccess',
  },
  {
    name: 'app',
    icon: 'gold',
    path: '/app/update/:id',
    component: './app/form',
    hideInMenu: true,
    access: 'systemEditAccess',
  },
  {
    name: 'provider',
    icon: 'safety',
    access: 'systemViewAccess',
    path: '/provider',
    routes: [
      {
        path: '/provider',
        redirect: '/provider/oidc',
      },
      {
        name: 'oidc',
        path: '/provider/oidc',
        component: './provider/oidc',
        access: 'systemViewAccess',

      },
      {
        name: 'oidc',
        path: '/provider/oidc/detail/:id',
        component: './provider/oidc/detail',
        hideInMenu: true,
        access: 'systemViewAccess',
      },
      {
        name: 'oidc',
        path: '/provider/oidc/create',
        component: './provider/oidc/form',
        hideInMenu: true,
        access: 'systemAdminAccess',
      },
      {
        name: 'oidc',
        path: '/provider/oidc/update/:id',
        component: './provider/oidc/form',
        hideInMenu: true,
        access: 'systemEditAccess',
      },
    ]

  },
    
  {
    name: 'personal',
    icon: 'user',
    path: '/personal/settings',
    access: 'systemPersonalAccess',
    component: './personal/settings',
  },
  {
    component: '404',
    path: '/*',
  },
];

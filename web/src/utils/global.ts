import { AccessTokenResponse } from '@/services/common.d';

const systemOauthUrl = 's';
const tokenKey = 'k';
const bindkey = 'b';

export const getCurrentOrgCode = (): string => {
  if (window.location.pathname.startsWith('/tenant')) {
    return window.location.pathname.split('/')?.[2];
  }
  return '';
};
export const getColorPrimary = (): string => {
  return '#1890ff';
};

// 获取 token
export const getToken = (): AccessTokenResponse => {
  return JSON.parse(localStorage.getItem(tokenKey) || '{}');
};

export const getAllLoginedOrgs = (): string[] => {
  return Object.keys(getToken()) as string[];
};
// 添加 token 到 localStorage
export const addToken = (token: AccessTokenResponse) => {
  token.timestamp = Date.now();
  localStorage.setItem(tokenKey, JSON.stringify(token));
};

// 从 localStorage 删除 token
export const deleteToken = () => {
  localStorage.removeItem(tokenKey);
  return;
};

export const trimSpace = (str: string): string => {
  const reg = /^\s+|\s+$/g;
  return str.replace(reg, '');
};

export const isSystemPersonalCenter = () => {
  return window.location.pathname.startsWith('/personal/');
};

export const getI18nLanguage = () => {
  return localStorage.getItem('umi_locale') || 'zh-CN';
};
export const getSearchParams = () => {
  return new URLSearchParams(window.location.search);
};

export const addBind = (userId: string | undefined) => {
  if (userId) {
    sessionStorage.setItem(bindkey, userId.toString());
  }
};
export const getBind = (): string => {
  return sessionStorage.getItem(bindkey) || '';
};
export const deleteBind = () => {
  sessionStorage.removeItem(bindkey);
};
export const setSystemOauthUrl = (url: string) => {
  sessionStorage.setItem(systemOauthUrl, url);
};
export const getSystemOauthUrl = (): string => {
  return sessionStorage.getItem(systemOauthUrl) || '';
};
export const deleteSystemOauthUrl = () => {
  sessionStorage.removeItem(systemOauthUrl);
};
export const waitTime = (time: number = 100) => {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve(true);
    }, time);
  });
};

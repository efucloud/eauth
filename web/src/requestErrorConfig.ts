import type { RequestConfig } from '@umijs/max';
import { notification } from 'antd';
import { getI18nLanguage, getToken, setSystemOauthUrl, deleteToken } from './utils/global';
import {  ResponseError } from '@/services/common.d';


const authHeaderInterceptor = (url: string, options: RequestConfig) => {
  options.headers = { ...options.headers,   };
  const token = getToken();
  if (token?.timestamp > Date.now()) {
    console.log('token过期');
    // 重定向到登录页
    window.location.href = `/user/login`;
  }
  if (token && token?.id_token) {
    // 添加全局token
    options.headers = { ...options.headers, Authorization: `Bearer ${token.id_token}` };
  }
  return {
    url: `${url}`,
    options: options,
  };
};

/**
 * @name 错误处理
 * pro 自带的错误处理， 可以在这里做自己的改动
 * @doc https://umijs.org/docs/max/request#配置
 */
export const errorConfig: RequestConfig = {
  // 错误处理： umi@3 的错误处理方案。
  errorConfig: {
    // 错误抛出
    errorThrower: (res) => {
      console.log(res);
    },
    // 错误接收及处理
    errorHandler: (error: any) => {
      console.log(error);
      const { response } = error;
      if (response && response.status) {
        const { data, status } = response;
        const errResponse = data as unknown as ResponseError;
        if (status === 303){
          window.location.href = '/';
        }else if (status === 401) {
          // 重定向到登录页面
          deleteToken();
          setSystemOauthUrl(window.location.href);
          window.location.href = '/user/login';

        } else {
          if (errResponse.alert) {
            notification.error({
              message: errResponse.alert,
              description: errResponse.detail,
            });
          } else {
            if (getI18nLanguage() === 'en-US') {
              notification.error({
                message: 'Server connection failed, please check your network connection settings.',
              });
            } else {
              notification.error({
                message: '服务端连接失败，请检查网络连接配置',
              });
            }

          }

        }
      } else {
        const lang = getI18nLanguage();
        if (lang === 'en-US') {
          notification.error({
            message: 'Error',
            description: "Can't connect to the server",
          });
        } else {
          notification.error({
            message: '错误',
            description: '服务器连接失败',
          });
        }

      }
    },
  },

  // 请求拦截器
  requestInterceptors: [authHeaderInterceptor],

  // 响应拦截器
  responseInterceptors: [
    (response) => {
      return response;
    },
  ],
};
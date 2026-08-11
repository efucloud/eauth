import { request } from '@umijs/max';

import { Dashboard } from './dashboard.d';

//系统面板
//系统面板
//请求方法: GET
//请求地址: /api/dashboard
export async function systemDashboard<Dashboard>(  options?: { [key: string]: any }) {
  return  request<Dashboard>(`/api/dashboard`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    ...(options || {}),
  });
}
//30天内应用认证统计
//30天内应用认证统计
//请求方法: GET
//请求地址: /api/dashboard/application-auth-30days
export async function applicationAuth30Days(  options?: { [key: string]: any }) {
  return  request(`/api/dashboard/application-auth-30days`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    ...(options || {}),
  });
}
//系统应用认证TOP10
//系统应用认证TOP10
//请求方法: GET
//请求地址: /api/dashboard/application-auth-top10
export async function applicationAuthTop10(  options?: { [key: string]: any }) {
  return  request(`/api/dashboard/application-auth-top10`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    ...(options || {}),
  });
}

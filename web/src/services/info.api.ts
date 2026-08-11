import { request } from '@umijs/max';


//健康检查
//健康检查
//请求方法: GET
//请求地址: /api/health
export async function health(  options?: { [key: string]: any }) {
  return  request(`/api/health`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    ...(options || {}),
  });
}
//查看应用信息
//查看应用的编译信息
//请求方法: GET
//请求地址: /api/info
export async function getAppInformation(  options?: { [key: string]: any }) {
  return  request(`/api/info`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    ...(options || {}),
  });
}

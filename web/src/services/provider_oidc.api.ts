import { request } from '@umijs/max';

import { ProviderOidcDetail, ProviderOidcUpdate, ProviderOidcCreate, ProviderOidcDetailList, ProviderOidcStatus } from './provider_oidc.d';
import { BatchOperationIds } from './common.d';

//更新OIDC认证提供商信息
//更新OIDC认证提供商信息
//请求方法: PUT
//请求地址: /api/provider-oidc
export async function updateProviderOidc<ProviderOidcDetail>(  data: ProviderOidcUpdate,   options?: { [key: string]: any }) {
  return  request<ProviderOidcDetail>(`/api/provider-oidc`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//创建OIDC认证提供商
//创建OIDC认证提供商信息
//请求方法: POST
//请求地址: /api/provider-oidc
export async function createProviderOidc<ProviderOidcDetail>(  data: ProviderOidcCreate,   options?: { [key: string]: any }) {
  return  request<ProviderOidcDetail>(`/api/provider-oidc`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//删除OIDC认证提供商
//删除OIDC认证提供商信息详情
//请求方法: DELETE
//请求地址: /api/provider-oidc
export async function deleteProviderOidc(  data: BatchOperationIds,   options?: { [key: string]: any }) {
  return  request(`/api/provider-oidc`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//获取OIDC认证提供商列表
//获取OIDC认证提供商信息
//请求方法: GET
//请求地址: /api/provider-oidc
//参数名: pageSize 参数类型: number 参数位置: query 是否必须: false  参数说明: 每页大小
//参数名: order 参数类型: string 参数位置: query 是否必须: false  参数说明: 排序
//参数名: name 参数类型: string 参数位置: query 是否必须: false  参数说明: 名称
//参数名: code 参数类型: string 参数位置: query 是否必须: false  参数说明: 编码
//参数名: current 参数类型: number 参数位置: query 是否必须: false  参数说明: 页码
export async function listProviderOidc<ProviderOidcDetailList>(
  params: {
    current?:number;// 页码
    pageSize?:number;// 每页大小
    order?:string;// 排序
    name?:string;// 名称
    code?:string;// 编码
  },
  options?: { [key: string]: any }) {
  return  request<ProviderOidcDetailList>(`/api/provider-oidc`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    params: params,
    ...(options || {}),
  });
}
//启用禁用
//启用禁用,修改账户状态
//请求方法: POST
//请求地址: /api/provider-oidc/status
export async function changeProviderOidcStatus(  data: ProviderOidcStatus,   options?: { [key: string]: any }) {
  return  request(`/api/provider-oidc/status`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//获取OIDC认证提供商详情
//获取OIDC认证提供商信息详情
//请求方法: GET
//请求地址: /api/provider-oidc/{id}
//参数名: id 参数类型: string 参数位置: path 是否必须: true  参数说明: 记录ID
export async function getProviderOidc<ProviderOidcDetail>(
  params: {
    id: string;// 记录ID
  },
  options?: { [key: string]: any }) {
  const { id, ...rest } = params;
  return  request<ProviderOidcDetail>(`/api/provider-oidc/${id}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...rest },
    ...(options || {}),
  });
}

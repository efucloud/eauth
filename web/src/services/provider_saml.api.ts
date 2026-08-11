import { request } from '@umijs/max';

import { ProviderSamlDetailList, ProviderSamlStatus, ProviderSamlDetail, ProviderSamlUpdate, ProviderSamlCreate } from './provider_saml.d';
import { BatchOperationIds } from './common.d';

//获取SAML认证提供商列表
//获取SAML认证提供商信息
//请求方法: GET
//请求地址: /api/provider-saml
//参数名: pageSize 参数类型: number 参数位置: query 是否必须: false  参数说明: 每页大小
//参数名: order 参数类型: string 参数位置: query 是否必须: false  参数说明: 排序
//参数名: name 参数类型: string 参数位置: query 是否必须: false  参数说明: 名称
//参数名: code 参数类型: string 参数位置: query 是否必须: false  参数说明: 编码
//参数名: current 参数类型: number 参数位置: query 是否必须: false  参数说明: 页码
export async function listProviderSaml<ProviderSamlDetailList>(
  params: {
    current?:number;// 页码
    pageSize?:number;// 每页大小
    order?:string;// 排序
    name?:string;// 名称
    code?:string;// 编码
  },
  options?: { [key: string]: any }) {
  return  request<ProviderSamlDetailList>(`/api/provider-saml`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    params: params,
    ...(options || {}),
  });
}
//启用禁用
//启用禁用,修改SAML提供商状态
//请求方法: POST
//请求地址: /api/provider-saml/status
export async function changeProviderSamlStatus(  data: ProviderSamlStatus,   options?: { [key: string]: any }) {
  return  request(`/api/provider-saml/status`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//获取SAML认证提供商详情
//获取SAML认证提供商信息详情
//请求方法: GET
//请求地址: /api/provider-saml/{id}
//参数名: id 参数类型: string 参数位置: path 是否必须: true  参数说明: 记录ID
export async function getProviderSaml<ProviderSamlDetail>(
  params: {
    id: string;// 记录ID
  },
  options?: { [key: string]: any }) {
  const { id, ...rest } = params;
  return  request<ProviderSamlDetail>(`/api/provider-saml/${id}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...rest },
    ...(options || {}),
  });
}
//更新SAML认证提供商信息
//更新SAML认证提供商信息
//请求方法: PUT
//请求地址: /api/provider-saml
export async function updateProviderSaml<ProviderSamlDetail>(  data: ProviderSamlUpdate,   options?: { [key: string]: any }) {
  return  request<ProviderSamlDetail>(`/api/provider-saml`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//删除SAML认证提供商
//删除SAML认证提供商信息详情
//请求方法: DELETE
//请求地址: /api/provider-saml
export async function deleteProviderSaml(  data: BatchOperationIds,   options?: { [key: string]: any }) {
  return  request(`/api/provider-saml`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//创建SAML认证提供商
//创建SAML认证提供商信息
//请求方法: POST
//请求地址: /api/provider-saml
export async function createProviderSaml<ProviderSamlDetail>(  data: ProviderSamlCreate,   options?: { [key: string]: any }) {
  return  request<ProviderSamlDetail>(`/api/provider-saml`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}

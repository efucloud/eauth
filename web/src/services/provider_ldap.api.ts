import { request } from '@umijs/max';

import { ProviderLdapDetail, ProviderLdapUpdate, ProviderLdapDetailList, ProviderLdapCreate } from './provider_ldap.d';
import { BatchOperationIds } from './common.d';

//获取LDAP认证提供商详情
//获取LDAP认证提供商信息详情
//请求方法: GET
//请求地址: /api/provider-ldap/{id}
//参数名: id 参数类型: string 参数位置: path 是否必须: true  参数说明: 记录ID
export async function getProviderLdap<ProviderLdapDetail>(
  params: {
    id: string;// 记录ID
  },
  options?: { [key: string]: any }) {
  const { id, ...rest } = params;
  return  request<ProviderLdapDetail>(`/api/provider-ldap/${id}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...rest },
    ...(options || {}),
  });
}
//更新LDAP认证提供商信息
//更新LDAP认证提供商信息
//请求方法: PUT
//请求地址: /api/provider-ldap
export async function updateProviderLdap<ProviderLdapDetail>(  data: ProviderLdapUpdate,   options?: { [key: string]: any }) {
  return  request<ProviderLdapDetail>(`/api/provider-ldap`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//删除LDAP认证提供商
//删除LDAP认证提供商信息详情
//请求方法: DELETE
//请求地址: /api/provider-ldap
export async function deleteProviderLdap(  data: BatchOperationIds,   options?: { [key: string]: any }) {
  return  request(`/api/provider-ldap`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//获取LDAP认证提供商列表
//获取LDAP认证提供商信息
//请求方法: GET
//请求地址: /api/provider-ldap
//参数名: current 参数类型: number 参数位置: query 是否必须: false  参数说明: 页码
//参数名: pageSize 参数类型: number 参数位置: query 是否必须: false  参数说明: 每页大小
//参数名: order 参数类型: string 参数位置: query 是否必须: false  参数说明: 排序
//参数名: name 参数类型: string 参数位置: query 是否必须: false  参数说明: 名称
//参数名: code 参数类型: string 参数位置: query 是否必须: false  参数说明: 编码
export async function listProviderLdap<ProviderLdapDetailList>(
  params: {
    current?:number;// 页码
    pageSize?:number;// 每页大小
    order?:string;// 排序
    name?:string;// 名称
    code?:string;// 编码
  },
  options?: { [key: string]: any }) {
  return  request<ProviderLdapDetailList>(`/api/provider-ldap`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    params: params,
    ...(options || {}),
  });
}
//创建LDAP认证提供商
//创建LDAP认证提供商信息
//请求方法: POST
//请求地址: /api/provider-ldap
export async function createProviderLdap<ProviderLdapDetail>(  data: ProviderLdapCreate,   options?: { [key: string]: any }) {
  return  request<ProviderLdapDetail>(`/api/provider-ldap`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}

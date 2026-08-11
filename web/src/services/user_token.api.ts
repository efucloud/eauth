import { request } from '@umijs/max';

import { UserTokenDetailList, UserTokenDetail, UserTokenUpdate } from './user_token.d';
import { UserAuthProfileDetailList } from './user_auth_profile.d';
import { BatchOperationIds } from './common.d';

//获取系统用户令牌列表
//获取系统用户令牌信息
//请求方法: GET
//请求地址: /api/user-token
//参数名: current 参数类型: number 参数位置: query 是否必须: false  参数说明: 页码
//参数名: pageSize 参数类型: number 参数位置: query 是否必须: false  参数说明: 每页大小
//参数名: order 参数类型: string 参数位置: query 是否必须: false  参数说明: 排序
export async function listUserToken<UserTokenDetailList>(
  params: {
    pageSize?:number;// 每页大小
    order?:string;// 排序
    current?:number;// 页码
  },
  options?: { [key: string]: any }) {
  return  request<UserTokenDetailList>(`/api/user-token`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    params: params,
    ...(options || {}),
  });
}
//获取个人的所有访问令牌
//获取个人的所有访问令牌
//请求方法: GET
//请求地址: /api/user-token/personal/authed
//参数名: current 参数类型: number 参数位置: query 是否必须: false  参数说明: 页码
//参数名: pageSize 参数类型: number 参数位置: query 是否必须: false  参数说明: 每页大小
//参数名: order 参数类型: string 参数位置: query 是否必须: false  参数说明: 排序
export async function getUserPersonalUserToken<UserAuthProfileDetailList>(
  params: {
    current?:number;// 页码
    pageSize?:number;// 每页大小
    order?:string;// 排序
  },
  options?: { [key: string]: any }) {
  return  request<UserAuthProfileDetailList>(`/api/user-token/personal/authed`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    params: params,
    ...(options || {}),
  });
}
//删除系统用户令牌
//删除系统用户令牌信息详情
//请求方法: DELETE
//请求地址: /api/user-token
export async function deleteUserToken(  data: BatchOperationIds,   options?: { [key: string]: any }) {
  return  request(`/api/user-token`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//获取系统用户令牌详情
//获取系统用户令牌信息详情
//请求方法: GET
//请求地址: /api/user-token/{id}
//参数名: id 参数类型: string 参数位置: path 是否必须: true  参数说明: 记录ID
export async function getUserToken<UserTokenDetail>(
  params: {
    id: string;// 记录ID
  },
  options?: { [key: string]: any }) {
  const { id, ...rest } = params;
  return  request<UserTokenDetail>(`/api/user-token/${id}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...rest },
    ...(options || {}),
  });
}
//更新系统用户令牌信息
//更新系统用户令牌信息
//请求方法: PUT
//请求地址: /api/user-token
export async function updateUserToken<UserTokenDetail>(  data: UserTokenUpdate,   options?: { [key: string]: any }) {
  return  request<UserTokenDetail>(`/api/user-token`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}

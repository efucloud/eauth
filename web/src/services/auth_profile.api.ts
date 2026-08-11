import { request } from '@umijs/max';

import { UserAuthProfileDetailList, UserAuthProfileDetail, UserAuthProfileStatus } from './user_auth_profile.d';
import { BatchOperationIds } from './common.d';

//获取第三方认证列表
//获取第三方认证信息
//请求方法: GET
//请求地址: /api/auth-profile
//参数名: code 参数类型: string 参数位置: query 是否必须: false  参数说明: 编码
//参数名: current 参数类型: number 参数位置: query 是否必须: false  参数说明: 页码
//参数名: pageSize 参数类型: number 参数位置: query 是否必须: false  参数说明: 每页大小
//参数名: order 参数类型: string 参数位置: query 是否必须: false  参数说明: 排序
//参数名: name 参数类型: string 参数位置: query 是否必须: false  参数说明: 名称
export async function listUserAuthProfile<UserAuthProfileDetailList>(
  params: {
    current?:number;// 页码
    pageSize?:number;// 每页大小
    order?:string;// 排序
    name?:string;// 名称
    code?:string;// 编码
  },
  options?: { [key: string]: any }) {
  return  request<UserAuthProfileDetailList>(`/api/auth-profile`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    params: params,
    ...(options || {}),
  });
}
//获取第三方认证详情
//获取第三方认证信息详情
//请求方法: GET
//请求地址: /api/auth-profile/{id}
//参数名: id 参数类型: string 参数位置: path 是否必须: true  参数说明: 记录ID
export async function getUserAuthProfile<UserAuthProfileDetail>(
  params: {
    id: string;// 记录ID
  },
  options?: { [key: string]: any }) {
  const { id, ...rest } = params;
  return  request<UserAuthProfileDetail>(`/api/auth-profile/${id}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...rest },
    ...(options || {}),
  });
}
//启用禁用
//启用禁用,修改认证方式状态
//请求方法: POST
//请求地址: /api/auth-profile/status
export async function changeUserAuthProfileStatus(  data: UserAuthProfileStatus,   options?: { [key: string]: any }) {
  return  request(`/api/auth-profile/status`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//获取某个人的第三方认证方式
//获取某个人的第三方认证方式
//请求方法: GET
//请求地址: /api/auth-profile/user/{userId}
//参数名: userId 参数类型: number 参数位置: path 是否必须: true  参数说明: 用户ID
export async function getUserAuthProfileByUserId<UserAuthProfileDetailList>(
  params: {
    userId: number;// 用户ID
  },
  options?: { [key: string]: any }) {
  const { userId, ...rest } = params;
  return  request<UserAuthProfileDetailList>(`/api/auth-profile/user/${userId}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...rest },
    ...(options || {}),
  });
}
//获取个人的第三方认证方式
//获取个人的第三方认证方式
//请求方法: GET
//请求地址: /api/auth-profile/personal/authed
export async function getUserPersonalAuthProfile<UserAuthProfileDetailList>(  options?: { [key: string]: any }) {
  return  request<UserAuthProfileDetailList>(`/api/auth-profile/personal/authed`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    ...(options || {}),
  });
}
//删除第三方认证
//删除第三方认证信息详情
//请求方法: DELETE
//请求地址: /api/auth-profile
export async function deleteUserAuthProfile(  data: BatchOperationIds,   options?: { [key: string]: any }) {
  return  request(`/api/auth-profile`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}

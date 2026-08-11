import { request } from '@umijs/max';

import { ShortUser, UserResetPassword } from './user.d';
import { ExistResponse } from './common.d';

//判断手机号码在系统中是否存在
//判断手机号码在系统中是否存在
//请求方法: GET
//请求地址: /api/phone
//参数名: phone 参数类型: string 参数位置: query 是否必须: false  参数说明: 手机号码
export async function systemPhoneExist<ExistResponse>(
  params: {
    phone?:string;// 手机号码
  },
  options?: { [key: string]: any }) {
  return  request<ExistResponse>(`/api/phone`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    params: params,
    ...(options || {}),
  });
}
//发送重置邮件
//发送重置邮件
//请求方法: POST
//请求地址: /api/email
export async function sendResetPwdEmail(  data: ExistResponse,   options?: { [key: string]: any }) {
  return  request(`/api/email`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//判断邮箱在系统中是否存在
//判断邮箱在系统中是否存在
//请求方法: GET
//请求地址: /api/email
//参数名: email 参数类型: string 参数位置: query 是否必须: false  参数说明: 邮箱
export async function systemEmailExist<ExistResponse>(
  params: {
    email?:string;// 邮箱
  },
  options?: { [key: string]: any }) {
  return  request<ExistResponse>(`/api/email`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    params: params,
    ...(options || {}),
  });
}
//获取重置密码用户信息
//获取重置密码用户信息
//请求方法: GET
//请求地址: /api/user/{code}
//参数名: code 参数类型: string 参数位置: path 是否必须: true  参数说明: 重置校验码
export async function systemValidateUserInfo<ShortUser>(
  params: {
    code: string;// 重置校验码
  },
  options?: { [key: string]: any }) {
  const { code, ...rest } = params;
  return  request<ShortUser>(`/api/user/${code}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...rest },
    ...(options || {}),
  });
}
//重置密码
//重置密码
//请求方法: POST
//请求地址: /api/resetpwd
export async function systemUserResetPassword(  data: UserResetPassword,   options?: { [key: string]: any }) {
  return  request(`/api/resetpwd`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}

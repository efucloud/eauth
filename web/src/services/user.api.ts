import { request } from '@umijs/max';

import { UserDetail, UserCreate, UserDetailList, UserStatus, UserResetPassword, UserUpdate, UserRole } from './user.d';
import { BatchOperationIds } from './common.d';

//创建用户
//创建用户信息
//请求方法: POST
//请求地址: /api/user
export async function createUser<UserDetail>(  data: UserCreate,   options?: { [key: string]: any }) {
  return  request<UserDetail>(`/api/user`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//获取用户列表
//获取用户列表
//请求方法: GET
//请求地址: /api/user
//参数名: role 参数类型: string 参数位置: query 是否必须: false  参数说明: 系统角色
//参数名: username 参数类型: string 参数位置: query 是否必须: false  参数说明: 账户名英文
//参数名: email 参数类型: string 参数位置: query 是否必须: false  参数说明: 邮箱
//参数名: id 参数类型: string 参数位置: query 是否必须: false  参数说明: 数据库记录ID
//参数名: order 参数类型: string 参数位置: query 是否必须: false  参数说明: 排序
//参数名: nickname 参数类型: string 参数位置: query 是否必须: false  参数说明: 昵称
//参数名: phone 参数类型: string 参数位置: query 是否必须: false  参数说明: 电话号码
//参数名: jobNumber 参数类型: string 参数位置: query 是否必须: false  参数说明: 工号
//参数名: current 参数类型: number 参数位置: query 是否必须: false  参数说明: 页码
//参数名: pageSize 参数类型: number 参数位置: query 是否必须: false  参数说明: 每页大小
export async function listUser<UserDetailList>(
  params: {
    order?:string;// 排序
    nickname?:string;// 昵称
    phone?:string;// 电话号码
    jobNumber?:string;// 工号
    current?:number;// 页码
    pageSize?:number;// 每页大小
    role?:string;// 系统角色
    username?:string;// 账户名英文
    email?:string;// 邮箱
    id?:string;// 数据库记录ID
  },
  options?: { [key: string]: any }) {
  return  request<UserDetailList>(`/api/user`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    params: params,
    ...(options || {}),
  });
}
//删除用户
//删除用户信息详情
//请求方法: DELETE
//请求地址: /api/user
export async function deleteUser(  data: BatchOperationIds,   options?: { [key: string]: any }) {
  return  request(`/api/user`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//启用禁用
//启用禁用,修改账户状态
//请求方法: POST
//请求地址: /api/user/status
export async function changeUserStatus(  data: UserStatus,   options?: { [key: string]: any }) {
  return  request(`/api/user/status`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//获取用户详情
//获取用户信息详情
//请求方法: GET
//请求地址: /api/user/{id}
//参数名: id 参数类型: string 参数位置: path 是否必须: true  参数说明: 记录ID
export async function getUser<UserDetail>(
  params: {
    id: string;// 记录ID
  },
  options?: { [key: string]: any }) {
  const { id, ...rest } = params;
  return  request<UserDetail>(`/api/user/${id}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...rest },
    ...(options || {}),
  });
}
//重置密码
//管理员给用户重置密码
//请求方法: POST
//请求地址: /api/user/reset/password
export async function resetUserPassword(  data: UserResetPassword,   options?: { [key: string]: any }) {
  return  request(`/api/user/reset/password`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//更新用户信息
//更新用户信息
//请求方法: PUT
//请求地址: /api/user
export async function updateUser<UserDetail>(  data: UserUpdate,   options?: { [key: string]: any }) {
  return  request<UserDetail>(`/api/user`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//系统角色设置
//系统角色设置,admin: 管理员，view: 查看者， edit: 编辑者， none: 无权限，仅为系统成员
//请求方法: POST
//请求地址: /api/user/role
export async function setUserRole<UserDetail>(  data: UserRole,   options?: { [key: string]: any }) {
  return  request<UserDetail>(`/api/user/role`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}

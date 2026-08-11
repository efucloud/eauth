import { request } from '@umijs/max';

import { ApplicationDetail, ApplicationCreate, ApplicationStatus, ApplicationUpdate, ApplicationDetailList } from './application.d';
import { BatchOperationIds } from './common.d';

//创建普通应用
//创建普通应用信息
//请求方法: POST
//请求地址: /api/application
export async function createApplication<ApplicationDetail>(  data: ApplicationCreate,   options?: { [key: string]: any }) {
  return  request<ApplicationDetail>(`/api/application`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//启用禁用
//启用禁用,修改应用状态
//请求方法: POST
//请求地址: /api/application/status
export async function changeApplicationStatus(  data: ApplicationStatus,   options?: { [key: string]: any }) {
  return  request(`/api/application/status`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//获取普通应用详情
//获取普通应用信息详情
//请求方法: GET
//请求地址: /api/application/{id}
//参数名: id 参数类型: string 参数位置: path 是否必须: true  参数说明: 记录ID
export async function getApplication<ApplicationDetail>(
  params: {
    id: string;// 记录ID
  },
  options?: { [key: string]: any }) {
  const { id, ...rest } = params;
  return  request<ApplicationDetail>(`/api/application/${id}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...rest },
    ...(options || {}),
  });
}
//删除普通应用
//删除普通应用信息详情
//请求方法: DELETE
//请求地址: /api/application
export async function deleteApplication(  data: BatchOperationIds,   options?: { [key: string]: any }) {
  return  request(`/api/application`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//更新普通应用信息
//更新普通应用信息
//请求方法: PUT
//请求地址: /api/application
export async function updateApplication<ApplicationDetail>(  data: ApplicationUpdate,   options?: { [key: string]: any }) {
  return  request<ApplicationDetail>(`/api/application`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//获取普通应用列表
//获取普通应用信息
//请求方法: GET
//请求地址: /api/application
//参数名: current 参数类型: number 参数位置: query 是否必须: false  参数说明: 页码
//参数名: pageSize 参数类型: number 参数位置: query 是否必须: false  参数说明: 每页大小
//参数名: order 参数类型: string 参数位置: query 是否必须: false  参数说明: 排序
//参数名: name 参数类型: string 参数位置: query 是否必须: false  参数说明: 名称
//参数名: code 参数类型: string 参数位置: query 是否必须: false  参数说明: 编码
//参数名: id 参数类型: string 参数位置: query 是否必须: false  参数说明: 数据库记录ID
//参数名: search 参数类型: string 参数位置: query 是否必须: false  参数说明: 搜索
export async function listApplication<ApplicationDetailList>(
  params: {
    pageSize?:number;// 每页大小
    order?:string;// 排序
    id?:string;// 数据库记录ID
    current?:number;// 页码
    name?:string;// 名称
    code?:string;// 编码
    search?:string;// 搜索
  },
  options?: { [key: string]: any }) {
  return  request<ApplicationDetailList>(`/api/application`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    params: params,
    ...(options || {}),
  });
}

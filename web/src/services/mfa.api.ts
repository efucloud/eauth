import { request } from '@umijs/max';

import { MultiFactorAuthDetail, PersonalBoundMFA } from './mult_factor_auth.d';
import { BatchOperationIds } from './common.d';

//个人重新绑定MFA
//个人重新绑定MFA， 在个人中心页面操作 
//请求方法: POST
//请求地址: /api/mfa/bound
export async function boundMultiFactorAuthPersonal<MultiFactorAuthDetail>(  data: PersonalBoundMFA,   options?: { [key: string]: any }) {
  return  request<MultiFactorAuthDetail>(`/api/mfa/bound`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//重置个人MFA信息
//重置个人MFA信息，会重新生成新的
//请求方法: POST
//请求地址: /api/mfa/reset
export async function resetMultiFactorAuthPersonal<MultiFactorAuthDetail>(  data: PersonalBoundMFA,   options?: { [key: string]: any }) {
  return  request<MultiFactorAuthDetail>(`/api/mfa/reset`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//获取个人MFA信息列表
//获取个人MFA信息列表
//请求方法: GET
//请求地址: /api/mfa/personal
export async function getMultiFactorAuthPersonal<MultiFactorAuthDetail>(  options?: { [key: string]: any }) {
  return  request<MultiFactorAuthDetail>(`/api/mfa/personal`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    ...(options || {}),
  });
}
//删除MFA数据
//删除MFA数据信息详情
//请求方法: DELETE
//请求地址: /api/mfa
export async function deleteUserMultiFactorAuth(  data: BatchOperationIds,   options?: { [key: string]: any }) {
  return  request(`/api/mfa`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}

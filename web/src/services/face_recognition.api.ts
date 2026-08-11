import { request } from '@umijs/max';

import { ProviderOidcDetail } from './provider_oidc.d';
import { FaceRecognitionCreate, FaceRecognitionStatus, FaceRecognitionDetailList } from './face_recognition.d';
import { AccessTokenResponse, LoginByFaceIdData, BatchOperationIds } from './common.d';

//获取个人人脸识别信息列表
//获取个人人脸识别信息列表
//请求方法: GET
//请求地址: /api/face-recognition/personal
export async function getFaceRecognitionPersonal(  options?: { [key: string]: any }) {
  return  request(`/api/face-recognition/personal`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    ...(options || {}),
  });
}
//创建人脸识别数据
//创建人脸识别数据信息
//请求方法: POST
//请求地址: /api/face-recognition
export async function createUserFaceRecognition<ProviderOidcDetail>(  data: FaceRecognitionCreate,   options?: { [key: string]: any }) {
  return  request<ProviderOidcDetail>(`/api/face-recognition`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//人脸识别登录
//通过人脸识别数据登录
//请求方法: POST
//请求地址: /api/face-recognition/login
export async function loginByFaceIdData<AccessTokenResponse>(  data: LoginByFaceIdData,   options?: { [key: string]: any }) {
  return  request<AccessTokenResponse>(`/api/face-recognition/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//启用禁用
//启用禁用,修改人脸数据状态
//请求方法: POST
//请求地址: /api/face-recognition/status
export async function changeUserFaceRecognitionStatus(  data: FaceRecognitionStatus,   options?: { [key: string]: any }) {
  return  request(`/api/face-recognition/status`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//删除人脸识别数据
//删除人脸识别数据信息详情
//请求方法: DELETE
//请求地址: /api/face-recognition
export async function deleteUserFaceRecognition(  data: BatchOperationIds,   options?: { [key: string]: any }) {
  return  request(`/api/face-recognition`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//获取人脸识别数据列表
//获取人脸识别数据列表
//请求方法: GET
//请求地址: /api/face-recognition
//参数名: current 参数类型: number 参数位置: query 是否必须: false  参数说明: 页码
//参数名: pageSize 参数类型: number 参数位置: query 是否必须: false  参数说明: 每页大小
//参数名: order 参数类型: string 参数位置: query 是否必须: false  参数说明: 排序
//参数名: userId 参数类型: number 参数位置: query 是否必须: false  参数说明: 用户ID
export async function listUserFaceRecognition<FaceRecognitionDetailList>(
  params: {
    current?:number;// 页码
    pageSize?:number;// 每页大小
    order?:string;// 排序
    userId?:number;// 用户ID
  },
  options?: { [key: string]: any }) {
  return  request<FaceRecognitionDetailList>(`/api/face-recognition`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    params: params,
    ...(options || {}),
  });
}
//判断用户是否有人脸识别数据
//根据用户名判断用户是否有人脸识别数据
//请求方法: GET
//请求地址: /api/face-recognition/check/{username}
//参数名: username 参数类型: string 参数位置: path 是否必须: true  参数说明: 登录用户名
export async function checkUserHasFaceIdData(
  params: {
    username: string;// 登录用户名
  },
  options?: { [key: string]: any }) {
  const { username, ...rest } = params;
  return  request(`/api/face-recognition/check/${username}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...rest },
    ...(options || {}),
  });
}

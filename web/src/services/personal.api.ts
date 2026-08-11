import { request } from '@umijs/max';

import { SetPassword } from './user.d';

//更新个人头像
//更新个人头像
//请求方法: POST
//请求地址: /api/personal/self/avatar
export async function updateSelfAvatar(
  params: {
    file: File;// 上传的头像文件
  },
  options?: { [key: string]: any }) {
  const { file, ...rest } = params;
  const formData = new FormData();
  formData.append('file', file);
  return  request(`/api/personal/self/avatar`, {
    method: 'POST',
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    data: formData,
    ...(options || {}),
  });
}
//设置个人密码
//设置个人密码
//请求方法: POST
//请求地址: /api/personal/set/pwd
export async function setPassword(  data: SetPassword,   options?: { [key: string]: any }) {
  return  request(`/api/personal/set/pwd`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}

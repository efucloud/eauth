import { JsonMap } from './common.d';
// UserAuthProfileTempCreate 系统用户认证方式创建
export type UserAuthProfileTempCreate = { 
  //创建时间
  createdAt: string;
  //请求码
  //最大长度: 50
  code: string;
  //认证类型provider中的code,email,phone
  //最大长度: 255
  provider?: string;
  //第三方登录用户的id，邮箱，手机号
  //最大长度: 255
  loginId?: string;
  //第三方登录用户名
  //最大长度: 255
  loginName?: string;
  //昵称
  //最大长度: 255
  nickname?: string;
  //第三方认证的头像
  avatar?: string;
  //用户主页
  //最大长度: 500
  home?: string;
  //第三方认证的全部用户信息
  properties?: JsonMap;
}; 
// UserAuthProfileTempDetail 系统用户认证方式详情
export type UserAuthProfileTempDetail = { 
  //主键
  //最大长度: 50
  id: string;
  //创建时间
  createdAt: string;
  //请求码
  //最大长度: 50
  code: string;
  //认证类型provider中的code,email,phone
  //最大长度: 255
  provider?: string;
  //第三方登录用户的id，邮箱，手机号
  //最大长度: 255
  loginId?: string;
  //第三方登录用户名
  //最大长度: 255
  loginName?: string;
  //昵称
  //最大长度: 255
  nickname?: string;
  //第三方认证的头像
  avatar?: string;
  //用户主页
  //最大长度: 500
  home?: string;
  //第三方认证的全部用户信息
  properties?: JsonMap;
}; 
// UserAuthProfileTempDetailList 系统用户认证方式列表响应
export type UserAuthProfileTempDetailList = { 
  //当前页数据
  data?: UserAuthProfileTempDetail[];
  //数据库满足条件的数据总数
  total: number;
}; 

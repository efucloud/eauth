import { JsonMap } from './common.d';
// UserAuthProfileCreate 系统用户认证方式创建
export type UserAuthProfileCreate = { 
  //用户ID
  //最大长度: 50
  userId: string;
  //认证类型provider中的code,email,phone
  //最大长度: 255
  provider: string;
  //第三方登录用户的id，邮箱，手机号
  //最大长度: 255
  loginId: string;
  //第三方登录用户名
  //最大长度: 255
  loginName?: string;
  //昵称
  //最大长度: 255
  nickname?: string;
  //是否有效
  //默认值: true
  enable: boolean;
  //第三方认证的头像
  avatar?: string;
  //用户主页
  //最大长度: 500
  home?: string;
  //第三方认证的全部用户信息
  properties?: JsonMap;
  //最近一次使用时间
  //最大长度: 255
  latestUsedTime?: string;
}; 
// UserAuthProfileDetail 系统用户认证方式详情
export type UserAuthProfileDetail = { 
  //主键
  //最大长度: 50
  id: string;
  //创建时间
  createdAt: string;
  //创建时间
  updatedAt: string;
  //用户ID
  //最大长度: 50
  userId: string;
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
  //是否有效
  //默认值: true
  enable: boolean;
  //第三方认证的头像
  avatar?: string;
  //用户主页
  //最大长度: 500
  home?: string;
  //第三方认证的全部用户信息
  properties?: JsonMap;
  //最近一次使用时间
  //最大长度: 255
  latestUsedTime?: string;
}; 
// UserAuthProfileDetailList 系统用户认证方式列表响应
export type UserAuthProfileDetailList = { 
  //当前页数据
  data?: UserAuthProfileDetail[];
  //数据库满足条件的数据总数
  total: number;
}; 
// UserAuthProfileStatus 禁用后，用户将不能使用该认证方式登陆系统
export type UserAuthProfileStatus = { 
  //主键
  ids: string[];
  //是否有效
  enable: boolean;
}; 
// UserAuthProfileUpdate 系统用户认证方式修改
export type UserAuthProfileUpdate = { 
  //记录ID
  //最大长度: 50
  id: string;
  //用户ID
  //最大长度: 50
  userId: string;
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
  //最近一次使用时间
  //最大长度: 255
  latestUsedTime?: string;
}; 

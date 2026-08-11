import { ShortUser } from './user.d';
// UserTokenCreate 系统用户Token创建
export type UserTokenCreate = { 
  //用户ID
  //最大长度: 50
  userId: string;
  //客户端ID
  //最大长度: 255
  clientId: string;
  //过期时间(时间戳)
  expired?: number;
  //过期时间
  expiredTime?: string;
  //RefreshToken
  //最大长度: 50
  refreshToken?: string;
  //session key， token MD5
  //最大长度: 50
  sessionKey?: string;
  //Token
  token?: string;
}; 
// UserTokenDetail 系统用户Token详情
export type UserTokenDetail = { 
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
  //用户
  user?: ShortUser;
  //客户端ID
  //最大长度: 255
  clientId: string;
  //过期时间(时间戳)
  expired?: number;
  //过期时间
  expiredTime?: string;
  //RefreshToken
  //最大长度: 50
  refreshToken?: string;
  //session key， token MD5
  //最大长度: 50
  sessionKey?: string;
  //Token
  token?: string;
}; 
// UserTokenDetailList 系统用户Token列表响应
export type UserTokenDetailList = { 
  //当前页数据
  data?: UserTokenDetail[];
  //数据库满足条件的数据总数
  total: number;
}; 
// UserTokenUpdate 系统用户Token修改
export type UserTokenUpdate = { 
  //主键
  //最大长度: 50
  id: string;
  //用户ID
  //最大长度: 50
  userId: string;
  //客户端ID
  //最大长度: 255
  clientId: string;
  //过期时间(时间戳)
  expired?: number;
  //过期时间
  expiredTime?: string;
  //RefreshToken
  //最大长度: 50
  refreshToken?: string;
  //session key， token MD5
  //最大长度: 50
  sessionKey?: string;
  //Token
  token?: string;
}; 

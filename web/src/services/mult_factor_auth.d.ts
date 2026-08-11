// MultiFactorAuthCreate 用户MFA信息创建
export type MultiFactorAuthCreate = { 
  //所属用户
  //最大长度: 50
  userId: string;
  //密钥
  //最大长度: 50
  secret: string;
  //二维码
  image: string;
  //状态：是否已绑定
  //默认值: unbound
  //最大长度: 50
  status: string;
}; 
// MultiFactorAuthDetail 用户MFA信息详情
export type MultiFactorAuthDetail = { 
  //主键
  //最大长度: 50
  id: string;
  //创建时间
  createdAt: string;
  //创建时间
  updatedAt: string;
  //所属用户
  //最大长度: 50
  userId: string;
  //密钥
  //最大长度: 50
  secret: string;
  //二维码
  image: string;
  //状态：是否已绑定
  //默认值: unbound
  //最大长度: 50
  status: string;
}; 
// MultiFactorAuthDetailList 用户MFA信息
export type MultiFactorAuthDetailList = { 
  //当前页数据
  data?: MultiFactorAuthDetail[];
  //数据库满足条件的数据总数
  total: number;
}; 
// MultiFactorAuthStatus 禁用后，用户将不能使用该认证方式登陆系统
export type MultiFactorAuthStatus = { 
  //主键
  id?: string;
  //状态：是否已绑定
  //默认值: unbound
  //最大长度: 50
  status: string;
}; 
// PersonalBoundMFA 个人绑定MFA
export type PersonalBoundMFA = { 
  code?: string;
  client?: string;
}; 

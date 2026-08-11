// ProviderOidcCreate OIDC提供商创建
export type ProviderOidcCreate = { 
  //提供商名称
  //最大长度: 255
  name: string;
  //图标
  //最大长度: 255
  category: string;
  //client ID，在提供商创建应用时生成
  //最大长度: 255
  clientId: string;
  //ClientSecret，在提供商创建应用时生成
  //最大长度: 255
  clientSecret: string;
  //颁发者地址
  //最大长度: 500
  issuer?: string;
  //作用域
  scopes?: string[];
  //是否有效
  //默认值: true
  enable: boolean;
  //授权地址
  //最大长度: 255
  authorizationEndpoint: string;
  //令牌获取地址
  //最大长度: 255
  tokenEndpoint: string;
  //用户信息获取地址
  //最大长度: 255
  userinfoEndpoint?: string;
}; 
// ProviderOidcDetail OIDC提供商详情
export type ProviderOidcDetail = { 
  //主键
  //最大长度: 50
  id: string;
  //创建时间
  createdAt: string;
  //创建时间
  updatedAt: string;
  //提供商名称
  //最大长度: 255
  name: string;
  //提供商类型，根据该字段显示不同的图标
  //最大长度: 255
  category: string;
  //client ID，在提供商创建应用时生成
  //最大长度: 255
  clientId: string;
  //ClientSecret，在提供商创建应用时生成
  //最大长度: 255
  clientSecret: string;
  //颁发者地址
  //最大长度: 500
  issuer?: string;
  //作用域
  scopes?: string[];
  //是否有效
  //默认值: true
  enable: boolean;
  //授权地址
  //最大长度: 255
  authorizationEndpoint: string;
  //令牌获取地址
  //最大长度: 255
  tokenEndpoint: string;
  //用户信息获取地址
  //最大长度: 255
  userinfoEndpoint?: string;
}; 
// ProviderOidcDetailList OIDC提供商列表响应
export type ProviderOidcDetailList = { 
  //当前页数据
  data?: ProviderOidcDetail[];
  //数据库满足条件的数据总数
  total: number;
}; 
// ProviderOidcStatus 认证提供商状态
// 状态为disable时将不在用户前端显示
export type ProviderOidcStatus = { 
  //主键
  ids: number[];
  //是否有效
  enable: boolean;
}; 
// ProviderOidcUpdate OIDC提供商修改
export type ProviderOidcUpdate = { 
  //主键
  //最大长度: 50
  id: string;
  //提供商名称
  //最大长度: 255
  name: string;
  //图标
  //最大长度: 255
  category: string;
  //client ID，在提供商创建应用时生成
  //最大长度: 255
  clientId: string;
  //ClientSecret，在提供商创建应用时生成
  //最大长度: 255
  clientSecret: string;
  //颁发者地址
  //最大长度: 500
  issuer?: string;
  //作用域
  scopes?: string[];
  //是否有效
  //默认值: true
  enable: boolean;
  //授权地址
  //最大长度: 255
  authorizationEndpoint: string;
  //令牌获取地址
  //最大长度: 255
  tokenEndpoint: string;
  //用户信息获取地址
  //最大长度: 255
  userinfoEndpoint?: string;
}; 

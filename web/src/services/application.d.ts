// ApplicationCreate 普通应用创建
export type ApplicationCreate = { 
  //应用名称
  //最大长度: 255
  name: string;
  //应用编码
  //最大长度: 255
  code: string;
  //应用描述
  //最大长度: 255
  description?: string;
  //应用主页
  //最大长度: 1000
  home?: string;
  //是否有效
  //默认值: true
  enable: boolean;
  //Logo
  //最大长度: 1000
  logo?: string;
  //OIDC client ID，在提供商创建应用时生成
  //最大长度: 255
  clientId: string;
  //OIDC 客户端密钥
  clientSecret: string;
  //OIDC 回调地址
  //最大长度: 1000
  redirectUri: string;
  //OIDC 重定向地址匹配类型
  //默认值: equal
  //最大长度: 50
  redirectUriMatchType: string;
}; 
// ApplicationDetail 普通应用详情
export type ApplicationDetail = { 
  //主键
  //最大长度: 50
  id: string;
  //创建时间
  createdAt: string;
  //创建时间
  updatedAt: string;
  //应用名称
  //最大长度: 255
  name: string;
  //应用编码
  //最大长度: 255
  code: string;
  //应用描述
  //最大长度: 255
  description?: string;
  //应用主页
  //最大长度: 1000
  home?: string;
  //是否有效
  //默认值: true
  enable: boolean;
  //Logo
  //最大长度: 1000
  logo?: string;
  //OIDC client ID，在提供商创建应用时生成
  //最大长度: 255
  clientId: string;
  //OIDC 客户端密钥
  clientSecret: string;
  //OIDC 回调地址
  //最大长度: 1000
  redirectUri: string;
  //OIDC 重定向地址匹配类型
  //默认值: equal
  //最大长度: 50
  redirectUriMatchType: string;
}; 
// ApplicationDetailList 普通应用列表响应
export type ApplicationDetailList = { 
  //当前页数据
  data?: ShortApplication[];
  //数据库满足条件的数据总数
  total: number;
}; 
// ApplicationStatus 普通应用状态
// 状态为disable时将不在用户前端显示，同时普通应用中的应用将不能认证
export type ApplicationStatus = { 
  //主键
  ids: number[];
  //是否有效
  enable: boolean;
}; 
// ApplicationUpdate 普通应用修改
export type ApplicationUpdate = { 
  //主键
  //最大长度: 50
  id: string;
  //应用名称
  //最大长度: 255
  name: string;
  //应用编码
  //最大长度: 255
  code: string;
  //应用描述
  //最大长度: 255
  description?: string;
  //应用主页
  //最大长度: 1000
  home?: string;
  //是否有效
  //默认值: true
  enable: boolean;
  //Logo
  //最大长度: 1000
  logo?: string;
  //OIDC 回调地址
  //最大长度: 1000
  redirectUri: string;
  //OIDC 重定向地址匹配类型
  //默认值: equal
  //最大长度: 50
  redirectUriMatchType: string;
}; 
// ShortApplication 简单应用信息
export type ShortApplication = { 
  //主键
  //最大长度: 50
  id: string;
  //创建时间
  createdAt: string;
  //创建时间
  updatedAt: string;
  //应用名称
  //最大长度: 255
  name: string;
  //应用编码
  //最大长度: 255
  code: string;
  //应用描述
  //最大长度: 255
  description?: string;
  //应用主页
  //最大长度: 1000
  home?: string;
  //是否有效
  //默认值: true
  enable: boolean;
  //Logo
  //最大长度: 1000
  logo?: string;
  //OIDC client ID，在提供商创建应用时生成
  //最大长度: 255
  clientId: string;
  //OIDC 客户端密钥
  clientSecret: string;
  //OIDC 回调地址
  //最大长度: 1000
  redirectUri: string;
  //OIDC 重定向地址匹配类型
  //默认值: equal
  //最大长度: 50
  redirectUriMatchType: string;
}; 

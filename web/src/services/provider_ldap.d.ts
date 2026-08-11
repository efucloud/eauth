import { LdapConfig } from './common.d';
// ProviderLdapCreate LDAP提供商创建
export type ProviderLdapCreate = { 
  //名称
  //最大长度: 255
  name: string;
  //编码
  //最大长度: 255
  code: string;
  //提供商类型
  //最大长度: 255
  category: string;
  //是否有效
  //默认值: true
  enable: boolean;
  //认证提供商
  //最大长度: 255
  provider: string;
  //LDAP配置
  ldapConfig?: LdapConfig;
}; 
// ProviderLdapDetail LDAP提供商详情
export type ProviderLdapDetail = { 
  //主键
  //最大长度: 50
  id: string;
  //创建时间
  createdAt: string;
  //创建时间
  updatedAt: string;
  //名称
  //最大长度: 255
  name: string;
  //图标
  //最大长度: 255
  icon: string;
  //是否有效
  //默认值: true
  enable: boolean;
  //认证提供商
  //最大长度: 255
  provider: string;
  //LDAP配置
  ldapConfig?: LdapConfig;
}; 
// ProviderLdapDetailList LDAP提供商列表响应
export type ProviderLdapDetailList = { 
  //当前页数据
  data?: ProviderLdapDetail[];
  //数据库满足条件的数据总数
  total: number;
}; 
// ProviderLdapUpdate LDAP提供商修改
export type ProviderLdapUpdate = { 
  //主键
  //最大长度: 50
  id: string;
  //名称
  //最大长度: 255
  name: string;
  //编码
  //最大长度: 255
  code: string;
  //最大长度: 255
  category: string;
  //是否有效
  //默认值: true
  enable: boolean;
  //认证提供商
  //最大长度: 255
  provider: string;
  //LDAP配置
  ldapConfig?: LdapConfig;
}; 

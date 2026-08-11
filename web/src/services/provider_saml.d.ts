// ProviderSamlCreate SAML提供商创建
export type ProviderSamlCreate = { 
  //提供商名称
  //最大长度: 255
  name: string;
  //提供商编码
  //最大长度: 255
  category: string;
  //是否有效
  //默认值: true
  enable: boolean;
  //IdP EntityID
  //最大长度: 500
  entityId: string;
  //SAML SSO地址
  //最大长度: 1000
  ssoUrl: string;
  //回调地址(AssertionConsumerService URL)
  //最大长度: 1000
  acsUrl: string;
  //SAML证书
  certificate: string;
  //元数据地址
  //最大长度: 1000
  metadataUrl?: string;
  //元数据内容
  metadata?: string;
  //登录ID字段映射
  //最大长度: 255
  loginIdAttr?: string;
  //登录名字段映射
  //最大长度: 255
  loginNameAttr?: string;
  //邮箱字段映射
  //最大长度: 255
  emailAttr?: string;
  //手机号字段映射
  //最大长度: 255
  phoneAttr?: string;
  //昵称字段映射
  //最大长度: 255
  nicknameAttr?: string;
  //头像字段映射
  //最大长度: 255
  avatarAttr?: string;
}; 
// ProviderSamlDetail SAML提供商详情
export type ProviderSamlDetail = { 
  //主键
  //最大长度: 50
  id: string;
  //创建时间
  createdAt: string;
  //更新时间
  updatedAt: string;
  //提供商名称
  //最大长度: 255
  name: string;
  //提供商编码
  //最大长度: 255
  category: string;
  //是否有效
  //默认值: true
  enable: boolean;
  //IdP EntityID
  //最大长度: 500
  entityId: string;
  //SAML SSO地址
  //最大长度: 1000
  ssoUrl: string;
  //回调地址(AssertionConsumerService URL)
  //最大长度: 1000
  acsUrl: string;
  //SAML证书
  certificate: string;
  //元数据地址
  //最大长度: 1000
  metadataUrl?: string;
  //元数据内容
  metadata?: string;
  //登录ID字段映射
  //最大长度: 255
  loginIdAttr?: string;
  //登录名字段映射
  //最大长度: 255
  loginNameAttr?: string;
  //邮箱字段映射
  //最大长度: 255
  emailAttr?: string;
  //手机号字段映射
  //最大长度: 255
  phoneAttr?: string;
  //昵称字段映射
  //最大长度: 255
  nicknameAttr?: string;
  //头像字段映射
  //最大长度: 255
  avatarAttr?: string;
}; 
// ProviderSamlDetailList SAML提供商列表响应
export type ProviderSamlDetailList = { 
  //当前页数据
  data?: ProviderSamlDetail[];
  //数据库满足条件的数据总数
  total: number;
}; 
// ProviderSamlStatus 认证提供商状态
export type ProviderSamlStatus = { 
  //主键
  ids: number[];
  //是否有效
  enable: boolean;
}; 
// ProviderSamlUpdate SAML提供商修改
export type ProviderSamlUpdate = { 
  //主键
  //最大长度: 50
  id: string;
  //提供商名称
  //最大长度: 255
  name: string;
  //提供商编码
  //最大长度: 255
  category: string;
  //是否有效
  //默认值: true
  enable: boolean;
  //IdP EntityID
  //最大长度: 500
  entityId: string;
  //SAML SSO地址
  //最大长度: 1000
  ssoUrl: string;
  //回调地址(AssertionConsumerService URL)
  //最大长度: 1000
  acsUrl: string;
  //SAML证书
  certificate: string;
  //元数据地址
  //最大长度: 1000
  metadataUrl?: string;
  //元数据内容
  metadata?: string;
  //登录ID字段映射
  //最大长度: 255
  loginIdAttr?: string;
  //登录名字段映射
  //最大长度: 255
  loginNameAttr?: string;
  //邮箱字段映射
  //最大长度: 255
  emailAttr?: string;
  //手机号字段映射
  //最大长度: 255
  phoneAttr?: string;
  //昵称字段映射
  //最大长度: 255
  nicknameAttr?: string;
  //头像字段映射
  //最大长度: 255
  avatarAttr?: string;
}; 
export type applyProviderSamlDefaults = { 
}; 
export type normalizeSamlCertificate = { 
}; 
export type parseSamlEntityDescriptor = { 
}; 
export type parseSamlMetadata = { 
}; 
//是否有效
export type samlMetadataEntitiesDescriptor = { 
  EntityDescriptors?: samlMetadataEntityDescriptor[];
}; 
//是否有效
export type samlMetadataEntityDescriptor = { 
  EntityID?: string;
  IDPSSODescriptors?: samlMetadataIDPSSODescriptor[];
}; 
//是否有效
export type samlMetadataIDPSSODescriptor = { 
  SingleSignOnServices?: samlMetadataService[];
  KeyDescriptors?: samlMetadataKeyDescriptor[];
}; 
//是否有效
export type samlMetadataKeyDescriptor = { 
  Use?: string;
  KeyInfo?: samlMetadataKeyInfo;
}; 
//是否有效
export type samlMetadataKeyInfo = { 
  X509Data?: samlMetadataX509Data;
}; 
//是否有效
export type samlMetadataService = { 
  Binding?: string;
  Location?: string;
}; 
//是否有效
export type samlMetadataX509Data = { 
  Certificates?: string[];
}; 
export type selectSamlCertificate = { 
}; 
export type selectSamlSSOServiceURL = { 
}; 

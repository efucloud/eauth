// AccessTokenResponse 登录获取token
export type AccessTokenResponse = { 
  access_token: string;
  token_type : string;
  expires_in: number;
  refresh_token: string;
  id_token: string;
  timestamp: number;
  need: boolean;
  id: string;
  code?: string;
  username?: string;
  nickname?: string;
  phone?: string;
  email?: string;
  mfa: boolean;
  secret?: string;
  image?: string;
  //登录方式 未来动态认证时使用
  method?: string;
}; 
//录入名称
export type ApplicationRole = { 
  name?: string;
  role?: string;
  zhCn?: string;
  enUs?: string;
}; 
//录入名称
export type ApplicationRoles = { 
}; 
// ArrayFloat64 浮点数组
export type ArrayFloat64 = { 
}; 
// AuthedUserInfo 用户信息
export type AuthedUserInfo = { 
  //主键
  id: string;
  //用户名
  username: string;
  //昵称，如中文名
  nickname: string;
  //工号
  jobNumber: string;
  //系统角色
  role: string;
  //是否有效
  enable: boolean;
  //邮箱
  email: string;
  //手机号码
  phone: string;
  //默认语言
  language: string;
  //头像
  avatar: string;
  //有设置密码
  hasPassword: boolean;
  //密码强度
  passwordStrength?: string;
}; 
// BatchOperationIds 需要删除的列表,根据数据库id
export type BatchOperationIds = { 
  //需要删key 可以为数据库的id
  ids: string[];
}; 
// BatchOperationKeys 需要删除的列表，根据表唯一字符型字段
export type BatchOperationKeys = { 
  //需要删key 可以为数据库的id
  keys?: string[];
}; 
// CaptchaCheckData 行为验证码校验请求数据
export type CaptchaCheckData = { 
  code?: string;
  data?: string;
}; 
// CaptchaCheckResponse 行为验证码校验结果
export type CaptchaCheckResponse = { 
  success: boolean;
}; 
// CaptchaResponse 行为验证码生成数据
export type CaptchaResponse = { 
  code?: string;
  image?: string;
  prompt?: string;
  x?: number;
  y?: number;
  size?: number;
  width?: number;
  height?: number;
}; 
export type CheckFaceIdData = { 
}; 
export type ExistResponse = { 
  name?: string;
  exist: boolean;
}; 
// ExternalApp 外部应用
export type ExternalApp = { 
  //应用名称
  title: string;
  //应用地址
  url: string;
  //打开方式
  target: string;
  //应用描述
  desc: string;
  //子应用
  children?: ExternalApp[];
}; 
// JsonMap json对象
export type JsonMap = { 
}; 
// LdapConfig holds configuration options for LDAP logins.
export type LdapConfig = { 
  url?: string;
  baseDn?: string;
  uid?: string;
  bindUser?: string;
  bindPassword?: string;
  filter?: string;
  attributes?: string[];
}; 
// LoginByFaceIdData 人脸识别数据登录
export type LoginByFaceIdData = { 
  //自动登录,12小时，1周，15天,1个月，半年
  //-|12h|1w|15d|1m|0.5y
  rememberMe?: string;
  //用户名,邮箱、手机号码、工号
  username: string;
  //人脸识别数据
  faceIdData: Float64Array;
  //账户来源
  source?: string;
}; 
//用户主页
export type LoginByLDAP = { 
  //自动登录,12小时，1周，15天,1个月，半年
  rememberMe?: string;
  //绑定的系统用户
  bindId?: string;
  //用户名
  username: string;
  //密码
  password: string;
}; 
// LoginByOIDC OIDC登录
export type LoginByOIDC = { 
  rememberMe?: string;
  code?: string;
  bindId?: string;
  provider: string;
  redirectUri?: string;
}; 
// LoginBySAML SAML登录
export type LoginBySAML = { 
  rememberMe?: string;
  provider: string;
  samlResponse?: string;
  SAMLResponse?: string;
  relayState?: string;
  RelayState?: string;
}; 
// LoginByUsername 用户名密码登录
export type LoginByUsername = { 
  //自动登录,12小时，1周，15天,1个月，半年
  //-|12h|1w|15d|1m|0.5y
  rememberMe?: string;
  //用户名,邮箱、手机号码、工号
  username: string;
  //密码
  password: string;
  //账户来源
  source?: string;
}; 
// Value 实现 driver.Valuer 接口，Value 返回 json value
export type MfaCode = { 
  //用户ID
  userId: string;
  //验证码
  code?: string;
}; 
// ThirdAuthMethod 支持的第三方登录方式
export type OIDC = { 
  //名称
  name: string;
  //认证地址
  address: string;
  //提供商类型
  category: string;
}; 
// OidcCodeRequest 获取认证用户的OAuth的code
export type OidcCodeRequest = { 
  //客户端ID
  clientId?: string;
  //重定向地址
  redirectUri?: string;
  //状态码
  state?: string;
  //响应类型
  responseType?: string;
  //PKCE挑战值(camel风格)
  codeChallenge?: string;
  //PKCE挑战方法(camel风格)
  codeChallengeMethod?: string;
  //PKCE挑战值(标准风格)
  code_challenge?: string;
  //PKCE挑战方法(标准风格)
  code_challenge_method?: string;
  //OIDC nonce
  nonce?: string;
}; 
//OIDC nonce
export type OidcCodeResponse = { 
  code?: string;
  //状态码
  state?: string;
  //重定向地址
  redirectUri?: string;
}; 
// OidcRequestToken oidc认证获取token请求结构
export type OidcRequestToken = { 
  //客户端ID
  client_id: string;
  //客户端密钥
  client_secret: string;
  //类型
  grant_type: string;
  //请求码
  code: string;
  //重定向地址
  redirect_uri?: string;
  //PKCE验证码
  code_verifier?: string;
  //PKCE验证码(camel风格兼容)
  codeVerifier?: string;
}; 
//重定向地址
export type OpenIDConfiguration = { 
  issuer?: string;
  authorization_endpoint?: string;
  token_endpoint?: string;
  userinfo_endpoint?: string;
  jwks_uri?: string;
  introspection_endpoint?: string;
  revocation_endpoint?: string;
  response_type s_supported?: string[];
  response_modes_supported?: string[];
  grant_type s_supported?: string[];
  subject_type s_supported?: string[];
  id_token_signing_alg_values_supported?: string[];
  scopes_supported?: string[];
  claims_supported?: string[];
  code_challenge_methods_supported?: string[];
  request_parameter_supported: boolean;
  request_object_signing_alg_values_supported?: string[];
}; 
// Value 实现 driver.Valuer 接口，Value 返回 json value
export type ProviderExtend = { 
  workWeXinAgentId?: string;
}; 
// RegisterByOIDC OIDC认证时注册新用户时补全的信息
export type RegisterByOIDC = { 
  //请求码
  code: string;
  //默认用户名
  username: string;
  // 密码
  password: string;
  //昵称
  nickname: string;
  //手机号码
  phone: string;
  //邮箱
  email: string;
}; 
// ResponseError 错误响应
export type ResponseError = { 
  //错误英文编码
  message?: string;
  //错误详情信息
  detail?: string;
  //支持I18N的提示信息
  alert?: string;
  //当前请求地址
  requestUri?: string;
}; 
//提供商类型
export type SAML = { 
  //名称
  name: string;
  //认证地址
  address: string;
  //提供商类型
  category: string;
}; 
//验证码
export type SamlApplication = { 
  idpMetadata?: string;
  spMetadata?: string;
  privateKey?: string;
  certificate?: string;
  metadataUrl?: string;
  ssoUrl?: string;
  logoutUrl?: string;
  metadata?: string;
}; 
// SamlAuth saml认证
export type SamlAuth = { 
  SamlRequest?: string;
  RelayState?: string;
}; 
// TableListPagination 分页信息
export type TableListPagination = { 
  //总数
  total?: number;
  //每页数量
  pageSize?: number;
  //当前页
  current?: number;
}; 
// ThirdAuthMethod 支持的第三方登录方式
export type ThirdAuthMethod = { 
  oidcs: OIDC[];
  samls?: SAML[];
  mfa: boolean;
  faceRecognition: boolean;
}; 
//账户来源
export type ThirdAuthProfile = { 
  //所属用户
  userId?: string;
  //认证提供商
  provider?: string;
  //第三方认证信息中的ID
  loginId?: string;
  //第三方认证信息中的用户名
  username?: string;
  //第三方认证信息中的邮箱
  email?: string;
  //手机号码
  phone?: string;
  //第三方认证信息中的别名
  nickname?: string;
  //第三方认证信息中的头像
  avatar?: string;
  //第三方认证返回的所有用户信息
  properties?: JsonMap;
  //用户主页
  home?: string;
}; 
//登录方式 未来动态认证时使用
export type UserClaims = { 
  //系统用户ID
  id?: string;
  // 用户名 组织内唯一必须由DNS-1123标签格式的单元组成
  username?: string;
  // 昵称，如中文名
  nickname?: string;
  // 系统角色
  role?: string;
  nonce?: string;
  email?: string;
  phone?: string;
}; 
// UserFaceRecognition 人脸信息
export type UserFaceRecognition = { 
  //是否有效
  //默认值: true
  enable: boolean;
  //录入名称
  //最大长度: 50
  name: string;
}; 
// UserThirdMethod 用户认证方式
export type UserOIdcInfo = { 
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
  //最近一次使用时间
  //最大长度: 255
  latestUsedTime?: string;
}; 
// UserThirdMethod 用户认证方式
export type UserThirdMethod = { 
  oidcs: UserOIdcInfo[];
  webAuthns?: UserWebAuthn[];
  faceRecognitions?: UserFaceRecognition[];
  mfa: string;
}; 
//最近一次使用时间
export type UserWebAuthn = { 
  //是否有效
  //默认值: true
  enable: boolean;
  //密钥名称
  //最大长度: 50
  name: string;
}; 

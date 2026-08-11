import { request } from '@umijs/max';

import { UserUpdate } from './user.d';
import { AuthedUserInfo, AccessTokenResponse, LoginBySAML, LoginByOIDC, OidcCodeResponse, OidcCodeRequest, MfaCode, LoginByUsername, RegisterByOIDC, LoginByLDAP } from './common.d';

//获取用户信息
//获取用户信息
//请求方法: GET
//请求地址: /api/userinfo
export async function systemGetUserinfo<AuthedUserInfo>(  options?: { [key: string]: any }) {
  return  request<AuthedUserInfo>(`/api/userinfo`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    ...(options || {}),
  });
}
//获取Token
//获取Token，支持四种token获取方式
//请求方法: POST
//请求地址: /api/token
//参数名: client_id 参数类型: string 参数位置: query 是否必须: false  参数说明: client ID
//参数名: client_secret 参数类型: string 参数位置: query 是否必须: false  参数说明: client TotpSecret，使用场景:grant_type 为password;client_credentials
//参数名: grant_type 参数类型: string 参数位置: query 是否必须: false  参数说明: Grant Type
//参数名: code 参数类型: string 参数位置: query 是否必须: false  参数说明: StateCode
//参数名: redirect_uri 参数类型: string 参数位置: query 是否必须: false  参数说明: Redirect URI
//参数名: code_verifier 参数类型: string 参数位置: query 是否必须: false  参数说明: PKCE code verifier
//参数名: username 参数类型: string 参数位置: query 是否必须: false  参数说明: Username
//参数名: password 参数类型: string 参数位置: query 是否必须: false  参数说明: Password
export async function systemFetchTokenWithPost<AccessTokenResponse>(
  params: {
    username?:string;// Username
    password?:string;// Password
    client_id?:string;// client ID
    redirect_uri?:string;// Redirect URI
    client_secret?:string;// client TotpSecret，使用场景:grant_type 为password;client_credentials
    grant_type?:string;// Grant Type
    code?:string;// StateCode
    code_verifier?:string;// PKCE code verifier
  },
  options?: { [key: string]: any }) {
  return  request<AccessTokenResponse>(`/api/token`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json;application/x-www-form-urlencoded',
    },
    params: params,
    ...(options || {}),
  });
}
//openid jwks
//openid jwks
//请求方法: GET
//请求地址: /api/discovery/jwks
export async function systemPublicKeys(  options?: { [key: string]: any }) {
  return  request(`/api/discovery/jwks`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    ...(options || {}),
  });
}
//列出支持的第三方登录方式
//列出支持的第三方登录方式包括OIDC，LDAP，未来包括cas等
//请求方法: GET
//请求地址: /api/third/auth/methods
export async function thirdAuthMethods(  options?: { [key: string]: any }) {
  return  request(`/api/third/auth/methods`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    ...(options || {}),
  });
}
//SAML方式登录
//SAML回调后前端给到后端的SAMLResponse接口，用于解析SAML断言并登录系统
//请求方法: POST
//请求地址: /api/login/saml
export async function systemLoginBySaml<AccessTokenResponse>(  data: LoginBySAML,   options?: { [key: string]: any }) {
  return  request<AccessTokenResponse>(`/api/login/saml`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//更新个人信息
//用户更新个人信息
//请求方法: PUT
//请求地址: /api/self/info
export async function updateSelfInfo<AuthedUserInfo>(  data: UserUpdate,   options?: { [key: string]: any }) {
  return  request<AuthedUserInfo>(`/api/self/info`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//获取用户信息
//获取用户信息，在没有认证的情况下返回空内容
//请求方法: GET
//请求地址: /api/userinfo-no-error
export async function systemGetUserinfoNoError<AuthedUserInfo>(  options?: { [key: string]: any }) {
  return  request<AuthedUserInfo>(`/api/userinfo-no-error`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    ...(options || {}),
  });
}
//OIDC方式登录
//OIDC回调后前端给到后端的Code接口，用于换取第三方的token并获取用户信息，若用户在系统不存在，则根据组织是否允许自动注册来决定是否自动创建用户信息，若第一次是通过第三方登录，需要先设置密码，若组织设置了MFA则需要再次输入验证码，若用户没有绑定过验证器，则返回验证器的二维码和密钥
//请求方法: POST
//请求地址: /api/login/oidc
export async function systemLoginByOidc<AccessTokenResponse>(  data: LoginByOIDC,   options?: { [key: string]: any }) {
  return  request<AccessTokenResponse>(`/api/login/oidc`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//Token撤销
//校验Token是否合法并查询信息，比如，access_token是否还有效，谁颁发的，颁发给谁的，scope又哪些等等的元数据信息
//请求方法: POST
//请求地址: /api/revoke
export async function systemRevoke(
  params: {
    token: string;// 可以是access_token或者refresh_token
    token_type_hint: string;// 表示token的类型
    client_id: string;// client ID
    client_secret: string;// client TotpSecret，使用场景:grant_type 为password;client_credentials
  },
  options?: { [key: string]: any }) {
  const { token, token_type_hint, client_id, client_secret, ...rest } = params;
  const formData = new URLSearchParams();
  formData.append('token', token);
  formData.append('token_type_hint', token_type_hint);
  formData.append('client_id', client_id);
  formData.append('client_secret', client_secret);
  return  request(`/api/revoke`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
    data: formData.toString(),
    ...(options || {}),
  });
}
//获取code
//客户端获取code后，给后端获取用户信息
//请求方法: POST
//请求地址: /api/oidc/code
export async function getOidcCode<OidcCodeResponse>(  data: OidcCodeRequest,   options?: { [key: string]: any }) {
  return  request<OidcCodeResponse>(`/api/oidc/code`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//MFA认证
//MFA认证，使用MFA验证码获取token
//请求方法: POST
//请求地址: /api/mfa/validate
export async function mfaValidate<AccessTokenResponse>(  data: MfaCode,   options?: { [key: string]: any }) {
  return  request<AccessTokenResponse>(`/api/mfa/validate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//获取Token
//获取Token，支持四种token获取方式
//请求方法: GET
//请求地址: /api/token
//参数名: client_id 参数类型: string 参数位置: query 是否必须: false  参数说明: client ID
//参数名: client_secret 参数类型: string 参数位置: query 是否必须: false  参数说明: client TotpSecret，使用场景:grant_type 为password;client_credentials
//参数名: grant_type 参数类型: string 参数位置: query 是否必须: false  参数说明: Grant Type
//参数名: code 参数类型: string 参数位置: query 是否必须: false  参数说明: StateCode
//参数名: redirect_uri 参数类型: string 参数位置: query 是否必须: false  参数说明: Redirect URI
//参数名: code_verifier 参数类型: string 参数位置: query 是否必须: false  参数说明: PKCE code verifier
//参数名: username 参数类型: string 参数位置: query 是否必须: false  参数说明: Username
//参数名: password 参数类型: string 参数位置: query 是否必须: false  参数说明: Password
export async function systemFetchTokenWithGet<AccessTokenResponse>(
  params: {
    client_secret?:string;// client TotpSecret，使用场景:grant_type 为password;client_credentials
    redirect_uri?:string;// Redirect URI
    code_verifier?:string;// PKCE code verifier
    password?:string;// Password
    client_id?:string;// client ID
    grant_type?:string;// Grant Type
    code?:string;// StateCode
    username?:string;// Username
  },
  options?: { [key: string]: any }) {
  return  request<AccessTokenResponse>(`/api/token`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json;application/x-www-form-urlencoded',
    },
    params: params,
    ...(options || {}),
  });
}
//令牌内省
//校验Token是否合法并查询信息，比如，access_token是否还有效，谁颁发的，颁发给谁的，scope又哪些等等的元数据信息
//请求方法: POST
//请求地址: /api/introspect
export async function systemIntrospect(
  params: {
    token: string;// 可以是access_token或者refresh_token
    token_type_hint: string;// 表示token的类型
  },
  options?: { [key: string]: any }) {
  const { token, token_type_hint, ...rest } = params;
  const formData = new FormData();
  formData.append('token', token);
  formData.append('token_type_hint', token_type_hint);
  return  request(`/api/introspect`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: formData,
    ...(options || {}),
  });
}
//用户名/手机号码/邮箱 + 密码登录
//用户名/手机号码/邮箱 + 密码登录
//请求方法: POST
//请求地址: /api/login/username
export async function systemLoginByUsername<AccessTokenResponse>(  data: LoginByUsername,   options?: { [key: string]: any }) {
  return  request<AccessTokenResponse>(`/api/login/username`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//openid configuration
//openid configuration
//请求方法: GET
//请求地址: /api/.well-known/openid-configuration
export async function systemOidcConfig(  options?: { [key: string]: any }) {
  return  request(`/api/.well-known/openid-configuration`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    ...(options || {}),
  });
}
//OIDC认证新用户注册接口
//OIDC认证新用户注册接口，若系统用户注册且当前认证方式需在系统中创建新用户时
//请求方法: POST
//请求地址: /api/register/oidc
export async function registerByOIDC<AccessTokenResponse>(  data: RegisterByOIDC,   options?: { [key: string]: any }) {
  return  request<AccessTokenResponse>(`/api/register/oidc`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//刷新Token
//刷新Token
//请求方法: POST
//请求地址: /api/refresh_token
//参数名: refresh_token 参数类型: string 参数位置: query 是否必须: false  参数说明: Refresh Token
export async function systemRefreshTokenWithGet<AccessTokenResponse>(
  params: {
    refresh_token?:string;// Refresh Token
  },
  options?: { [key: string]: any }) {
  return  request<AccessTokenResponse>(`/api/refresh_token`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    params: params,
    ...(options || {}),
  });
}
//LDAP方式登录
//LDAP登录，认证成功节后，若用户在系统不存在，则根据组织是否允许自动注册来决定是否自动创建用户信息，若第一次是通过第三方登录，需要先设置密码，若组织设置了MFA则需要再次输入验证码，若用户没有绑定过验证器，则返回验证器的二维码和密钥
//请求方法: POST
//请求地址: /api/login/ldap
export async function systemLoginByLdap<AccessTokenResponse>(  data: LoginByLDAP,   options?: { [key: string]: any }) {
  return  request<AccessTokenResponse>(`/api/login/ldap`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
    ...(options || {}),
  });
}
//刷新Token
//刷新Token
//请求方法: GET
//请求地址: /api/refresh_token
//参数名: refresh_token 参数类型: string 参数位置: query 是否必须: false  参数说明: Refresh Token
export async function systemRefreshTokenWithPost(
  params: {
    refresh_token?:string;// Refresh Token
  },
  options?: { [key: string]: any }) {
  return  request(`/api/refresh_token`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    params: params,
    ...(options || {}),
  });
}

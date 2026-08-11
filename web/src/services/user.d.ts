//默认语言
export type SetPassword = { 
  //主键
  //最大长度: 50
  id: string;
  //密码
  newPassword: string;
  //旧密码
  oldPassword: string;
}; 
// ShortUser 简单账户详情
export type ShortUser = { 
  //主键
  //最大长度: 50
  id: string;
  //创建时间
  createdAt: string;
  //用户名
  //最大长度: 255
  username: string;
  //昵称，如中文名
  //最大长度: 255
  nickname: string;
  //工号
  //最大长度: 255
  jobNumber?: string;
  //系统角色
  //默认值: none
  //最大长度: 255
  role: string;
  //是否有效
  //默认值: true
  enable: boolean;
  //邮箱
  //最大长度: 255
  email: string;
  //手机号码
  //最大长度: 255
  phone: string;
  //默认语言
  //默认值: zh
  //最大长度: 255
  language: string;
  //头像
  //最大长度: 1000
  avatar?: string;
  //是否绑定MFA认证
  mfa: boolean;
}; 
// UserCreate 账户信息创建
// 未来账户信息修改只能从eauth中
export type UserCreate = { 
  //用户名
  //最大长度: 255
  username: string;
  //昵称，如中文名
  //最大长度: 255
  nickname: string;
  //工号
  //最大长度: 255
  jobNumber?: string;
  //密码
  password?: string;
  //数据库保存的加密密码
  //最大长度: 255
  passwordStore?: string;
  //密码强度
  //默认值: weak
  //最大长度: 20
  passwordStrength: string;
  //系统角色
  //默认值: none
  //最大长度: 255
  role: string;
  //是否有效
  //默认值: true
  enable: boolean;
  //邮箱
  //最大长度: 255
  email: string;
  //手机号码
  //最大长度: 255
  phone: string;
  //默认语言
  //默认值: zh
  //最大长度: 255
  language: string;
}; 
// UserDetail 账户详情
export type UserDetail = { 
  //主键
  //最大长度: 50
  id: string;
  //创建时间
  createdAt: string;
  //创建时间
  updatedAt: string;
  //用户名
  //最大长度: 255
  username: string;
  //昵称，如中文名
  //最大长度: 255
  nickname: string;
  //工号
  //最大长度: 255
  jobNumber?: string;
  //密码
  password?: string;
  //数据库保存的加密密码
  //最大长度: 255
  passwordStore?: string;
  //密码强度
  //默认值: weak
  //最大长度: 255
  passwordStrength: string;
  //系统角色
  //默认值: none
  //最大长度: 255
  role: string;
  //是否有效
  //默认值: true
  enable: boolean;
  //邮箱
  //最大长度: 255
  email: string;
  //手机号码
  //最大长度: 255
  phone: string;
  //默认语言
  //默认值: zh
  //最大长度: 255
  language: string;
  //头像
  //最大长度: 1000
  avatar?: string;
  //是否绑定MFA认证
  mfa: boolean;
}; 
// UserDetailList  账户列表响应
export type UserDetailList = { 
  //当前页数据
  data?: ShortUser[];
  //数据库满足条件的数据总数
  total: number;
}; 
//重置密码状态码
export type UserMfa = { 
  //主键
  id?: number;
  mfa: boolean;
}; 
// UserResetPassword 账户修改密码
export type UserResetPassword = { 
  //主键
  //最大长度: 50
  id: string;
  //密码
  password: string;
  //重置密码状态码
  code?: string;
}; 
// UserRole 账户系统角色设置
// 设置账户在系统中的角色
export type UserRole = { 
  //主键
  ids: string[];
  //用户在系统中的角色
  //默认值: none
  //最大长度: 255
  role: string;
}; 
// UserStatus 账户信息禁用/启用
// 账户禁用后，用户将不能登陆该系统
export type UserStatus = { 
  //主键
  ids: string[];
  //是否有效
  enable: boolean;
}; 
// UserUpdate 账户信息更新
// 更新账户信息，未来只能在eauth中更新
export type UserUpdate = { 
  //主键
  //最大长度: 50
  id: string;
  //用户名
  //最大长度: 255
  username: string;
  //昵称，如中文名
  //最大长度: 255
  nickname: string;
  //工号
  //最大长度: 255
  jobNumber?: string;
  //默认语言
  //默认值: zh
  //最大长度: 255
  language: string;
}; 

// ValidateCodeCreate 系统操作校验码创建
export type ValidateCodeCreate = { 
  //用户ID
  //最大长度: 50
  userId: string;
  //验证类型
  //最大长度: 50
  category: string;
  //验证码
  //最大长度: 50
  code?: string;
  //动作
  //最大长度: 50
  action: string;
  //过期时间
  expired: string;
}; 
// ValidateCodeDetail 系统操作校验码详情
export type ValidateCodeDetail = { 
  //主键
  //最大长度: 50
  id: string;
  //创建时间
  createdAt: string;
  //创建时间
  updatedAt: string;
  //用户ID
  //最大长度: 50
  userId?: string;
  //验证类型
  //最大长度: 50
  category: string;
  //验证码
  //最大长度: 50
  code?: string;
  //动作
  //最大长度: 50
  action: string;
  //过期时间
  expired: string;
}; 
// ValidateCodeDetailList 系统操作校验码列表响应
export type ValidateCodeDetailList = { 
  //当前页数据
  data?: ValidateCodeDetail[];
  //数据库满足条件的数据总数
  total: number;
}; 

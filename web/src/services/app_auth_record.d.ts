// AppAuthRecordCreate 应用认证记录创建
export type AppAuthRecordCreate = { 
  //应用ID
  //最大长度: 50
  applicationId?: string;
  //响应编码 返回给浏览器客户端
  //最大长度: 50
  code: string;
  //PKCE挑战值
  //最大长度: 255
  codeChallenge?: string;
  //PKCE挑战方式
  //最大长度: 20
  codeChallengeMethod?: string;
  //用户ID
  //最大长度: 50
  userId: string;
  //OIDC nonce
  //最大长度: 255
  nonce?: string;
}; 
// AppAuthRecordDetail 应用认证记录详情
export type AppAuthRecordDetail = { 
  //主键
  //最大长度: 50
  id: string;
  //创建时间
  createdAt: string;
  //创建时间
  updatedAt: string;
  //应用ID
  //最大长度: 50
  applicationId?: string;
  //响应编码 返回给浏览器客户端
  //最大长度: 50
  code: string;
  //PKCE挑战值
  //最大长度: 255
  codeChallenge?: string;
  //PKCE挑战方式
  //最大长度: 20
  codeChallengeMethod?: string;
  //用户ID
  //最大长度: 50
  userId: string;
  //OIDC nonce
  //最大长度: 255
  nonce?: string;
}; 
// AppAuthRecordDetailList 应用认证记录列表响应
export type AppAuthRecordDetailList = { 
  //当前页数据
  data?: AppAuthRecordDetail[];
  //数据库满足条件的数据总数
  total: number;
}; 

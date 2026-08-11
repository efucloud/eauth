// ConfigCreate 配置创建
export type ConfigCreate = { 
  //配置名称
  //最大长度: 255
  name: string;
  //配置编码
  //最大长度: 255
  code: string;
  //描述
  //最大长度: 255
  description?: string;
  //值
  value: string;
  //值类型
  //最大长度: 20
  valType: string;
}; 
// ConfigDetail 配置详情
export type ConfigDetail = { 
  //主键
  //最大长度: 50
  id: string;
  //创建时间
  createdAt: string;
  //创建时间
  updatedAt: string;
  //配置名称
  //最大长度: 255
  name: string;
  //配置编码
  //最大长度: 255
  code: string;
  //描述
  //最大长度: 255
  description?: string;
  //值
  value: string;
  //值类型
  //最大长度: 20
  valType: string;
}; 
// ConfigDetailList 配置列表响应
export type ConfigDetailList = { 
  //当前页数据
  data?: ConfigDetail[];
  //数据库满足条件的数据总数
  total: number;
}; 
// ConfigUpdate 配置修改
export type ConfigUpdate = { 
  //主键
  //最大长度: 50
  id: string;
  //配置名称
  //最大长度: 255
  name: string;
  //配置编码
  //最大长度: 255
  code: string;
  //描述
  //最大长度: 255
  description?: string;
  //值
  value: string;
  //值类型
  //最大长度: 20
  valType: string;
}; 

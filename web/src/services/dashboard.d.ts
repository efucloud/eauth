//用户角色
export type ApplicationAuthTop = { 
  name?: string;
  code?: string;
  //最大长度: 50
  applicationId?: string;
  value?: number;
  home?: string;
  description?: string;
  //默认值: tenant
  //最大长度: 50
  scope: string;
}; 
// Dashboard 系统统计数据
export type Dashboard = { 
  //人脸数据
  faceRecognition?: number;
  //认证方式
  authProfile?: any;//todo 可能需要手动完善结构[];
  //普通应用数据
  application?: any;//todo 可能需要手动完善结构[];
  //用户角色
  userRole?: any;//todo 可能需要手动完善结构[];
}; 
// DashboardData 看板数据
export type DashboardData = { 
  name: string;
  value: any;
}; 

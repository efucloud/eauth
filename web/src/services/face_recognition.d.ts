// FaceRecognitionCreate 用户人脸识别信息创建
export type FaceRecognitionCreate = { 
  //用户ID
  //最大长度: 50
  userId: string;
  //是否有效
  //默认值: true
  enable: boolean;
  //录入名称
  //最大长度: 50
  name: string;
  //人脸数据
  faceIdData?: Float64Array;
}; 
// FaceRecognitionDetail 用户人脸识别信息详情
export type FaceRecognitionDetail = { 
  //主键
  //最大长度: 50
  id: string;
  //创建时间
  createdAt: string;
  //创建时间
  updatedAt: string;
  //用户ID
  //最大长度: 50
  userId: string;
  //是否有效
  //默认值: true
  enable: boolean;
  //录入名称
  //最大长度: 50
  name: string;
  //最近一次使用时间
  //最大长度: 255
  latestUsedTime?: string;
}; 
// FaceRecognitionDetailList 用户人脸识别信息
export type FaceRecognitionDetailList = { 
  //当前页数据
  data?: FaceRecognitionDetail[];
  //数据库满足条件的数据总数
  total: number;
}; 
// FaceRecognitionStatus 禁用后，用户将不能使用该认证方式登陆系统
export type FaceRecognitionStatus = { 
  //主键
  ids: number[];
  //是否有效
  enable: boolean;
}; 

import { AuthedUserInfo } from "@/services/common.d";
import { getCurrentOrgCode, isSystemPersonalCenter } from "@/utils/global";
 
const systemPersonalAccess = () => {
  return isSystemPersonalCenter()  
}
const systemAdminAccess = (currentUser?: AuthedUserInfo) => {
  if (!currentUser    ) {
    return false;
  }
  return currentUser.role === 'admin' &&!isSystemPersonalCenter();
  
}
const systemEditAccess = (currentUser?: AuthedUserInfo) => {
  if (!currentUser    ) {
    return false;
  }
  return (currentUser.role === 'admin' || currentUser.role === 'edit')&&!isSystemPersonalCenter();;
}
const systemViewAccess = (currentUser?: AuthedUserInfo) => {
  if (!currentUser    ) {
    return false;
  }
  return currentUser.role !== 'none'&&!isSystemPersonalCenter();
}


const systemUserAccess = (currentUser?: AuthedUserInfo) => {
  return currentUser?.id ||0 > 0;
}
const systemConsoleAccess = (currentUser?: AuthedUserInfo) => {
  return currentUser?.role !== 'none';
}
/**
 * @see https://umijs.org/zh-CN/plugins/plugin-access
 * */
// 权限分类  
// 系统权限：只有企业用户才能拥有，即登录的用户为global User
// admin: 管理员 能查看操作系统所有数据
// edit: 编辑，能查看和编辑系统数据
// view: 查看 能查看系统数据
// none: 普通用户 即无法进入系统后台查看系统本身的数据，只能查看个人数据
 
export default function access(initialState: { currentUser?: AuthedUserInfo } | undefined) {
  const { currentUser } = initialState ?? {};
  return {
    systemPersonalAccess: systemPersonalAccess(),
    systemAdminAccess: systemAdminAccess(currentUser),
    systemEditAccess: systemEditAccess(currentUser),
    systemViewAccess: systemViewAccess(currentUser),
    systemConsoleAccess: systemConsoleAccess(currentUser),
    systemUserAccess: systemUserAccess(currentUser),
  };
}

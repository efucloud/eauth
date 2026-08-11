import React, { useEffect } from 'react';
import { deleteToken } from '@/utils/global';

const Logout: React.FC = () => {
  const loginOut = async () => {
    // 删除全局Token
    deleteToken();
    window.location.href = '/user/login';
  };
  useEffect(() => {
    loginOut();
  }, []);
  return (<>
  </>);
}
export default Logout;
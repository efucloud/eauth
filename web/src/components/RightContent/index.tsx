import { GlobalOutlined } from '@ant-design/icons';
import { SelectLang as UmiSelectLang } from '@umijs/max';


export type SiderTheme = 'light' | 'dark';

export const SelectLang = () => {
  return <UmiSelectLang icon={<GlobalOutlined />} />;
};


import {
  DingtalkOutlined,
  AlipayOutlined,
  WeiboOutlined,
  WechatOutlined,
  GithubOutlined,
  GitlabOutlined,
  WechatWorkOutlined,
  QqOutlined,
  SafetyOutlined,
  TikTokOutlined,
  BaiduOutlined,
  BilibiliOutlined,
  QuestionCircleOutlined,
} from '@ant-design/icons';
import { Space, Divider, Image, Tooltip } from 'antd';
import { FormattedMessage } from '@umijs/max';
import { OIDC } from '@/services/common.d';


import React from 'react';
import type { CSSProperties } from 'react';
import { getColorPrimary } from '@/utils/global';

const iconStyle: CSSProperties = {
  color: 'rgba(0, 0, 0, 0.2)',
  fontSize: '25px',
  verticalAlign: 'middle',
  cursor: 'pointer',
};
const iconDivStyle: CSSProperties = {
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  flexDirection: 'column',
  height: 25,
  width: 25,
  // border: '1px solid #D4D8DD',
  // borderRadius: '40%',
}
interface ProviderProps {
  provider?: string; // 提供商
  height?: number; // 高度
}

export const Provider: React.FC<ProviderProps> = (props) => {
  if (props.height === 0) {
    props.height = 25
  }
  if (props?.height && props.height > Number(iconDivStyle.height)) {
    iconDivStyle.height = props.height;
    iconDivStyle.width = props.height;
  }

  return (
    <>
      <Space align='center' size={24} style={{ display: 'flex', justifyContent: 'center' }}>
        {props.provider === 'github' && <div style={iconDivStyle}>
          <GithubOutlined style={{ ...iconStyle, color: 'black' }} />
        </div>}
        {props.provider === 'gitlab' && <div style={iconDivStyle}>
          <GitlabOutlined style={{ ...iconStyle, color: 'rgb(252, 109, 38)' }} />
        </div>}
        {props.provider === 'wechat' && <div style={iconDivStyle}>
          <WechatOutlined style={{ ...iconStyle, color: 'rgb(51, 204, 0)' }} />
        </div>}
        {props.provider === 'alipay' && <div style={iconDivStyle}>
          <AlipayOutlined style={{ ...iconStyle, color: '#1677FF' }} />
        </div>}

        {props.provider === 'weibo' && <div style={iconDivStyle}>
          <WeiboOutlined style={{ ...iconStyle, color: '#e71f19' }} />
        </div>}
        {props.provider === 'wechatWork' && <div style={iconDivStyle}>
          <WechatWorkOutlined style={{ ...iconStyle, color: '#0082ef' }} />
        </div>}
        {props.provider === 'qq' && <div style={iconDivStyle}>
          <QqOutlined style={{ ...iconStyle, color: '#0082ef' }} />
        </div>}
        {props.provider === 'dingTalk' && <div style={iconDivStyle}>
          <DingtalkOutlined style={{ ...iconStyle, color: '#0082ef' }} />
        </div>}
        {props.provider === 'gitee' && <div style={iconDivStyle}>
          <Image preview={false} width={props.height} src='/logo/gitee.png' alt='gitee' />
        </div>}
        {props.provider === 'tiktok' && <div style={iconDivStyle}>
          <TikTokOutlined style={{ ...iconStyle, color: '#0082ef' }} />
        </div>}
        {props.provider === 'baidu' && <div style={iconDivStyle}>
          <BaiduOutlined style={{ ...iconStyle, color: '#0082ef' }} />
        </div>}
        {props.provider === 'feiShu' && <div style={iconDivStyle}>
          <Image preview={false} width={props.height} src='/logo/feishu.webp' alt='feiShu' />
        </div>}
        {props.provider === 'bilibili' && <div style={iconDivStyle}>
          <BilibiliOutlined style={{ ...iconStyle, color: '#0082ef' }} />
        </div>}
        {props.provider === 'custom' && <div style={iconDivStyle}>
          <SafetyOutlined style={{ ...iconStyle, color: '#0082ef' }} />
        </div>}

      </Space>
    </>
  );
};


export const ProviderInLine: React.FC<ProviderProps> = (props) => {
  return (
    <>
      <Space align='center' size={24} style={{ display: 'left' }}>
        {props.provider === 'github' &&
          <Space> <GithubOutlined style={{ ...iconStyle, color: 'black' }} /><FormattedMessage id='model.provider.oidc.category.github'  /></Space>}
        {props.provider === 'gitlab' &&
          <Space> <GitlabOutlined style={{ ...iconStyle, color: 'rgb(252, 109, 38)' }} /><FormattedMessage id='model.provider.oidc.category.gitlab'  /></Space>
        }
        {props.provider === 'wechat' &&
          <Space>   <WechatOutlined style={{ ...iconStyle, color: 'rgb(51, 204, 0)' }} /><FormattedMessage id='model.provider.oidc.category.wechat'  /> </Space>
        }
        {props.provider === 'alipay' &&
          <Space>   <AlipayOutlined style={{ ...iconStyle, color: '#1677FF' }} /><FormattedMessage id='model.provider.oidc.category.alipay'  /> </Space>
        }

        {props.provider === 'weibo' &&
          <Space>   <WeiboOutlined style={{ ...iconStyle, color: '#e71f19' }} /><FormattedMessage id='model.provider.oidc.category.weibo'  /> </Space>
        }
        {props.provider === 'wechatWork' &&
          <Space>  <WechatWorkOutlined style={{ ...iconStyle, color: '#0082ef' }} /><FormattedMessage id='model.provider.oidc.category.wechatWork'  /> </Space>
        }
        {props.provider === 'qq' &&
          <Space>  <QqOutlined style={{ ...iconStyle, color: '#0082ef' }} /><FormattedMessage id='model.provider.oidc.category.qq'  /> </Space>
        }
        {props.provider === 'dingTalk' &&
          <Space>  <DingtalkOutlined style={{ ...iconStyle, color: '#0082ef' }} /><FormattedMessage id='model.provider.oidc.category.dingTalk'  /> </Space>
        }
        {props.provider === 'tiktok' &&
          <Space>  <TikTokOutlined style={{ ...iconStyle, color: '#0082ef' }} /><FormattedMessage id='model.provider.oidc.category.tiktok'  /> </Space>
        }
        {props.provider === 'baidu' &&
          <Space>  <BaiduOutlined style={{ ...iconStyle, color: '#0082ef' }} /><FormattedMessage id='model.provider.oidc.category.baidu'  /> </Space>
        }
        {props.provider === 'bilibili' &&
          <Space>  <BilibiliOutlined style={{ ...iconStyle, color: '#0082ef' }} /><FormattedMessage id='model.provider.oidc.category.bilibili'  /> </Space>
        }
        {props.provider === 'gitee' &&
          <Space> <Image preview={false} width={25} src='/logo/gitee.png' alt='gitee' /><FormattedMessage id='model.provider.oidc.category.gitee'  /> </Space>
        }
        {props.provider === 'feiShu' &&
          <Space> <Image preview={false} width={25} src='/logo/feishu.webp' alt='gitee' /><FormattedMessage id='model.provider.oidc.category.feiShu'  /> </Space>
        }
        {props.provider === 'custom' &&
          <Space> <SafetyOutlined style={{ ...iconStyle, color: '#0082ef' }} /> <FormattedMessage id='model.provider.oidc.category.custom'  /> </Space>
        }
      </Space>
    </>
  );
};

export const ProviderTwoLine: React.FC<ProviderProps> = (props) => {

  const widthParam=props?.height?props.height:40;
  const fontSizeParam= props?.height?`${props.height}px`:'40px';
  return (

    <div style={{ textAlign: 'center' }}>
      <Space align='center' size={24} style={{ display: 'left' }}>
        {props.provider === 'github' &&
          <div>
            <GithubOutlined style={{ ...iconStyle, fontSize: `${fontSizeParam}`, color: 'black' }} />
            <br />
            <FormattedMessage id='model.provider.oidc.category.github'  />
          </div>}
        {props.provider === 'gitlab' &&
          <div>
            <GitlabOutlined style={{ ...iconStyle, fontSize: `${fontSizeParam}`, color: 'rgb(252, 109, 38)' }} />
            <br />
            <FormattedMessage id='model.provider.oidc.category.gitlab'  />
          </div>
        }
        {props.provider === 'wechat' &&
          <div>
            <WechatOutlined style={{ ...iconStyle, fontSize: `${fontSizeParam}`, color: 'rgb(51, 204, 0)' }} />
            <br />
            <FormattedMessage id='model.provider.oidc.category.wechat'  />
          </div>

        }
        {props.provider === 'alipay' &&
          <div>
            <AlipayOutlined style={{ ...iconStyle, fontSize: `${fontSizeParam}`, color: '#1677FF' }} />
            <br />
            <FormattedMessage id='model.provider.oidc.category.alipay'  />
          </div>
        }

        {props.provider === 'weibo' &&
          <div>
            <WeiboOutlined style={{ ...iconStyle, fontSize: `${fontSizeParam}`, color: '#e71f19' }} />
            <br />
            <FormattedMessage id='model.provider.oidc.category.weibo'  />
          </div>
        }
        {props.provider === 'wechatWork' &&
          <div>
            <WechatWorkOutlined style={{ ...iconStyle, fontSize: `${fontSizeParam}`, color: '#0082ef' }} />
            <br />
            <FormattedMessage id='model.provider.oidc.category.wechatWork'  />
          </div>
        }
        {props.provider === 'qq' &&
          <div>
            <QqOutlined style={{ ...iconStyle, fontSize: `${fontSizeParam}`, color: '#0082ef' }} />
            <br />
            <FormattedMessage id='model.provider.oidc.category.qq'  />
          </div>
        }
        {props.provider === 'dingTalk' &&
          <div>
            <DingtalkOutlined style={{ ...iconStyle, fontSize: `${fontSizeParam}`, color: '#0082ef' }} />
            <br />
            <FormattedMessage id='model.provider.oidc.category.dingTalk'  /> </div>
        }
        {props.provider === 'tiktok' &&
          <div>
            <TikTokOutlined style={{ ...iconStyle, fontSize:`${fontSizeParam}`, color: '#0082ef' }} />
            <br />
            <FormattedMessage id='model.provider.oidc.category.tiktok'  />
          </div>
        }
        {props.provider === 'baidu' &&
          <div>
            <BaiduOutlined style={{ ...iconStyle, fontSize: `${fontSizeParam}`, color: '#0082ef' }} />
            <br />
            <FormattedMessage id='model.provider.oidc.category.baidu'  />
          </div>
        }
        {props.provider === 'bilibili' &&
          <div>
            <BilibiliOutlined style={{ ...iconStyle, fontSize: `${fontSizeParam}`, color: '#0082ef' }} />
            <br />
            <FormattedMessage id='model.provider.oidc.category.bilibili'  />
          </div>
        }
        {props.provider === 'gitee' &&
          <div>
            <Image preview={false} width={widthParam} src='/logo/gitee.png' alt='gitee' />
            <br />
            <FormattedMessage id='model.provider.oidc.category.gitee'  />
          </div>
        }
        {props.provider === 'feiShu' &&
          <div>
            <Image preview={false} width={widthParam} src='/logo/feishu.webp' alt='gitee' />
            <br />
            <FormattedMessage id='model.provider.oidc.category.feiShu'  />
          </div>
        }
        {props.provider === 'custom' &&
          <div>
            <SafetyOutlined style={{ ...iconStyle, fontSize: `${fontSizeParam}`, color: '#0082ef' }} />
            <br />
            <FormattedMessage id='model.provider.oidc.category.custom'  />
          </div>
        }
      </Space>
    </div>
  );
};

interface LoginMethods {
  oidcs: OIDC[], // oidc 登录方式
  allowRegistry: boolean  // 是否允许注册
}

export const ProviderActionIcons: React.FC<LoginMethods> = (props) => {
  const colorPrimary=getColorPrimary();
  return (
    <>
      {props.oidcs.length > 0 && <>
        <Divider plain><FormattedMessage id="pages.login.loginWith" defaultMessage="其他登录方式" />
          {props.allowRegistry ? <Tooltip title={<FormattedMessage id="pages.login.allowRegistry.true.description"  />}>
            <QuestionCircleOutlined />
          </Tooltip> : <Tooltip title={<FormattedMessage id="pages.login.allowRegistry.false.description"  />}>
            <QuestionCircleOutlined />
          </Tooltip>}
        </Divider>
        <Space align='center' size={24} style={{ display: 'flex', justifyContent: 'center' }}>
          {props.oidcs.map((item) => (
            <div key={item.category} style={iconDivStyle}>
              <a style={{color:colorPrimary}} href={item.address}>
                {item.category === 'github' && <GithubOutlined style={{ ...iconStyle, color: 'black' }} />}
                {item.category === 'gitlab' && <GitlabOutlined style={{ ...iconStyle, color: 'rgb(252, 109, 38)' }} />}
                {item.category === 'wechat' && <WechatOutlined style={{ ...iconStyle, color: 'rgb(51, 204, 0)' }} />}
                {item.category === 'alipay' && <AlipayOutlined style={{ ...iconStyle, color: '#1677FF' }} />}
                {item.category === 'weibo' && <WeiboOutlined style={{ ...iconStyle, color: '#e71f19' }} />}
                {item.category === 'wechatWork' && <WechatWorkOutlined style={{ ...iconStyle, color: '#0082ef' }} />}
                {item.category === 'qq' && <QqOutlined style={{ ...iconStyle, color: '#0082ef' }} />}
                {item.category === 'dingTalk' && <DingtalkOutlined style={{ ...iconStyle, color: '#0082ef' }} />}
                {item.category === 'baidu' && <BaiduOutlined style={{ ...iconStyle, color: '#0082ef' }} />}
                {item.category === 'tiktok' && <TikTokOutlined style={{ ...iconStyle, color: '#0082ef' }} />}
                {item.category === 'bilibili' && <BilibiliOutlined style={{ ...iconStyle, color: '#0082ef' }} />}
                {item.category === 'gitee' && <Image preview={false} width={25} src='/logo/gitee.png' alt='gitee' />}
                {item.category === 'feishu' && <Image preview={false} width={25} src='/logo/feishu.webp' alt='feishu' />}
                {item.category === 'custom' && <SafetyOutlined style={{ ...iconStyle, color: '#0082ef' }} />}
              </a>
            </div>
          ))}
        </Space>
      </>
      }
    </>
  );
}
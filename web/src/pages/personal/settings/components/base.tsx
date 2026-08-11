import { UploadOutlined } from '@ant-design/icons';
import {
  ProForm,
  ProFormText,
  ProFormSelect,
} from '@ant-design/pro-components';
import { Button, message, Upload, Row, Col } from 'antd';
import React, { useEffect, useState } from 'react';
import { systemGetUserinfo, updateSelfInfo } from '@/services/oauth.api';
import { updateSelfAvatar } from '@/services/personal.api';
import useStyles from './index.style';
import { AuthedUserInfo } from "@/services/common.d";
import { FormattedMessage, useIntl } from '@umijs/max';
import { UserUpdate } from '@/services/user';



const BaseView: React.FC = () => {
  const [currentUser, setCurrentUser] = useState<AuthedUserInfo>();
  const intl = useIntl();
  const { styles } = useStyles();
  const getCurrentUser = async () => {
    const data = await systemGetUserinfo();
    setCurrentUser(data as AuthedUserInfo);
  }
  useEffect(() => { getCurrentUser(); }, []);
  const uploadImage = async (options: any) => {

    const { file } = options;
    const res = await updateSelfAvatar({ file }) as string;
    if (currentUser && currentUser !== undefined) {
      setCurrentUser((currentUser) => ({
        ...currentUser,
        avatar: res,
      }));
    }
  };
  // 头像组件 方便以后独立，增加裁剪之类的功能
  const AvatarView = ({ avatar }: { avatar: string }) => (
    <>
      <div className={styles.avatar_title}><FormattedMessage id='pages.avatar' /></div>
      <div className={styles.avatar}>
        <img src={avatar} alt="avatar" />
      </div>
      <Upload showUploadList={false} accept='.png,.jpg,.jpeg' maxCount={1} customRequest={uploadImage} >
        <div className={styles.button_view}>
          <Button>
            <UploadOutlined />
            <FormattedMessage id='pages.change.avatar' />
          </Button>
        </div>
      </Upload>
    </>
  );
  const languages = [
    {
      label: intl.formatMessage({ id: 'model.user.language.zh' }),
      value: 'zh-CN',
    },
    {
      label: intl.formatMessage({ id: 'model.user.language.en' }),
      value: 'en-US',
    },
  ]
  const getAvatarURL = () => {
    if (currentUser) {
      if (currentUser.avatar) {
        return currentUser.avatar;
      }
      return '';
    }
    return '';
  };
  const handleFinish = async (values: UserUpdate) => {
    values.id = currentUser?.id || 0;
    await updateSelfInfo(values);
    message.success(intl.formatMessage({ id: 'pages.operation.update.success' }));
    getCurrentUser();
  };
  return (
    <div  >
      {currentUser ?
        <>
          <div  >
            <ProForm
              layout="vertical"
              onFinish={handleFinish}
              submitter={{
                searchConfig: {
                  submitText: intl.formatMessage({ id: 'pages.personal.settings.baseinfo.update' }),
                },
                render: (_, dom) => dom[1],
              }}
              initialValues={currentUser}
            >
              <Row gutter={{ xs: 16, sm: 32, md: 64, lg: 128 }}>
                <Col lg={16} md={16} sm={24} >
                  <Row gutter={{ xs: 16, sm: 32, md: 64, lg: 128 }}>
                    <Col lg={12} md={12} sm={24} >
                      <ProFormText
                        width="md"
                        name="username"
                        label={intl.formatMessage({ id: 'model.user.username' })}
                        rules={[
                          {
                            required: true,
                            message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.user.username' }),
                          },
                        ]}
                      /></Col>
                    <Col lg={12} md={12} sm={24} >
                      <ProFormText
                        width="md"
                        label={intl.formatMessage({ id: 'model.user.nickname' })}
                        name="nickname"
                        rules={[
                          {
                            required: true,
                            message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'model.user.nickname' }),
                          },
                        ]}
                      /></Col>
                    <Col lg={12} md={12} sm={24} >
                      <ProFormText
                        width="md"
                        label={intl.formatMessage({ id: 'model.user.jobNumber' })}
                        name="jobNumber"
                      /></Col>
                    <Col lg={12} md={12} sm={24} >
                      <ProFormSelect
                        label={intl.formatMessage({ id: 'model.user.language' })}
                        name="language"
                        options={languages}
                      /></Col>
                  </Row>

                </Col>
                <Col lg={8} md={8} sm={24} >
                  <AvatarView avatar={getAvatarURL()} /></Col>
              </Row>
            </ProForm>
          </div>

        </> : null}
    </div>
  );
};
export default BaseView;

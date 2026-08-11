import { Divider, Tooltip, List, Popconfirm, Button, message, Card } from 'antd';
import React, { Fragment, useEffect, useState, Suspense, lazy } from 'react';
import { Typography } from 'antd';
const { Text, Paragraph, Link } = Typography;
import * as faceapi from "face-api.js";
import { QuestionCircleOutlined } from "@ant-design/icons";
import { getColorPrimary } from '@/utils/global';
import { IntlShape } from "react-intl";
import { ModalForm, ProDescriptions, ProFormText, } from '@ant-design/pro-components';
import { FormattedMessage, useIntl } from '@umijs/max';
import { createUserFaceRecognition, getFaceRecognitionPersonal, deleteUserFaceRecognition, changeUserFaceRecognitionStatus } from '@/services/face_recognition.api'
import { SystemFaceRecognitionCreate, SystemFaceRecognitionStatus, SystemFaceRecognitionDetailList, SystemFaceRecognitionDetail } from '@/services/system_face_recognition.d'
import { BatchOperationIds } from '@/services/common.d';
const FaceRecognitionModal = lazy(() => import('@/components/FaceRecognition'));


const FaceRecognition: React.FC = () => {
  const colorPrimary = getColorPrimary();
  const intl = useIntl();
  const [modalVisible, setModalVisible] = useState(false);
  const [faceName, setFaceName] = useState('');
  const [faceNameModalVisible, setFaceNameModalVisible] = useState(false);
  const [personalFaceIdDatas, setPersonalFaceIdDatas] = useState<SystemFaceRecognitionDetail[]>([]);
  const getPersonalFaceIdData = async () => {
    const data = await getFaceRecognitionPersonal();
    const list = data as SystemFaceRecognitionDetailList;
    if (list.data) {
      setPersonalFaceIdDatas(list.data);
    }
  };
  useEffect(() => {
    getPersonalFaceIdData();

  }, []);
  const addFaceId = async (faceIdData: number[]) => {
    let data = {} as SystemFaceRecognitionCreate;
    data.enable = true;
    data.name = faceName;
    data.faceIdData = faceIdData;
    await createUserFaceRecognition(data);
    message.success(intl.formatMessage({ id: 'pages.operation.update.success' }));
    getPersonalFaceIdData();
  };
  const onFaceNameFinish = async (values: Record<string, any>) => {
    if (values?.name) {
      setFaceName(values.name);
      setFaceNameModalVisible(false);
      setModalVisible(true);
    }
  }
  const handleEnable = async (intl: IntlShape, enable: boolean, selectedRows: SystemFaceRecognitionDetail[]) => {
    const hide = message.loading(intl.formatMessage({ id: 'pages.operation.updating' }));
    if (!selectedRows) return true;

    try {
      await changeUserFaceRecognitionStatus({
        ids: selectedRows.map((row) => row.id),
        enable: enable,
      } as SystemFaceRecognitionStatus);
      hide();
      message.success(intl.formatMessage({ id: 'pages.operation.update.success' }));
      await getPersonalFaceIdData();
      return true;
    } catch (error) {
      hide();
      message.error(intl.formatMessage({ id: 'pages.operation.update.failed' }));
      return false;
    }
  };
  const handleDelete = async (intl: IntlShape, selectedRows: SystemFaceRecognitionDetail[]) => {
    const hide = message.loading(intl.formatMessage({ id: 'pages.operation.updating' }));
    if (!selectedRows) return true;

    try {
      await deleteUserFaceRecognition({
        ids: selectedRows.map((row) => row.id),
      } as BatchOperationIds);
      hide();
      message.success(intl.formatMessage({ id: 'pages.operation.update.success' }));
      await getPersonalFaceIdData();
      return true;
    } catch (error) {
      hide();
      message.error(intl.formatMessage({ id: 'pages.operation.update.failed' }));
      return false;
    }
  };

  return (
    <Fragment>
      <Button type='primary' onClick={() => setFaceNameModalVisible(true)}><FormattedMessage id='pages.operation.add' /></Button>
      <Divider style={{ color: colorPrimary }}></Divider>
      <List
        itemLayout="horizontal"
        grid={{ gutter: 16, column: 3 }}
        dataSource={personalFaceIdDatas}
        renderItem={(item) => (
          <List.Item    >
            <Card
              title={item.name}
              extra={<>
                <Popconfirm
                  key={item.id + '-unbind'}
                  title=''
                  description={intl.formatMessage({ id: 'pages.personal.settings.unbind.description' })}

                  okText={intl.formatMessage({ id: 'pages.operation.confirm' })}
                  cancelText={intl.formatMessage({ id: 'pages.operation.cancel' })}
                  onConfirm={() => {
                    handleDelete(intl, [item]);
                  }}
                >
                  <a><FormattedMessage id="pages.personal.settings.unbind" defaultMessage="解绑" /></a>
                </Popconfirm>&nbsp;&nbsp;
                {item.enable ? <Popconfirm
                  key={item.id + '-disable'}
                  title=''
                  description={intl.formatMessage({ id: 'pages.personal.settings.bind.disable.description' })}
                  okText={intl.formatMessage({ id: 'pages.operation.confirm' })}
                  cancelText={intl.formatMessage({ id: 'pages.operation.cancel' })}
                  onConfirm={() => {
                    handleEnable(intl, false, [item]);
                  }}
                >
                  <a style={{color:colorPrimary}} className='delete'><FormattedMessage id="pages.personal.settings.bind.disable" defaultMessage="禁用" /></a>
                </Popconfirm> : <Popconfirm
                  key={item.id + '-enable'}
                  title=''
                  description={intl.formatMessage({ id: 'pages.personal.settings.bind.enable.description' })}
                  okText={intl.formatMessage({ id: 'pages.operation.confirm' })}
                  cancelText={intl.formatMessage({ id: 'pages.operation.cancel' })}
                  onConfirm={() => {
                    handleEnable(intl, true, [item]);
                  }}
                >
                  <a style={{color:colorPrimary}} ><FormattedMessage id="pages.personal.settings.bind.enable" defaultMessage="启用" /></a>
                </Popconfirm>}
              </>}
            >
              <Card.Meta
                description={
                  <ProDescriptions column={1}>
                    <ProDescriptions.Item label={intl.formatMessage({ id: 'model.createdAt' })} valueType='dateTime'>{item.createdAt}</ProDescriptions.Item>

                  </ProDescriptions>
                }
              />
            </Card>

          </List.Item>
        )}
      />
      {modalVisible === true && <Suspense>
        <FaceRecognitionModal visible={modalVisible} onSuccess={(faceIdData: number[]) => {
          setModalVisible(false);
          message.success(intl.formatMessage({ id: 'pages.personal.settings.faceRecognition.face.data.extract.success' }));
          addFaceId(faceIdData);
        }} onCancel={() => { setModalVisible(false) }} />
      </Suspense>}
      <ModalForm
        title={intl.formatMessage({ id: 'pages.personal.settings.faceRecognition.face.data.name.title' })}
        width="400px"
        open={faceNameModalVisible}
        onOpenChange={setFaceNameModalVisible}
        onFinish={onFaceNameFinish}

      >
        <ProFormText
          label={<>
            <FormattedMessage id='pages.personal.settings.faceRecognition.face.data.name' />
            <Tooltip title={<FormattedMessage id='pages.personal.settings.faceRecognition.face.data.name.description' />}>
              <QuestionCircleOutlined />
            </Tooltip> </>}
          name="name"
          rules={[
            {
              required: true,
              message: intl.formatMessage({ id: 'pages.input.text.tips' }) + intl.formatMessage({ id: 'pages.personal.settings.faceRecognition.face.data.name' }),
            },
          ]}
        />
      </ModalForm>
    </Fragment>
  );
};


export default FaceRecognition;
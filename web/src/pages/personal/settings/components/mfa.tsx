import React, { Fragment, useEffect, useState } from 'react';
import { ProDescriptions } from '@ant-design/pro-components';
import { FormattedMessage, useIntl } from '@umijs/max';
import { Button, Image, Input, Popconfirm } from "antd";
import { getMultiFactorAuthPersonal, resetMultiFactorAuthPersonal, boundMultiFactorAuthPersonal } from '@/services/mfa.api';
import { MultiFactorAuthDetail, PersonalBoundMFA } from '@/services/mult_factor_auth.d';

const Mfa: React.FC = () => {
  const intl = useIntl();
  const [mfaData, setMfaData] = useState<MultiFactorAuthDetail>();
  const [value, setValue] = useState('');
  const getCurrentUserPersonalWebauthnInfos = async () => {
    const data = await getMultiFactorAuthPersonal() as MultiFactorAuthDetail;
    setMfaData(data);
  };
  useEffect(() => {
    getCurrentUserPersonalWebauthnInfos();
  }, []);
  const mfaBound = async () => {
    const data = await boundMultiFactorAuthPersonal({ code: value  } as PersonalBoundMFA) as MultiFactorAuthDetail;
    setMfaData(data);
  };
  const mfaResetBound = async () => {
    const data = await resetMultiFactorAuthPersonal({ code: value,  } as PersonalBoundMFA) as MultiFactorAuthDetail;
    setMfaData(data);
  };

  return (
    <Fragment>
      <div style={{ textAlign: 'center' }}>
        {mfaData && mfaData.id > 0 && <div>
          <Image
            preview={false}
            width={200} src={mfaData.image} />
          <br />
          <br />
          {mfaData.secret !== '' && <>
            <ProDescriptions column={1}>
              <ProDescriptions.Item copyable >{mfaData.secret}</ProDescriptions.Item>
            </ProDescriptions>
          </>}
          {mfaData.status === 'bound' && <>
          <Input style={{ width: '200px' }} placeholder={intl.formatMessage({ id: 'pages.validate.code' })} allowClear onChange={(e) => {
            setValue(e.target.value); }} />
          <Popconfirm
            title={<FormattedMessage id='pages.operation.mfa.reset'  />}
            description={<FormattedMessage id='pages.personal.settings.mfa.self.reset.description'  />}
            onConfirm={mfaResetBound}
            okText={intl.formatMessage({ id: 'pages.operation.confirm'})}
            cancelText={intl.formatMessage({ id: 'pages.operation.cancel'})}
          >&nbsp;&nbsp;<Button key="confirm" danger>
              <FormattedMessage id='pages.operation.mfa.reset'  />
            </Button></Popconfirm></> }
         
        </div>}
      </div>
    </Fragment>

  );
};
export default Mfa;


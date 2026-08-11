import { PageContainer} from '@ant-design/pro-components';
import { Row, Col, Card, List } from 'antd';
import { useIntl, FormattedMessage } from '@umijs/max';
import { StarOutlined,LinkOutlined } from '@ant-design/icons';
import React, { useEffect, useState } from 'react';
import { systemDashboard, applicationAuth30Days, applicationAuthTop10 } from '@/services/dashboard.api';
import RcResizeObserver from 'rc-resize-observer';
import { StatisticCard } from '@ant-design/pro-components';
import { Dashboard, DashboardData, ApplicationAuthTop } from '@/services/dashboard.d';
import { CustomChatPie, CustomChatLine } from '@/components';
import { getColorPrimary } from '@/utils/global';
const { Divider } = StatisticCard;

const DashboardIndex: React.FC = () => {
  const [responsive, setResponsive] = useState(false);
  const colorPrimary = getColorPrimary();

  const intl = useIntl();
  const [dashboard, setDashboard] = useState<Dashboard>();
  const [app30Days, setApp30Days] = useState<DashboardData[]>([]);
  const [appTop10, setAppTop10] = useState<ApplicationAuthTop[]>([]);
  const getDashboard = async () => {
    const data = await systemDashboard() as Dashboard;
    data?.authProfile?.forEach((item: DashboardData) => {
      item.intlName = intl.formatMessage({ id: `model.provider.oidc.category.${item.name}`, defaultMessage: item.name });
    });
    data?.organizationCategory?.forEach((item: DashboardData) => {
      item.intlName = intl.formatMessage({ id: `dashboard.organization.category.${item.name}`, defaultMessage: item.name });
    });
    
    data?.userRole?.forEach((item: DashboardData) => {
      item.intlName = intl.formatMessage({ id: `model.user.role.${item.name}`, defaultMessage: item.name });
    });
    setDashboard(data);
  }
  const getApplicationAuth30Days = async () => {
    const data = await applicationAuth30Days() as DashboardData[];
    if (data) {
      data.sort((a, b) => new Date(a.name).getTime() - new Date(b.name).getTime());
      data.forEach((item: DashboardData) => {
        item[intl.formatMessage({id:'pages.dashboard.application.auth.time'})] =item.name;
        delete item.name;
        item[intl.formatMessage({id:'pages.dashboard.application.auth.count'})]=item.value;
        delete item.value;
      })
      setApp30Days(data);
    }
  };
  const getAuthTop10 = async () => {
    const data = await applicationAuthTop10() as ApplicationAuthTop[];
    if (data) {
      setAppTop10(data);
    }
  }
  useEffect(() => {
    getDashboard();
    getApplicationAuth30Days();
    getAuthTop10();
  }, []);

  return (
    <>
      <PageContainer title={false} >
        {dashboard &&
          <RcResizeObserver
            key="resize-observer"
            onResize={(offset) => {
              setResponsive(offset.width < 596);
            }}
          >
            <StatisticCard.Group direction={responsive ? 'column' : 'row'}  >
              <StatisticCard
                title={intl.formatMessage({ id: 'dashboard.oidc' })}
                key='oidc'
                chart={<CustomChatPie data={dashboard?.authProfile || []} />}
              />
             
              <Divider type={responsive ? 'horizontal' : 'vertical'} />
              <StatisticCard
                title={intl.formatMessage({ id: 'model.user.role' })}
                key='role'
                chart={<CustomChatPie data={dashboard?.userRole || []} />}
              />
            </StatisticCard.Group>
          </RcResizeObserver>
        }
        <Row gutter={24} style={{ marginTop: 24 }}>
          <Col xl={16} lg={24} md={24} sm={24} xs={24}>
            <Card
              title={<FormattedMessage id='pages.dashboard.application.auth.30days'  />}
            >
              {app30Days && <CustomChatLine data={app30Days} xField={intl.formatMessage({id:'pages.dashboard.application.auth.time'})} yField={intl.formatMessage({id:'pages.dashboard.application.auth.count'})} />}
            </Card>
          </Col>
          <Col xl={8} lg={24} md={24} sm={24} xs={24}>
            <Card
              title={<FormattedMessage id='pages.dashboard.application.auth.top10'  />}
            >
              <List
                itemLayout="horizontal"
                dataSource={appTop10}
                renderItem={(item) => (
                  <List.Item>
                    <List.Item.Meta
                      title={<>{item.name}({item.code})&nbsp;&nbsp;<StarOutlined style={{ color: colorPrimary }} />{item.value}&nbsp;&nbsp;<a style={{color:colorPrimary}} href={item.home} rel="noreferrer" target='_blank'><LinkOutlined /></a></>}
                      description={<>{item.description}</>}
                    />
                  </List.Item>
                )}
              />
            </Card>
          </Col>
        </Row>
      </PageContainer>
    </>
  );
};
export default DashboardIndex;



import { PageContainer } from '@ant-design/pro-components';
import React, { useEffect, useState } from 'react';
import { Outlet, FormattedMessage, useIntl } from '@umijs/max';
import { Input } from 'antd';
import { listApplication } from '@/services/application.api'
import { ApplicationDetail, ApplicationDetailList } from '@/services/application';
import { Card, List, Typography, Avatar, Skeleton, Divider, Tooltip } from 'antd';
const { Paragraph } = Typography;
import InfiniteScroll from 'react-infinite-scroll-component';
import debounce from 'lodash/debounce';
import { getColorPrimary } from '@/utils/global';
import useStyles from './style.style';
import { User } from '@/components';


const IndexView: React.FC = () => {
  const intl = useIntl();
  const { styles } = useStyles();
  const [debouncedValue, setDebouncedValue] = useState('');
  const [current, setCurrent] = useState(1);
  const [total, setTotal] = useState(0);
  const [apps, setApps] = useState<ApplicationDetail[]>([]);
  const colorPrimary = getColorPrimary();

  const debouncedHandleChange = debounce((value) => {
    setDebouncedValue(value);
    setApps([]);
    setCurrent(0);
  }, 1000); // 500 毫秒防抖

  const loadMoreData = async () => {
    const data = await listApplication({  name: debouncedValue, current: current }) as ApplicationDetailList;
    setCurrent(current + 1);
    if (data?.data) {
      setApps([...apps, ...data.data]);
    }
    setTotal(data.total);
  }
  useEffect(() => { loadMoreData(); }, [debouncedValue]);

  return (
    <PageContainer
      title={false}
      content={
        <div style={{ textAlign: 'center' }}>
          <Input
            placeholder={intl.formatMessage({ id: 'pages.input.application' })}
            size="large"
            allowClear={true}
            style={{ maxWidth: 522, width: '100%' }}
            onChange={(e) => {
              debouncedHandleChange(e.target.value);
            }}
            onSearch={loadMoreData}
          />
        </div>
      }
    >
      <InfiniteScroll
        dataLength={apps.length}
        next={() => { loadMoreData(); }}
        hasMore={apps.length < total}
        loader={<Skeleton avatar paragraph={{ rows: 1 }} active />}
        endMessage={<Divider plain><FormattedMessage id='pages.no.more.data' /> 🤐</Divider>}
        scrollableTarget="scrollableDiv"
      >
        <List<ApplicationDetail>
          rowKey='id'
          grid={{ gutter: 16, xs: 1, sm: 2, md: 3, lg: 3, xl: 4, xxl: 4 }}
          dataSource={apps}
          renderItem={(item) => (
            <List.Item>
              <Card
                hoverable
                className={styles.card}
                onClick={() => {
                  if(item.home){
                    window.open(item.home, '_blank');
                  }
                }
                }
              >
                <Card.Meta
                  avatar={<img alt="" className={styles.cardAvatar} src={item.logo || '/logo.png'} />}
                  title={<a>{item.name}</a>}
                  description={<Tooltip color={colorPrimary} placement='bottom' title={item.description}>
                    <Paragraph 
                      className={styles.item}
                      ellipsis={{
                        rows: 4,
                      }}
                    >
                      {item.description}
                    </Paragraph>
                  </Tooltip>
                  }
                />
              </Card>
            </List.Item>
          )}
        />
      </InfiniteScroll>
      <Outlet />
    </PageContainer>
  );

};

export default IndexView;
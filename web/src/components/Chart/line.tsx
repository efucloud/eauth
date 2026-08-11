import React from 'react';
import { DashboardData } from '@/services/dashboard.d';
import { FormattedMessage } from '@umijs/max';
import { Line } from '@ant-design/charts';

interface ChartProps {
    data: DashboardData[],
    height?: number,
    position?: string,
    xField: string,
    yField: string
}
export const CustomChatLine: React.FC<ChartProps> = (props) => {
    if (props.data.length > 0) {
        const config = {
            data: props.data || [],
            xField: props.xField,
            yField: props.yField,
            point: {
              shapeField: 'square',
              sizeField: 4,
            },
            interaction: {
              tooltip: {
                marker: false,
              },
            },
            style: {
              lineWidth: 2,
            },

        };
        return <Line {...config} />;
    } else {
        return <div style={{ textAlign: 'center', display: 'flex', justifyContent: 'center', height: '100%' }}><FormattedMessage id='pages.no.data'  /></div>;
    }

};


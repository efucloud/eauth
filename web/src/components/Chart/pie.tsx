import React from 'react';
import { DashboardData } from '@/services/dashboard.d';
import { FormattedMessage } from '@umijs/max';
import { Pie } from '@ant-design/charts';

interface ChartProps {
    data: DashboardData[],
    height?: number,
    position?: string,
    angleField?: string,
    colorField?: string
}


export const CustomChatPie: React.FC<ChartProps> = (props) => {
    if (props.data.length > 0) {
        const config = {
            data: props.data || [],
            height: props?.height || 100,
            angleField: props.angleField || 'value',
            colorField: props.colorField || 'intlName',
            label: {
                text: props.angleField || 'value',
                style: {
                    fontWeight: 'bold',
                },
            },
            legend: {
                color: {
                    title: false,
                    position: props?.position || 'right',
                    rowPadding: 5,
                },
            },

        };
        return <Pie {...config} />;
    } else {
        return <div style={{ textAlign: 'center', display: 'flex', justifyContent: 'center', height: '100%' }}><FormattedMessage id='pages.no.data'  /></div>;
    }

};


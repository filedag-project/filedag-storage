import dashboardStore from '@/store/modules/dashboard';
import { observer } from 'mobx-react';
import { useEffect } from 'react';
import styles from './style.module.scss';
import * as echarts from 'echarts';
import { formatBytes } from '@/utils';
const Dashboard = (props:any) => {
  useEffect(()=>{
    dashboardStore.fetchRequestOverview().then(res=>{
      initGetBytesChart();
      initPutBytesChart();
      initGetCountChart();
      initPutCountChart();
    })
  },[]);
  const initGetBytesChart = ()=>{
    const el = document.getElementById('getBytesChart');
    if(el){
      var myChart = echarts.init(el);
      myChart.setOption({
        title: {
          text: 'Get object bytes',
          left: 'center'
        },
        tooltip: {
          trigger: 'item',
          formatter: function (params) {
            const _value = formatBytes(params.value??0);
            return params.marker + `${params.name}：${_value}`;
          },
        },
        legend: {
          orient: 'vertical',
          left: 'left'
        },
        series: [
          {
            name: 'Get object bytes',
            type: 'pie',
            radius: '50%',
            data: dashboardStore.get_obj_bytes,
            emphasis: {
              itemStyle: {
                shadowBlur: 10,
                shadowOffsetX: 0,
                shadowColor: 'rgba(0, 0, 0, 0.5)'
              }
            }
          }
        ]
      });
    }
    
  };
  const initPutBytesChart = ()=>{
    const el = document.getElementById('putBytesChart');
    if(el){
      var myChart = echarts.init(el);
      myChart.setOption({
        title: {
          text: 'Put object bytes',
          left: 'center'
        },
        tooltip: {
          trigger: 'item',
          formatter: function (params) {
            const _value = formatBytes(params.value??0);
            return params.marker + `${params.name}：${_value}`;
          },
        },
        legend: {
          orient: 'vertical',
          left: 'left'
        },
        series: [
          {
            name: 'Put object bytes',
            type: 'pie',
            radius: '50%',
            data: dashboardStore.put_obj_bytes,
            emphasis: {
              itemStyle: {
                shadowBlur: 10,
                shadowOffsetX: 0,
                shadowColor: 'rgba(0, 0, 0, 0.5)'
              }
            }
          }
        ]
      });
    }
    
  };
  const initGetCountChart = ()=>{
    const el = document.getElementById('getCountChart');
    if(el){
      var myChart = echarts.init(el);
      myChart.setOption({
        title: {
          text: 'Get object count',
          left: 'center'
        },
        tooltip: {
          trigger: 'item',
          formatter: function (params) {
            return params.marker + `${params.name}：${params.value}`;
          },
        },
        legend: {
          orient: 'vertical',
          left: 'left'
        },
        series: [
          {
            name: 'Get object count',
            type: 'pie',
            radius: '50%',
            data: dashboardStore.get_obj_count,
            emphasis: {
              itemStyle: {
                shadowBlur: 10,
                shadowOffsetX: 0,
                shadowColor: 'rgba(0, 0, 0, 0.5)'
              }
            }
          }
        ]
      });
    }
    
  };
  const initPutCountChart = ()=>{
    const el = document.getElementById('putCountChart');
    if(el){
      var myChart = echarts.init(el);
      myChart.setOption({
        title: {
          text: 'Put object count',
          left: 'center'
        },
        tooltip: {
          trigger: 'item',
          formatter: function (params) {
            return params.marker + `${params.name}：${params.value}`;
          },
        },
        legend: {
          orient: 'vertical',
          left: 'left'
        },
        series: [
          {
            name: 'Put object count',
            type: 'pie',
            radius: '50%',
            data: dashboardStore.put_obj_count,
            emphasis: {
              itemStyle: {
                shadowBlur: 10,
                shadowOffsetX: 0,
                shadowColor: 'rgba(0, 0, 0, 0.5)'
              }
            }
          }
        ]
      });
    }
    
  };
  
  return <div className={styles.dashboard}>
    <div className={styles.chartWrap}>
      <div className={styles.chart} id="getBytesChart"></div>
      <div className={styles.chart} id="putBytesChart"></div>
    </div>
    <div className={styles.chartWrap}>
      <div className={styles.chart} id="getCountChart"></div>
      <div className={styles.chart} id="putCountChart"></div>
    </div>
  </div>;
};

export default observer(Dashboard);

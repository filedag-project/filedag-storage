import dashboardStore from '@/store/modules/dashboard';
import { observer } from 'mobx-react';
import { useEffect } from 'react';
import styles from './style.module.scss';
import * as echarts from 'echarts';
import { formatBytes } from '@/utils';
const Dashboard = (props:any) => {
  useEffect(()=>{
    dashboardStore.fetchStorePool();
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
          text: 'Get object bytes(Top 20)',
          left: 20,
          top:20
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
          right: 20,
          bottom:20,
          formatter: function (name) {
            const obj = dashboardStore.top_20_get_obj_bytes.find(n=>n.name === name);
            const _value = formatBytes(obj?.value??0);
            return `${name}: ${_value}` ;
          }
        },
        series: [
          {
            name: 'Get object bytes',
            type: 'pie',
            radius: 100,
            center:[200,'55%'],
            data: dashboardStore.top_20_get_obj_bytes,
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
          text: 'Put object bytes(Top 20)',
          left: 20,
          top:20
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
          right: 20,
          bottom:20,
          formatter: function (name) {
            const obj = dashboardStore.top_20_put_obj_bytes.find(n=>n.name === name);
            const _value = formatBytes(obj?.value??0);
            return `${name}: ${_value}` ;
          }
        },
        series: [
          {
            name: 'Put object bytes',
            type: 'pie',
            radius: 100,
            center:[200,'55%'],
            data: dashboardStore.top_20_put_obj_bytes,
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
          text: 'Get object count(Top 20)',
          left: 20,
          top:20
        },
        tooltip: {
          trigger: 'item',
          formatter: function (params) {
            return params.marker + `${params.name}：${params.value}`;
          },
        },
        legend: {
          orient: 'vertical',
          right: 20,
          bottom:20,
          formatter: function (name) {
            const obj = dashboardStore.top_20_get_obj_count.find(n=>n.name === name);
            const _value = formatBytes(obj?.value??0);
            return `${name}: ${_value}` ;
          }
        },
        series: [
          {
            name: 'Get object count',
            type: 'pie',
            radius: 100,
            center:[200,'55%'],
            data: dashboardStore.top_20_get_obj_count,
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
          text: 'Put object count(Top 20)',
          left: 20,
          top:20
        },
        tooltip: {
          trigger: 'item',
          formatter: function (params) {
            return params.marker + `${params.name}：${params.value}`;
          },
        },
        legend: {
          orient: 'vertical',
          right: 20,
          bottom:20,
          formatter: function (name) {
            const obj = dashboardStore.top_20_put_obj_count.find(n=>n.name === name);
            const _value = formatBytes(obj?.value??0);
            return `${name}: ${_value}` ;
          }
        },
        series: [
          {
            name: 'Put object count',
            type: 'pie',
            radius: 100,
            center:[200,'55%'],
            data: dashboardStore.top_20_put_obj_count,
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
    <div className={styles.boxWrap}>
      <div className={styles.box}>
        <div className={styles.left}>
          <div className={styles.valueGroup}>
            <span className={styles.value}>
              {dashboardStore.bucketsCount}
            </span>
            <span className={styles.unit}></span>
          </div>
          <div className={styles.label}>Buckets</div>
        </div>
        <div className={styles.right}>
          <img src={require('@/assets/images/dashboard/icon-1.png')} alt="" />
        </div>
      </div>
      <div className={styles.box}>
        <div className={styles.left}>
          <div className={styles.valueGroup}>
            <span className={styles.value}>
                {dashboardStore.objectsCount}
            </span>
            <span className={styles.unit}></span>
          </div>
          <div className={styles.label}>Objects</div>
          
        </div>
        <div className={styles.right}>
          <img src={require('@/assets/images/dashboard/icon-2.png')} alt="" />
        </div>
      </div>
      <div className={styles.box}>
        <div className={styles.left}>
          <div className={styles.valueGroup}>
            <span className={styles.value}>
              {formatBytes(dashboardStore.totalCaptivity)}
            </span>
            <span className={styles.unit}></span>
          </div>
          <div className={styles.label}>Total Storage</div>
          
        </div>
        <div className={styles.right}>
          <img src={require('@/assets/images/dashboard/icon-3.png')} alt="" />
        </div>
      </div>
      <div className={styles.box}>
        <div className={styles.left}>
          <div  className={styles.valueGroup}>
            <span className={styles.value}>
              {formatBytes(dashboardStore.objectsTotalSize)}
            </span>
            <span className={styles.unit}></span>
          </div>
          <div className={styles.label}>Use Storage</div>
        </div>
        <div className={styles.right}>
          <img src={require('@/assets/images/dashboard/icon-4.png')} alt="" />
        </div>
      </div>

    </div>
    <div className={styles.chartWrap}>
      <div className={styles.chart} id="getBytesChart"></div>
      <div className={styles.chart} id="putBytesChart"></div>
      <div className={styles.chart} id="getCountChart"></div>
      <div className={styles.chart} id="putCountChart"></div>
    </div>
  </div>;
};

export default observer(Dashboard);
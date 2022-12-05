import { Axios, HttpMethods } from '@/api/https';
import { SignModel } from '@/models/SignModel';
import dashboardStore from '@/store/modules/dashboard';
import globalStore from '@/store/modules/global';
import { SaveOutlined } from '@ant-design/icons';
import { observer } from 'mobx-react';
import { useEffect } from 'react';
import styles from './style.module.scss';
const Dashboard = (props:any) => {
  useEffect(()=>{
    globalStore.fetchIsAdmin();
    dashboardStore.fetchRequestOverview();
    dashboardStore.fetchUserInfos();
  },[]);
  
  return <div className={styles.dashboard}>
    <div className={styles.boxWrap}>
      <div className={styles.box}>
        <div className={styles.top}>
          <span className={styles.label}>Objects</span>
          <SaveOutlined />
        </div>
        <div className={styles.bottom}>
          <span className={styles.value}>
            {dashboardStore.overview.get_obj_bytes}
          </span>
          <span className={styles.unit}></span>
        </div>
      </div>
    </div>
  </div>;
};

export default observer(Dashboard);

import { observer } from 'mobx-react';
import { RestOutlined,SaveOutlined,FolderOutlined,FileTextOutlined } from '@ant-design/icons';
import dashboardStore from "@/store/modules/dashboard";
import styles from './style.module.scss';
import { useEffect } from 'react';

const Dashboard = (props:any) => {
  useEffect(()=>{
    dashboardStore.fetchUserInfo();
  },[])
  return <div className={styles.dashboard}>
    <div className={styles.box}>
      <div className={styles.top}>
        <span className={styles.label}>Buckets</span>
        <RestOutlined />
      </div>
      <div className={styles.bottom}>
        <span className={styles.value}>
          {dashboardStore.userInfo.buckets}
        </span>
        <span className={styles.unit}></span>
      </div>
    </div>
    <div className={styles.box}>
      <div className={styles.top}>
        <span className={styles.label}>Objects</span>
        <SaveOutlined />
      </div>
      <div className={styles.bottom}>
        <span className={styles.value}>
          {dashboardStore.userInfo.objects}
        </span>
        <span className={styles.unit}></span>
      </div>
    </div>
    <div className={styles.box}>
      <div className={styles.top}>
        <span className={styles.label}>Total Storage</span>
        <FolderOutlined />
      </div>
      <div className={styles.bottom}>
        <span className={styles.value}>
          {dashboardStore.userInfo.total_storage_capacity}
        </span>
        <span className={styles.unit}></span>
      </div>
    </div>
    <div className={styles.box}>
      <div className={styles.top}>
        <span className={styles.label}>Use Storage</span>
        <FileTextOutlined />
      </div>
      <div className={styles.bottom}>
        <span className={styles.value}>
          {dashboardStore.userInfo.use_storage_capacity}
        </span>
        <span className={styles.unit}></span>
      </div>
    </div>
  </div>;
};

export default observer(Dashboard);

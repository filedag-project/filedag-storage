import { observer } from 'mobx-react';
import { RestOutlined,SaveOutlined,FolderOutlined,FileTextOutlined } from '@ant-design/icons';
import styles from './style.module.scss';
import { useEffect } from 'react';
import overviewStore from '@/store/modules/overview';


const Overview = (props:any) => {
  useEffect(()=>{
    overviewStore.fetchUserInfo();
  },[]);
  return <div className={styles.overview}>
    <div className={styles.box}>
      <div className={styles.top}>
        <span className={styles.label}>Buckets</span>
        <RestOutlined />
      </div>
      <div className={styles.bottom}>
        <span className={styles.value}>
          {overviewStore.userInfo.buckets}
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
          {overviewStore.userInfo.objects}
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
          {overviewStore.userInfo.total_storage_capacity}
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
          {overviewStore.userInfo.use_storage_capacity}
        </span>
        <span className={styles.unit}></span>
      </div>
    </div>
  </div>;
};

export default observer(Overview);

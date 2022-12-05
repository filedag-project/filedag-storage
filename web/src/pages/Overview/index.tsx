import { observer } from 'mobx-react';
import { RestOutlined,SaveOutlined,FolderOutlined,FileTextOutlined } from '@ant-design/icons';
import styles from './style.module.scss';
import globalStore from '@/store/modules/global';

const Overview = (props:any) => {
  return <div className={styles.overview}>
    <div className={styles.box}>
      <div className={styles.top}>
        <span className={styles.label}>Buckets</span>
        <RestOutlined />
      </div>
      <div className={styles.bottom}>
        <span className={styles.value}>
          {globalStore.userInfo.buckets}
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
          {globalStore.userInfo.objects}
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
          {globalStore.userInfo.total_storage_capacity}
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
          {globalStore.userInfo.use_storage_capacity}
        </span>
        <span className={styles.unit}></span>
      </div>
    </div>
  </div>;
};

export default observer(Overview);

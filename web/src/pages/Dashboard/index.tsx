import { observer } from 'mobx-react';
import styles from './style.module.scss';
const Dashboard = (props:any) => {
  return <div className={styles.dashboard}>
   Dashboard
  </div>;
};

export default observer(Dashboard);

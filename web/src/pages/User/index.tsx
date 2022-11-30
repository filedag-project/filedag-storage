import { observer } from "mobx-react";
import styles from './style.module.scss';
const User = (props:any) => {
  
  return <div className={styles.user}>
    User
  </div>;
};

export default observer(User);

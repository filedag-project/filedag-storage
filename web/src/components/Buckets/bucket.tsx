import styles from './bucket.module.scss';
import { Button } from 'antd';
import bucketsStore from '@/store/modules/buckets';
import { RouterPath } from '@/router/RouterConfig';
import { useNavigate } from 'react-router-dom';
import { observer } from 'mobx-react';
import iconSetting from '@/assets/images/common/icon-setting.png';
import iconView from '@/assets/images/common/icon-view.png';
import iconDelete from '@/assets/images/common/icon-delete.png';
const Bucket = (props:any) => {
  const navigate = useNavigate();
  const { data :{ Name: bucket,CreationDate } } = props;
  const openDelete = ()=>{
    bucketsStore.SET_DELETE_NAME(bucket);
    bucketsStore.SET_DELETE_SHOW(true);
  }

  const viewObjects = ()=>{
    navigate(`${RouterPath.bucketDetail}`,{
      state:{ bucket }
    });
  }
  const openPower = ()=>{
    navigate(`${RouterPath.power}`,{ state :{ bucket }});
  }
  return (
    <>
      <div className={styles['bucket']}>
        <div className={styles['top']}>
          <div className={styles['info']}>
            <div className={styles['name']}>{ bucket }</div>
            <div className={styles['create']}>Created: {CreationDate}</div>
          </div>
          <div className={styles['action']}>
            <div onClick={openPower}>
              <Button className={'setting-btn'} type="primary" icon={<img src={iconSetting} alt='' />}>Setting</Button>
            </div>
            <div onClick={viewObjects}>
              <Button className={'browse-btn'} type="primary" icon={<img src={iconView} alt='' />}>View</Button>
            </div>
            <div onClick={openDelete}>
              <Button className={'delete-btn'} type="primary" icon={<img src={iconDelete} alt='' />}>Delete</Button>
            </div>
          </div>
        </div>
        <div className={styles['bottom']}></div>
      </div>
      
    </>
  );
};

export default observer(Bucket);

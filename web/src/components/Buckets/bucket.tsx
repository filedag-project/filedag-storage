import styles from './bucket.module.scss';
import { FolderViewOutlined,DeleteOutlined } from '@ant-design/icons';
import { Button } from 'antd';
import bucketsStore from '@/store/modules/buckets';
import { RouterPath } from '@/router/RouterConfig';
import { useHistory } from 'react-router';
import { observer } from 'mobx-react';
const Bucket = (props:any) => {
  const history = useHistory();
  const { data } = props;
  const openDelete = ()=>{
    const {Name} = data;
    bucketsStore.SET_DELETE_NAME(Name);
    bucketsStore.SET_DELETE_SHOW(true);
  }

  const viewObjects = ()=>{
    const {Name,CreationDate} = data;
    const path = `/${Name}`
    history.push({
      pathname: RouterPath.objects,
      state: { bucket:  path,Created:CreationDate},
    });
  }
  return (
    <>
      <div className={styles['bucket']}>
        <div className={styles['top']}>
          <div className={styles['info']}>
            <div className={styles['name']}>{data.Name}</div>
            <div className={styles['create']}>Created: {data.CreationDate}</div>
          </div>
          <div className={styles['action']}>
            <div className={styles['browse']} onClick={viewObjects}>
              <Button type="primary" icon={<FolderViewOutlined />}>查看</Button>
            </div>
            <div className={styles['delete']} onClick={openDelete}>
              <Button type="primary" icon={<DeleteOutlined />}>删除</Button>
            </div>
          </div>
        </div>
        <div className={styles['bottom']}></div>
      </div>
      
    </>
  );
};

export default observer(Bucket);

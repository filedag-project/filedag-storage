import styles from './action.module.scss';
import { Upload, Button } from 'antd';
import { UploadOutlined } from '@ant-design/icons';
import objectsStore from '@/store/modules/objects';
import { escapeStr } from '@/utils';
import { observer } from 'mobx-react';

const Action = (props:any) => {
  const {bucket,Created} = props;
  const customChange = async (e)=>{
    const name = escapeStr(e.file.name);
    const path = `${bucket}/${name}`;
    await objectsStore.fetchUpload(path,e.file);
    objectsStore.fetchList(bucket);
  }
  return (
    <div className={styles.action}>
      <div className={styles.info}>
        <div className={styles.bucket}>Books</div>
        <div className={styles['date-size']}>
          <span>Created:{Created}</span>
          <span>Size:{objectsStore.totalSize}</span>
          <span>
            {objectsStore.totalObjects} {' '}
            object
            {objectsStore.totalObjects>1?'s':''}
          </span>
        </div>
      </div>
      <div className={styles.operation}>
        <Upload customRequest={customChange} showUploadList={false}>
          <Button type="primary" icon={<UploadOutlined />}>Upload</Button>
        </Upload>
      </div>
    </div>
  );
};

export default observer(Action);

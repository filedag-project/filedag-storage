import Action from '@/components/Buckets/action';
import Bucket from '@/components/Buckets/bucket';
import bucketsStore from '@/store/modules/buckets';
import { Modal } from 'antd';
import { observer } from 'mobx-react';
import { useEffect } from 'react';
import styles from './style.module.scss';

const Buckets = () => {
  useEffect(()=>{
    bucketsStore.fetchList();
  },[])

  const confirmDelete = ()=>{
    const path = `/${bucketsStore.deleteName}`;
    bucketsStore.fetchDelete(path).then(res=>{
      bucketsStore.fetchList();
      bucketsStore.SET_DELETE_SHOW(false);
    });
  }

  const cancelDelete = ()=>{
    bucketsStore.SET_DELETE_SHOW(false);
  }

  return (
    <div className={styles.buckets}>
      <Action></Action>
      {
        bucketsStore.formatList.map((item,index)=>{
          return <Bucket key={index} data={item}></Bucket>
        })
      }

      <Modal
        title="Delete"
        open={bucketsStore.deleteShow}
        onOk={confirmDelete}
        onCancel={cancelDelete}
        okText="Confirm"
        cancelText="Cancel"
      >
        <p>Are you sure to delete this data？</p>
      </Modal>
    </div>
  );
};

export default observer(Buckets);

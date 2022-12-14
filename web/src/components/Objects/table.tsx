import objectsStore from '@/store/modules/objects';
import { Button, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { observer } from 'mobx-react';
import iconShare from '@/assets/images/common/icon-view.png';
import iconView from '@/assets/images/common/icon-view.png';
import iconDelete from '@/assets/images/common/icon-delete.png';
import styles from './table.module.scss'

interface DataType {
  Name: string;
  LastModified: string;
  Size: string;
}
const ObjectTable = (props:any) => {
  const columns: ColumnsType<DataType> = [
    {
      title: 'Name',
      dataIndex: 'Name',
      key: 'name',
      width:150,
      render:(r)=>{
        return <div className='row-name'>{r}</div>
      }
    },
    {
      title: 'LastModified',
      dataIndex: 'LastModified',
      key: 'LastModified',
    },
    {
      title: 'Size',
      dataIndex: 'Size',
      key: 'Size',
    },
    {
      title:'ETag',
      dataIndex: 'ETag',
      key: 'ETag',
    },
    {
      title: 'Action',
      key: 'action',
      render: (_, record) => (
        <div className='row-action'>
          <Button className='share-btn' type="primary" icon={<img src={iconShare} alt='' />} onClick={()=>{shareObject(record.Name)}}>Share</Button>
          <Button className='view-btn' type="primary" icon={<img src={iconView} alt='' />} onClick={()=>{viewObject(record.Name)}}>View</Button>
          <Button className='delete-btn' type="primary" icon={<img src={iconDelete} alt='' />} onClick={()=>{deleteObject(record.Name)}}>Delete</Button>
        </div>
      ),
    },
  ];
  const deleteObject = (name:string)=>{
    objectsStore.SET_DELETE_SHOW(true);
    objectsStore.SET_ACTION_NAME(name);
  };

  const viewObject = (name:string)=>{
    objectsStore.SET_ACTION_NAME(name);
    objectsStore.SET_PREVIEW_SHOW(true);
    objectsStore.fetchObject(`/${props.bucket}`,name);
  }

  const shareObject = async (name:string)=>{
    objectsStore.SET_ACTION_NAME(name);
    objectsStore.SET_SHARE_SHOW(true);
    const url = `/${props.bucket}/${name}`;
    objectsStore.fetchShare(url,objectsStore.shareSecond);
  }
  
  return <div className={styles['object-table']}>
    <Table columns={columns} dataSource={objectsStore.formatList}  rowKey={(record) => record.Name} pagination={false}/>
  </div>;
};

export default observer(ObjectTable);

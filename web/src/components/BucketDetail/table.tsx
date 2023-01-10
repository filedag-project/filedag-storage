import bucketDetailStore from '@/store/modules/bucketDetail';
import { Button, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { observer } from 'mobx-react';
import iconShare from '@/assets/images/common/icon-view.png';
import iconView from '@/assets/images/common/icon-view.png';
import iconDelete from '@/assets/images/common/icon-delete.png';
import styles from './table.module.scss'
import { FileType } from '@/models/BucketModel';
import { useNavigate } from 'react-router-dom';
import { RouterPath } from '@/router/RouterConfig';

interface DataType {
  Name: string;
  LastModified: string;
  Size: string;
}
const ObjectTable = (props:any) => {
  const navigate = useNavigate();
  const { bucket,prefix } = props; 
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
      render: (_, record) => {
        return (
          <div className='row-action'>
            {
              record['Type'] === FileType.file?<Button className='share-btn' type="primary" icon={<img src={iconShare} alt='' />} onClick={()=>{shareObject(record['Name'])}}>Share</Button>:<></>
            }
            <Button className='view-btn' type="primary" icon={<img src={iconView} alt='' />} onClick={()=>{viewObject(record['Name'],record['Type'])}}>View</Button>
            <Button className='delete-btn' type="primary" icon={<img src={iconDelete} alt='' />} onClick={()=>{deleteObject(record['Name'])}}>Delete</Button>
          </div>
        )
      },
    },
  ];
  const deleteObject = (name:string)=>{
    bucketDetailStore.SET_DELETE_SHOW(true);
    bucketDetailStore.SET_ACTION_NAME(`${prefix}${name}`);
  };

  const viewObject = (name:string,type:string)=>{
    if(type === FileType.file){
      bucketDetailStore.fetchObject(`/${props.bucket}`,name).then(res=>{
        bucketDetailStore.SET_ACTION_NAME(name);
        bucketDetailStore.SET_PREVIEW_SHOW(true);
      })
    }else{
      navigate(`${RouterPath.bucketDetail}`,{ state :{ bucket,prefix:`${prefix}${name}`}})
    }
  }

  const shareObject = async (name:string)=>{
    bucketDetailStore.SET_ACTION_NAME(name);
    bucketDetailStore.SET_SHARE_SHOW(true);
    const url = `/${props.bucket}/${name}`;
    bucketDetailStore.fetchShare(url,bucketDetailStore.shareSecond);
  }
  
  return <div className={styles['object-table']}>
    <Table columns={columns} dataSource={bucketDetailStore.formatList}  rowKey={(record) => record.Name} pagination={false}/>
  </div>;
};

export default observer(ObjectTable);

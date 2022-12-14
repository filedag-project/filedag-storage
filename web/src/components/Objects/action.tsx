import styles from './action.module.scss';
import { Upload, Button, Progress } from 'antd';
import { CloseCircleOutlined } from '@ant-design/icons';
import objectsStore from '@/store/modules/objects';
import { observer } from 'mobx-react';
import _ from 'lodash';
import { useState } from 'react';
import { pieceBytes } from '@/config';
import iconUpload from '@/assets/images/common/icon-upload.png';

const Action = (props:any) => {
  const [progressShow,setProgressShow] = useState(false);
  const {bucket,Created} = props;
  const customChange = async (e)=>{
    const name = e.file.name;
    const path = `/${bucket}/${name}`;
    const _size = e.file.size;
    console.log(0);
    
    if(_size > pieceBytes){
      sliceUpload(path,e.file);
    }else{
      commonUpload(path,e.file);
    }
  }
  const sliceUpload = (path:string,file:File)=>{
    const totalPieces = Math.ceil(file.size / pieceBytes);
    objectsStore.fetchUploadId(path).then(async res=>{
      const uploadId:string = typeof res === 'string' ? res :'';
      let index = 0;
      let parts = ``;
      while(index < totalPieces){
        try{
          const end = (index+1) * pieceBytes;
          const _file = file.slice(index,end);
          setProgressShow(true);
          const slice = await objectsStore.fetchSliceUpload(path,index,uploadId,_file);
          const _progress= ((index+1) / totalPieces * 100 | 0)
          objectsStore.SET_PERCENTAGE(_progress);
          const etag = _.get(slice,'etag','');
          parts += `
            <Part>
              <ETag>${etag}</ETag>
              <PartNumber>${index}</PartNumber>
            </Part>
          `;
          index++;
        }catch(error){
          objectsStore.fetchAbort(path,uploadId)
          break;
        }
      }
      if(index === totalPieces){
        const body = `
          <?xml version="1.0" encoding="UTF-8"?>
          <CompleteMultipartUpload xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
            ${parts}
          </CompleteMultipartUpload> 
        `;
        objectsStore.fetchSliceUploadComplete(path,uploadId,body).then(result=>{
          console.log('fetchSliceUploadComplete');
          objectsStore.fetchList(bucket);
        })
      }
    });
  }
  const commonUpload = async (path:string,file:File)=>{
    setProgressShow(true);
    await objectsStore.fetchUpload(path,file);
    console.log(1);
    
    objectsStore.fetchList(bucket);
  }
  return (
    <div className={styles.action}>
      <div className={styles.info}>
        <div className={styles.bucket}>{bucket}</div>
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
          <Button className='bg-btn' type="primary" icon={<img src={iconUpload} alt='' />}>Upload</Button>
        </Upload>
        <div className={styles.progressWrap} style={{display:progressShow?'block':'none'}}>
          <div className={styles.title}>
            <span>Upload Progress</span>
            <CloseCircleOutlined onClick={()=> {
              setProgressShow(false);
              objectsStore.SET_PERCENTAGE(0);
            }}/>
          </div>
          <Progress percent={objectsStore.percentage}></Progress>
        </div>
      </div>
    </div>
  );
};

export default observer(Action);

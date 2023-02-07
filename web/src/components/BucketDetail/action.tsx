import styles from './action.module.scss';
import { Upload, Button, Progress } from 'antd';
import { CloseCircleOutlined, FolderOutlined } from '@ant-design/icons';
import bucketDetailStore from '@/store/modules/bucketDetail';
import { observer } from 'mobx-react';
import _ from 'lodash';
import { useState } from 'react';
import { PIECE_BYTES } from '@/config';
import iconUpload from '@/assets/images/common/icon-upload.png';
import { useNavigate } from 'react-router-dom';
import { RouterPath } from '@/router/RouterConfig';

const Action = (props:any) => {
  const navigate = useNavigate();
  const [progressShow,setProgressShow] = useState(false);
  const {bucket,prefix} = props;
  const prefixArray = (prefix.split('/')).filter(n=>n);
  
  
  const customChange = async (e)=>{
    const name = e.file.name.replace(':','/');
    const _prefix = prefix ? `${prefix}`:''
    const path = `/${bucket}/${_prefix}${name}`;
    const _size = e.file.size;
    setProgressShow(false);
    bucketDetailStore.SET_PERCENTAGE(0);
    if(_size > PIECE_BYTES){
      sliceUpload(path,e.file);
    }else{
      commonUpload(path,e.file);
    }
  }
  const sliceUpload = (path:string,file:File)=>{
    const totalPieces = Math.ceil(file.size / PIECE_BYTES);
    bucketDetailStore.fetchUploadId(path).then(async res=>{
      const uploadId:string = typeof res === 'string' ? res :'';
      let index = 0;
      let parts = ``;
      while(index < totalPieces){
        try{
          const end = (index+1) * PIECE_BYTES;
          const _file = file.slice(index,end);
          setProgressShow(true);
          const slice = await bucketDetailStore.fetchSliceUpload(path,index,uploadId,_file);
          const _progress= ((index+1) / totalPieces * 100 | 0)
          bucketDetailStore.SET_PERCENTAGE(_progress);
          const etag = _.get(slice,'etag','');
          parts += `
            <Part>
              <ETag>${etag}</ETag>
              <PartNumber>${index}</PartNumber>
            </Part>
          `;
          index++;
        }catch(error){
          bucketDetailStore.fetchAbort(path,uploadId)
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
        bucketDetailStore.fetchSliceUploadComplete(path,uploadId,body).then(async result=>{
          resetList();
          await bucketDetailStore.fetchList(bucket,prefix);
          closeProgress();
        })
      }
    });
  }

  const resetList = ()=>{
    bucketDetailStore.SET_CURRENT_PAGE(1);
    bucketDetailStore.SET_NEXT_CONTINUE_TOKEN('');
    bucketDetailStore.SET_CONTENTS_LIST([]);
    bucketDetailStore.SET_COMMON_PREFIXES_LIST([]);
  }

  const closeProgress = ()=>{
    setTimeout(()=>{
      setProgressShow(false);
      bucketDetailStore.SET_PERCENTAGE(0);
    },2000)
  }
  const commonUpload = async (path:string,file:File)=>{
    setProgressShow(true);
    await bucketDetailStore.fetchUpload(path,file);
    resetList();
    await bucketDetailStore.fetchList(bucket,prefix);
    closeProgress();
  }
  const objectClick = (object:string,index:number)=>{
    const str = prefixArray.slice(0,index+1);
    const _p_ = str.join('/') + `/`;
    console.log(_p_,index,'_p_');
    navigate(RouterPath.bucketDetail,{
      state :{ bucket, prefix: _p_}
    })
  }

  return (
    <div className={styles.action}>
      <div className={styles.info}>
        <div className={styles.bucket}>
          <span className={styles['object-name']} onClick={()=>{
            navigate(RouterPath.bucketDetail,{
              state :{ bucket }
            })
          }}>
            {bucket}
          </span>
          {
            prefixArray.map((n,index)=>{
              const to = ` > `;
              return <span className={styles['object-name']} key={'prefix'+index} onClick={()=>{ objectClick(n,index) }}>{to}{n}</span>
            })
          }
        </div>
        <div className={styles['date-size']}>
          <span>Object:{bucketDetailStore.keyCount}</span>
        </div>
      </div>
      <div className={styles.operation}>
        <Button className='upload-folder' type="primary" icon={<FolderOutlined />} onClick={()=>{
          bucketDetailStore.SET_ADD_FOLDER_SHOW(true);
        }}>
          Upload Folder
        </Button>
        <Upload customRequest={customChange} showUploadList={false}>
          <Button className='bg-btn' type="primary" icon={<img src={iconUpload} alt='' />}>Upload File</Button>
        </Upload>
        <div className={styles.progressWrap} style={{display:progressShow?'block':'none'}}>
          <div className={styles.title}>
            <span>Upload Progress</span>
            <CloseCircleOutlined onClick={()=> {
              setProgressShow(false);
              bucketDetailStore.SET_PERCENTAGE(0);
            }}/>
          </div>
          <Progress percent={bucketDetailStore.percentage}></Progress>
        </div>
      </div>
    </div>
  );
};

export default observer(Action);

import Action from "@/components/Objects/action";
import ObjectTable from "@/components/Objects/table";
import objectsStore from "@/store/modules/objects";
import { download, getExpiresDate } from "@/utils";
import { Modal, InputNumber, notification } from "antd";
import { observer } from "mobx-react";
import { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import ReactPlayer from 'react-player';
import styles from './style.module.scss';
import { CopyOutlined } from "@ant-design/icons";
import copy from 'copy-to-clipboard';
const Objects = (props:any) => {
  const [day,setDay]=useState(7);
  const [hour,setHour]=useState(0);
  const [minute,setMinute]=useState(0);
  const { state :{ bucket = '',create = '' }} = useLocation();
  
  useEffect(()=>{
    objectsStore.fetchList(bucket);
  },[bucket]);

  const confirmDelete = async ()=>{
    const name = objectsStore.actionName;
    const path = `/${bucket}/${name}`;
    await objectsStore.fetchDelete(path);
    objectsStore.SET_DELETE_SHOW(false);
    objectsStore.fetchList(bucket);
  };
  const cancelDelete = ()=>{
    objectsStore.SET_DELETE_SHOW(false);
  };

  const downloadPreview = ()=>{
    download(objectsStore.downloadFile,objectsStore.actionName);
  };

  const cancelPreview = ()=>{
    objectsStore.SET_PREVIEW_SHOW(false);
  };

  const previewDom = ()=>{
    if(objectsStore.contentType.includes('image')){
      return <img src={objectsStore.previewUrl} className="preview-image" alt="" />
    }else if(objectsStore.contentType.includes('text')){
      return <span className="preview-text">{objectsStore.previewText}</span>
    }else if(objectsStore.contentType.includes('video')){
      return <ReactPlayer
        width="400px"
        height="200px"
        url={objectsStore.previewVideo}
        loop
        controls
      ></ReactPlayer>
    }else{
      return <div>The current file cannot be previewed, please download it locally to view.</div>;
    }
  }

  const inputLimit = (value):string=>{
    const reg = /^\d+$/;
    if(!reg.test(value)){
      return '';
    }
    return value;
  }

  const dayChange = (value)=>{
    const _day = value??0;
    setDay(value??'');
    const s = _day*24*60*60;
    const _expiresDate = getExpiresDate(s);
    objectsStore.SET_EXPIRES_DATE(_expiresDate);
    objectsStore.SET_SHARE_SECOND(s);
    const url = `/${bucket}/${objectsStore.actionName}`;
    objectsStore.fetchShare(url,objectsStore.shareSecond);
  }
  const hourChange = (value)=>{
    const _hour = value??0;
    setHour(value??'');
    const s = _hour*60*60;
    const _expiresDate = getExpiresDate(s);
    objectsStore.SET_EXPIRES_DATE(_expiresDate);
    objectsStore.SET_SHARE_SECOND(s);
    const url = `/${bucket}/${objectsStore.actionName}`;
    objectsStore.fetchShare(url,objectsStore.shareSecond);
  }
  const minuteChange = (value)=>{
    const _minute = value??0;
    setMinute(value??'');
    const s = _minute*60;
    const _expiresDate = getExpiresDate(s);
    objectsStore.SET_EXPIRES_DATE(_expiresDate);
    objectsStore.SET_SHARE_SECOND(s);
    const url = `/${bucket}/${objectsStore.actionName}`;
    objectsStore.fetchShare(url,objectsStore.shareSecond);
  }


  return <div className={styles.objects}>
    <Action bucket={bucket} Created={create}></Action>
    <ObjectTable bucket={bucket}></ObjectTable>
    <Modal
      title="Delete"
      open={objectsStore.deleteShow}
      onOk={confirmDelete}
      onCancel={cancelDelete}
      okText="Confirm"
      cancelText="Cancel"
    >
        <p>Are you sure to delete this data？</p>
    </Modal>
    <Modal 
      title="Preview"
      open={objectsStore.previewShow}
      onOk={downloadPreview}
      onCancel={cancelPreview}
      okText="Download"
      cancelText="Cancel"
    >
      <div className="modal-content">
        {previewDom()}
      </div>
    </Modal>
    <Modal
      title="Share File"
      open={objectsStore.shareShow}
      onCancel={()=>{objectsStore.SET_SHARE_SHOW(false)}}
      footer={<></>}
    >
        <div className="share-description">This is a temporary URL with integrated access credentials for sharing objects valid for up to 7 days.The temporary URL expires after the configured time limit.</div>
        <div className="share-title">Active for</div>
        <div className="share-input">
          <span>Day</span>
          <InputNumber onChange={dayChange} formatter={inputLimit} min={0} max={7} value={day}></InputNumber>
          <span>Hours</span>
          <InputNumber onChange={hourChange} formatter={inputLimit} min={0} max={23} value={hour}></InputNumber>
          <span>Minutes</span>
          <InputNumber onChange={minuteChange} formatter={inputLimit} min={0} max={59} value={minute}></InputNumber>
        </div>
        <div className="share-subTitle">Link will be available until：{objectsStore.expiresDate}</div>
        <div className="share-link">
          <div className="text">{objectsStore.shareLink}</div>
          <div className="copy" onClick={()=>{
            copy(objectsStore.shareLink);
            notification.open({
              message: 'Copy success',
              description: 'Share URL Copied to clipboard',
            });
          }}><CopyOutlined /></div>
        </div>
    </Modal>
  </div>;
};

export default observer(Objects);

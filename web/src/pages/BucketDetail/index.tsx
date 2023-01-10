import Action from "@/components/BucketDetail/action";
import ObjectTable from "@/components/BucketDetail/table";
import bucketDetailStore from "@/store/modules/bucketDetail";
import { download, getExpiresDate } from "@/utils";
import { Modal, InputNumber, notification, Form, Input } from "antd";
import { observer } from "mobx-react";
import { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import ReactPlayer from 'react-player';
import styles from './style.module.scss';
import { CopyOutlined } from "@ant-design/icons";
import copy from 'copy-to-clipboard';

const BucketDetail = (props:any) => {
  const [addFolderForm] = Form.useForm();
  const [day,setDay]=useState(7);
  const [hour,setHour]=useState(0);
  const [minute,setMinute]=useState(0);
  const { state :{ bucket = '', prefix='' }} = useLocation();
  
  useEffect(()=>{
    bucketDetailStore.fetchList(bucket,prefix);
  },[bucket,prefix]);

  const confirmDelete = async ()=>{
    const name = bucketDetailStore.actionName;
    const path = `/${bucket}/${name}`;
    await bucketDetailStore.fetchDelete(path);
    bucketDetailStore.SET_DELETE_SHOW(false);
    bucketDetailStore.fetchList(bucket,prefix);
  };
  const cancelDelete = ()=>{
    bucketDetailStore.SET_DELETE_SHOW(false);
  };

  const downloadPreview = ()=>{
    download(bucketDetailStore.downloadFile,bucketDetailStore.actionName);
  };

  const cancelPreview = ()=>{
    bucketDetailStore.SET_PREVIEW_SHOW(false);
  };

  const previewDom = ()=>{
    console.log(bucketDetailStore,'bucketDetailStore');
    
    if(bucketDetailStore.contentType.includes('image')){
      return <img src={bucketDetailStore.previewUrl} className="preview-image" alt="" />
    }else if(bucketDetailStore.contentType.includes('text')){
      return <span className="preview-text">{bucketDetailStore.previewText}</span>
    }else if(bucketDetailStore.contentType.includes('video')){
      return <ReactPlayer
        width="400px"
        height="200px"
        url={bucketDetailStore.previewVideo}
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
    bucketDetailStore.SET_EXPIRES_DATE(_expiresDate);
    bucketDetailStore.SET_SHARE_SECOND(s);
    const url = `/${bucket}/${bucketDetailStore.actionName}`;
    bucketDetailStore.fetchShare(url,bucketDetailStore.shareSecond);
  }
  const hourChange = (value)=>{
    const _hour = value??0;
    setHour(value??'');
    const s = _hour*60*60;
    const _expiresDate = getExpiresDate(s);
    bucketDetailStore.SET_EXPIRES_DATE(_expiresDate);
    bucketDetailStore.SET_SHARE_SECOND(s);
    const url = `/${bucket}/${bucketDetailStore.actionName}`;
    bucketDetailStore.fetchShare(url,bucketDetailStore.shareSecond);
  }
  const minuteChange = (value)=>{
    const _minute = value??0;
    setMinute(value??'');
    const s = _minute*60;
    const _expiresDate = getExpiresDate(s);
    bucketDetailStore.SET_EXPIRES_DATE(_expiresDate);
    bucketDetailStore.SET_SHARE_SECOND(s);
    const url = `/${bucket}/${bucketDetailStore.actionName}`;
    bucketDetailStore.fetchShare(url,bucketDetailStore.shareSecond);
  }

  const uploadFolder = async ()=>{
    try{
      await addFolderForm.validateFields();
      const folderName = addFolderForm.getFieldValue('folderName');
      const _prefix = prefix ? `${prefix}`:''
      const path = `/${bucket}/${_prefix}${folderName}/`;
      bucketDetailStore.fetchUploadFolder(path).then(res=>{
        bucketDetailStore.SET_ADD_FOLDER_SHOW(false);
        bucketDetailStore.fetchList(bucket,prefix);
      })
    }catch(error){

    }
  }


  return <div className={styles.objects}>
    <Action bucket={bucket} prefix={prefix}></Action>
    <ObjectTable bucket={bucket} prefix={prefix}></ObjectTable>
    <Modal
      title="Delete"
      open={bucketDetailStore.deleteShow}
      onOk={confirmDelete}
      onCancel={cancelDelete}
      okText="Confirm"
      cancelText="Cancel"
    >
        <p>Are you sure to delete this data？</p>
    </Modal>
    <Modal 
      title="Preview"
      open={bucketDetailStore.previewShow}
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
      title="Create Folder"
      open={bucketDetailStore.addFolderShow}
      onOk={uploadFolder}
      onCancel={()=>{
        bucketDetailStore.SET_ADD_FOLDER_SHOW(false);
        addFolderForm.resetFields();
      }}
      okText="Confirm"
      cancelText="Cancel"
    >
      <div className="modal-content">
      <Form form={addFolderForm} autoComplete="off">
        <Form.Item
            name="folderName"
            rules={[
                {required: true, message: 'Please input folder name.'},
                ({ getFieldValue }) => ({
                  validator(_, value) {
                    if (value && value.includes('/')) {
                      return Promise.reject(new Error('Do not include /'));
                    }else{
                      return Promise.resolve();
                    }
                  },
                }),
            ]}
        >
            <Input
                placeholder="please enter folder name"
            />
        </Form.Item>
      </Form>
      </div>
    </Modal>
    
    <Modal
      title="Share File"
      open={bucketDetailStore.shareShow}
      onCancel={()=>{bucketDetailStore.SET_SHARE_SHOW(false)}}
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
        <div className="share-subTitle">Link will be available until：{bucketDetailStore.expiresDate}</div>
        <div className="share-link">
          <div className="text">{bucketDetailStore.shareLink}</div>
          <div className="copy" onClick={()=>{
            copy(bucketDetailStore.shareLink);
            notification.open({
              message: 'Copy success',
              description: 'Share URL Copied to clipboard',
            });
          }}><CopyOutlined /></div>
        </div>
    </Modal>
  </div>;
};

export default observer(BucketDetail);

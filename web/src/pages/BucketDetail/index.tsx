import Action from "@/components/BucketDetail/action";
import ObjectTable from "@/components/BucketDetail/table";
import bucketDetailStore from "@/store/modules/bucketDetail";
import { download, getExpiresDate } from "@/utils";
import { Modal, InputNumber, notification, Form, Input, Pagination, Divider } from "antd";
import { observer } from "mobx-react";
import { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import ReactPlayer from 'react-player';
import styles from './style.module.scss';
import { CopyOutlined, LeftOutlined, RightOutlined } from "@ant-design/icons";
import copy from 'copy-to-clipboard';
import classNames from "classnames";
import { PAGE_SIZE } from "@/config";

const BucketDetail = (props:any) => {
  const [addFolderForm] = Form.useForm();
  const [day,setDay]=useState(7);
  const [hour,setHour]=useState(0);
  const [minute,setMinute]=useState(0);
  const { state :{ bucket = '', prefix='' }} = useLocation();
  
  useEffect(()=>{
    resetList();
    bucketDetailStore.fetchList(bucket,prefix);
  },[bucket,prefix]);

  const resetList = ()=>{
    bucketDetailStore.SET_CURRENT_PAGE(1);
    bucketDetailStore.SET_NEXT_CONTINUE_TOKEN('');
    bucketDetailStore.SET_CONTENTS_LIST([]);
    bucketDetailStore.SET_COMMON_PREFIXES_LIST([]);
  }

  const confirmDelete = async ()=>{
    const name = bucketDetailStore.actionName;
    const path = `/${bucket}/${name}`;
    await bucketDetailStore.fetchDelete(path);
    bucketDetailStore.SET_DELETE_SHOW(false);
    resetList();
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
        resetList();
        addFolderForm.resetFields()
        bucketDetailStore.fetchList(bucket,prefix);
      })
    }catch(error){

    }
  }

  const pageRender = ()=>{
    const current = bucketDetailStore.currentPage;
    const total = bucketDetailStore.formatList.length;
    const sizes = Math.ceil(total/PAGE_SIZE);
    const ellipsis = ()=> <div className={styles.ellipsis}>•••</div>;
   
    if(current < 5){
      return <>
        {
          Array(sizes>5?5:sizes).fill(0).map((n,index)=>{
            return <div className={classNames(styles.pagerItem,bucketDetailStore.currentPage === index+1?styles.active:'')} key={index} 
            onClick={()=>{
              pageChange(index+1);
            }}>{index+1}</div>
          })
        }
        {sizes>5?ellipsis():<></>}
      </>
    }
    
    if(current >= 5){
      let list:number[] = [];
      if(sizes - current > 2){
        list = [-2,-1,0,1,2]
      }
      if(sizes - current === 2){
        list = [-2,-1,0,1]
      }
      if(sizes - current === 1){
        list = [-3,-2,-1,0]
      }
      if(sizes - current === 0){
        list = [-4,-3,-2,-1]
      }
      console.log(list,'list12');
      
      const first = ()=> {
        return <div className={classNames(styles.pagerItem,bucketDetailStore.currentPage === 1?styles.active:'')} key={1} onClick={()=>{
          pageChange(1);
        }}>{1}</div>
      }
      const last = ()=> {
        return <div className={classNames(styles.pagerItem,bucketDetailStore.currentPage === sizes?styles.active:'')} key={sizes} onClick={()=>{
          pageChange(sizes)
        }}>{sizes}</div>
      }
      
      return <>
      {
        first()
      }
      {
        current - 2 > 1 ? ellipsis():<></>
      }
      {
        list.map((n,index)=>{
          return <div className={classNames(styles.pagerItem,bucketDetailStore.currentPage === current - n?styles.active:'')} key={current - n} onClick={()=>{
            pageChange(current + n);
          }}>{current + n}</div>
        })
      }
      {
        sizes - current > 3 ? ellipsis():<></>
      }
      {
        last()
      }
    </>
    }
  }

  const prev = ()=>{
    const current = bucketDetailStore.currentPage;
    return <div className={
      classNames(styles.pagerItem,bucketDetailStore.currentPage === 1 ? styles.disabled:'')
      }
      onClick={()=>{
        if(current === 1) return;
        bucketDetailStore.SET_CURRENT_PAGE(current - 1);
      }}
    >
      <LeftOutlined></LeftOutlined>
    </div>
  }

  const next = ()=>{
    const current = bucketDetailStore.currentPage;
    const total = bucketDetailStore.formatList.length;
    const sizes = Math.ceil(total/PAGE_SIZE);
    return <div className={
      classNames(styles.pagerItem,!bucketDetailStore.isTruncated && bucketDetailStore.currentPage === sizes ? styles.disabled:'')}
      onClick={()=>{
        if(bucketDetailStore.isTruncated && current === sizes){
          bucketDetailStore.SET_CURRENT_PAGE(current + 1);
          bucketDetailStore.fetchList(bucket,prefix);
        }

        if(current < sizes ) {
          bucketDetailStore.SET_CURRENT_PAGE(current + 1);
        }
      }}
      >
      <RightOutlined></RightOutlined>
    </div>
  }

  const pageChange = (page)=>{
    console.log(page,'page21');
    
    bucketDetailStore.SET_CURRENT_PAGE(page);
  }


  return <div className={styles.objects}>
    <Action bucket={bucket} prefix={prefix}></Action>
    <ObjectTable bucket={bucket} prefix={prefix}></ObjectTable>
    <div className={styles.pager}>
      {
        prev()
      }
      {
        pageRender()
      }
      {
        next()
      }
    </div>

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

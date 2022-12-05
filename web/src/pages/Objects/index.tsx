import Action from "@/components/Objects/action";
import ObjectTable from "@/components/Objects/table";
import objectsStore from "@/store/modules/objects";
import { download } from "@/utils";
import { Modal } from "antd";
import { observer } from "mobx-react";
import { useEffect } from "react";
import { useLocation } from "react-router";
import ReactPlayer from 'react-player';
import styles from './style.module.scss';

interface LocationParams {
  bucket: string;
  Created:string;
}

const Objects = (props:any) => {
  const {
    state: { bucket,Created },
  } = useLocation<LocationParams>();
  useEffect(()=>{
    objectsStore.fetchList(bucket);
  },[bucket])

  const confirmDelete = async ()=>{
    const name = objectsStore.deleteName;
    const path = `${bucket}/${name}`;
    await objectsStore.fetchDelete(path);
    objectsStore.SET_DELETE_SHOW(false);
    objectsStore.fetchList(bucket);
  };
  const cancelDelete = ()=>{
    objectsStore.SET_DELETE_SHOW(false);
  };

  const downloadPreview = ()=>{
    download(objectsStore.downloadFile,objectsStore.downloadName);
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


  return <div className={styles.objects}>
    <Action bucket={bucket} Created={Created}></Action>
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
  </div>;
};

export default observer(Objects);

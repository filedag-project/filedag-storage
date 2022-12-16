import { HttpMethods, Axios } from '@/api/https';
import { SignModel } from '@/models/SignModel';
import { action, computed, makeObservable, observable } from 'mobx';
import _ from 'lodash';
import { formatBytes, formatDate, getExpiresDate } from '@/utils';
import { PreSignModel } from '@/models/PreSignModel';
import presignV4 from '@/api/presign';
import { FileType } from '@/models/BucketModel';

class BucketDetailStore {
  contentsList:any[] = [];
  commonPrefixesList:any[] = [];
  deleteShow:boolean = false;
  previewShow:boolean = false;
  previewText:string = '';
  previewUrl:string = '';
  actionName:string = '';
  downloadFile:Blob = new Blob();
  previewVideo: string = '';
  contentType: string ='';
  shareShow:boolean = false;
  addFolderShow:boolean = false;
  shareLink:string='';
  expiresDate:string = getExpiresDate(7*24*60*60);
  maxDay:number = 7;
  maxSecond:number = 7*24*60*60;
  shareSecond:number = 7*24*60*60;
  percentage:number = 0;
  constructor() {
    makeObservable(this, {
      deleteShow:observable,
      contentsList: observable,
      commonPrefixesList: observable,
      actionName:observable,
      previewShow:observable,
      previewText:observable,
      previewUrl:observable,
      shareShow:observable,
      downloadFile:observable,
      previewVideo:observable,
      contentType:observable,
      shareLink:observable,
      shareSecond:observable,
      percentage:observable,
      addFolderShow:observable,
      formatList: computed,
      totalObjects:computed,
      fetchList: action,
      fetchUpload:action,
      fetchDelete:action,
      fetchShare:action,
      SET_CONTENTS_LIST:action,
      SET_ADD_FOLDER_SHOW:action,
      fetchUploadId:action,
      fetchSliceUpload:action,
      fetchSliceUploadComplete:action,
      fetchUploadFolder:action,
      fetchAbort:action,
    });
  }

  get formatList() {
    const contentsList = this.contentsList.map(n=>{
      const _size = _.get(n,'Size._text','0');
      const _size_ = formatBytes(Number(_size));
      const _ETag = _.get(n,'ETag._text','');
      const _ETag_ = _ETag.replace(/"/g,'');
      const _LastModified = _.get(n,'LastModified._text','');
      const _LastModified_ = formatDate(_LastModified);
      return {
        Name:_.get(n,'Key._text',''),
        LastModified:_LastModified_,
        Size:_size_,
        ETag:_ETag_,
        Type:FileType.file
      }
    }) || [];
    const commonPrefixesList = this.commonPrefixesList.map(n=>{
      return {
        Name:_.get(n,'Prefix._text',''),
        LastModified:'-',
        Size:'-',
        ETag:'-',
        Type:FileType.folder
      }
    }) || [];
    return [...contentsList,...commonPrefixesList];
  }

  get totalObjects(){
    return this.contentsList.length??0 + this.commonPrefixesList.length??0;
  }

  SET_DELETE_SHOW(data:boolean){
    this.deleteShow = data;
  }
  SET_PREVIEW_SHOW(data:boolean){
    this.previewShow = data;
  }
  SET_ACTION_NAME(data:string){
    this.actionName = data;
  }
  SET_ADD_FOLDER_SHOW(data:boolean){
    this.addFolderShow = data;
  }
  SET_OBJECTS_URL(data:string){
    this.previewUrl = data;
  }
  SET_OBJECTS_TEXT(data:string){
    this.previewText = data;
  }
  SET_DOWNLOAD_FILE(data:Blob){
    this.downloadFile = data;
  }
 
  SET_SHARE_SHOW(data:boolean){
    this.shareShow = data;
  }

  SET_SHARE_LINK(data:string){
    this.shareLink = data;
  }

  SET_SHARE_SECOND(data:number){
    this.shareSecond = data;
  }

  SET_EXPIRES_DATE(data:string){
    this.expiresDate = data;
  }

  SET_CONTENTS_LIST(data:any[]){
    this.contentsList = data
  }
  SET_COMMON_PREFIXES_LIST(data:any[]){
    this.commonPrefixesList = data;
  }

  SET_PERCENTAGE(data:number){
    this.percentage = data;
  }

  fetchList(bucket:string,prefix:string) {
    return new Promise(async (resolve) => {
      const params:SignModel = {
        service: 's3',
        body: '',
        protocol: 'http',
        method: HttpMethods.get,
        applyChecksum: true,
        path:bucket,
        region: '',
        query:{
          prefix
        }
      }
      const res = await Axios.axiosXMLStream(params);
      const contents = _.get(res,'ListBucketResult.Contents',[]);
      const commonPrefixes = _.get(res,'ListBucketResult.CommonPrefixes',[]);
      const _contents = Array.isArray(contents) ? contents : [contents];
      const _commonPrefixes = Array.isArray(commonPrefixes) ? commonPrefixes : [commonPrefixes];
      this.SET_CONTENTS_LIST(_contents);
      this.SET_COMMON_PREFIXES_LIST(_commonPrefixes);
      resolve(res)
    })
  }

  fetchObject(bucket,object) {
    this.reset();
    return new Promise(async (resolve) => {
      const params:SignModel = {
        service: 's3',
        body: '',
        protocol: 'http',
        method: HttpMethods.get,
        applyChecksum: true,
        path:`${bucket}/${object}`,
        region: ''
      }
      const res = await Axios.axiosStream(params);
      const contentType = _.get(res,'headers.content-type');
      this.contentType = contentType;
      const body = _.get(res,'body');
      const blob = await new Response(body, { 
        headers: { "Content-Type": contentType } 
      }).blob();
      
      this.downloadFile = blob;
      this.actionName = object;
      if(contentType.includes('image')){
        const objectURL:string = URL.createObjectURL(blob);
        this.previewUrl = objectURL;
      }
      if(contentType.includes('text')){
        const text = await blob.text();
        this.previewText = text;
      }
      if(contentType.includes('video')){
        const objectURL:string = URL.createObjectURL(blob);
        this.previewVideo = objectURL;
      }
      resolve(blob);
    })
  }

  fetchUpload(path:string,body){
    return new Promise(async (resolve) => {
      const params:SignModel = {
        service: 's3',
        body: body,
        protocol: 'http',
        method: HttpMethods.put,
        applyChecksum: true,
        path,
        region: ''
      }
      const res = await Axios.axiosUpload(params,(event)=>{
        if (event.lengthComputable) {
          var complete = (event.loaded / event.total * 100 | 0);
          this.percentage = complete;
        }
      });
      resolve(res)
    })
  }

  fetchDelete(path:string){
    return new Promise(async (resolve) => {
      const params:SignModel = {
        service: 's3',
        body: '',
        protocol: 'http',
        method: HttpMethods.delete,
        applyChecksum: true,
        path,
        region: ''
      }
      const res = await Axios.axiosXMLStream(params);
      resolve(res)
    })
  }

  async fetchShare(url:string,expiresIn:number){
    const params:PreSignModel = {
      region:'',
      expiresIn,
      path:url,
    }
    const res = await presignV4(params);
    const { headers,path,query } = res;
    let _params = ''
    for(var key in query){
      _params+=`${key}=${query[key]}&`
    }
    const str = `${headers.host}${path}?${_params}`;
    this.shareLink = str;
  }

  fetchUploadId(path:string){
    return new Promise(async (resolve) => {
      const params:SignModel = {
        service: 's3',
        body: '',
        protocol: 'http',
        method: HttpMethods.post,
        applyChecksum: true,
        path,
        region: '',
        query:{
          uploads:''
        }
      }
      const res = await Axios.axiosXMLStream(params);
      const id:string = _.get(res,'InitiateMultipartUploadResult.UploadId._text','');
      resolve(id);
    })
  }

  fetchSliceUpload(path:string,index:number,uploadId:string,body){
     return new Promise(async (resolve,reject) => {
       const params:SignModel = {
         service: 's3',
         body: body,
         protocol: 'http',
         method: HttpMethods.put,
         applyChecksum: true,
         path,
         region: '',
         query:{
          partNumber:index.toString(),
          uploadId
         }
       }
       Axios.axiosXMLStream(params).then(res=>{
          resolve(res)
       }).catch(()=>{
          reject();
       });
     })
   }

   fetchAbort(path:string,uploadId:string){
    return new Promise(async (resolve, reject) => {
      const params:SignModel = {
        service: 's3',
        body: '',
        protocol: 'http',
        method: HttpMethods.delete,
        applyChecksum: true,
        path,
        region: '',
        query:{
         uploadId
        }
      }
      const res = await Axios.axiosXMLStream(params);
      resolve(res);
    })
    
   }

   fetchSliceUploadComplete(path:string,uploadId,body:string){
     return new Promise(async (resolve) => {
       const params:SignModel = {
          service: 's3',
          body: body,
          protocol: 'http',
          method: HttpMethods.post,
          applyChecksum: true,
          path,
          region: '',
          query:{
            uploadId
          }
       }
       const res = await Axios.axiosXMLStream(params);
       resolve(res)
     })
   }

   fetchUploadFolder(path:string){
    return new Promise(async (resolve) => {
      const params:SignModel = {
         service: 's3',
         body: '',
         protocol: 'http',
         method: HttpMethods.put,
         applyChecksum: true,
         path,
         region: '',
         query:{}
      }
      const res = await Axios.axiosXMLStream(params);
      resolve(res)
    })
   }

  reset(){
    this.previewUrl = '';
    this.previewText = '';
    this.downloadFile = new Blob();
    this.actionName = '';
    this.previewVideo = '';
    this.contentType = '';
  }
}

const bucketDetailStore = new BucketDetailStore();

export default bucketDetailStore;

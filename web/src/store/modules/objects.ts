import { HttpMethods, Axios } from '@/api/https';
import { SignModel } from '@/models/SignModel';
import { action, computed, makeObservable, observable } from 'mobx';
import _ from 'lodash';
import { formatBytes, formatDate } from '@/utils';
class ObjectsStore {
  objectsList:any[] = [];
  deleteShow:boolean = false;
  previewShow:boolean = false;
  previewText:string = '';
  previewUrl:string = '';
  deleteName:string = '';
  previewName:string = '';
  downloadFile:Blob = new Blob();
  downloadName:string = '';
  previewVideo: string = '';
  contentType: string ='';

  constructor() {
    makeObservable(this, {
      deleteShow:observable,
      objectsList: observable,
      deleteName:observable,
      previewShow:observable,
      previewName:observable,
      previewText:observable,
      previewUrl:observable,
      downloadFile:observable,
      downloadName:observable,
      previewVideo:observable,
      contentType:observable,
      formatList: computed,
      totalSize:computed,
      totalObjects:computed,
      fetchList: action,
      fetchUpload:action,
      fetchDelete:action
    });
  }

  get formatList() {
    return this.objectsList.map(n=>{
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
      }
    }) || [];
  }

  get totalSize(){
    const res = this.objectsList.reduce((total,current)=>{
      const _current = _.get(current,'Size._text','0');
      return Number(total) + Number(_current);
    },0)
    return formatBytes(res);
  }

  get totalObjects(){
    return this.objectsList.length??0
  }

  SET_DELETE_SHOW(data:boolean){
    this.deleteShow = data;
  }
  SET_PREVIEW_SHOW(data:boolean){
    this.previewShow = data;
  }
  SET_PREVIEW_NAME(data:string){
    this.previewName = data;
  }
  SET_DELETE_NAME(data:string){
    this.deleteName = data;
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
  SET_DOWNLOAD_NAME(data:string){
    this.downloadName = data;
  }

  fetchList(bucket) {
    return new Promise(async (resolve) => {
      const params:SignModel = {
        service: 's3',
        body: '',
        protocol: 'http',
        method: HttpMethods.get,
        applyChecksum: true,
        path:bucket,
        region: ''
      }
      const res = await Axios.axiosXMLStream(params);
      const _list:[] = _.get(res,'ListBucketResult.Contents',[]);
      if(Array.isArray(_list)){
        this.objectsList = _list;
      }else{
        this.objectsList = [_list];
      }
      resolve(_list)
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
      this.downloadName = object;
      
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
      const res = await Axios.axiosXMLStream(params);
      const _list:[] = _.get(res,'ListBucketResult.Contents',[]);
      if(Array.isArray(_list)){
        this.objectsList = _list;
      }else{
        this.objectsList = [_list];
      }
      resolve(_list)
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
      const _list:[] = _.get(res,'ListBucketResult.Contents',[]);
      if(Array.isArray(_list)){
        this.objectsList = _list;
      }else{
        this.objectsList = [_list];
      }
      resolve(_list)
    })
  }

  reset(){
    this.previewUrl = '';
    this.previewText = '';
    this.downloadFile = new Blob();
    this.downloadName = '';
    this.previewVideo = '';
    this.contentType = '';
  }
}

const objectsStore = new ObjectsStore();

export default objectsStore;

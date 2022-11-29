import { action, computed, makeObservable, observable } from 'mobx';
import { SignModel } from '@/models/SignModel';
import _ from 'lodash';
import { HttpMethods, Axios } from '@/api/https';
import { formatDate } from '@/utils';
class BucketsStore {
  bucketList:any[] = [];
  deleteShow: boolean = false;
  deleteName:string = '';
  constructor() {
    makeObservable(this, {
      bucketList: observable,
      deleteShow :observable,
      deleteName:observable,
      formatList: computed,
      fetchList: action,
      fetchCreate:action,
      fetchDelete:action,
    });
  }

  get formatList() {
    return this.bucketList.map((n)=>{
      return {
        CreationDate: formatDate(_.get(n,'CreationDate._text','')),
        Name:_.get(n,'Name._text','')
      }
    });
  }

  SET_DELETE_SHOW(data:boolean){
    this.deleteShow = data;
  }

  SET_DELETE_NAME(data:string){
    this.deleteName = data;
  }

  fetchList() {
    return new Promise(async (resolve) => {
      const params:SignModel = {
        service: 's3',
        applyChecksum: true,
        body: '',
        protocol: 'http',
        method: HttpMethods.get,
        path:'/',
        region: ''
      }
      const res = await Axios.axiosXMLStream(params);
      const _list:[] = _.get(res,'ListAllMyBucketsResult.Buckets.Bucket',[]);
      if(Array.isArray(_list)){
        this.bucketList = _list;
      }else{
        this.bucketList = [_list];
      }
      resolve(_list)
    })
  }

  fetchCreate(path:string) {
    return new Promise(async (resolve) => {
      const params:SignModel = {
        service: 's3',
        body: '',
        protocol: 'http',
        method: HttpMethods.put,
        applyChecksum: true,
        path: path,
        region: ''
      }
      const res = await Axios.axiosXMLStream(params);
      resolve(res);
    })
    
  }

  fetchDelete(path:string) {
    return new Promise(async (resolve) => {
      const params:SignModel = {
        service: 's3',
        body: '',
        protocol: 'http',
        method: HttpMethods.delete,
        applyChecksum: true,
        path: path,
        region: ''
      }
      const res = await Axios.axiosXMLStream(params);
      resolve(res);
    })
    
  }

}

const bucketsStore = new BucketsStore();

export default bucketsStore;

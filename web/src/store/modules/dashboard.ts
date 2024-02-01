import { action, computed, makeObservable, observable } from 'mobx';
import { SignModel } from '@/models/SignModel';
import _ from 'lodash';
import { HttpMethods, Axios } from '@/api/https';
import {  objType } from '@/models/DashboardModel';

class DashboardStore {
  get_obj_bytes: objType[] = [];
  put_obj_bytes: objType[] = [];
  get_obj_count: objType[] = [];
  put_obj_count: objType[] = [];
  bucketsCount:number = 0;
  objectsCount:number = 0;
  totalCaptivity:number = 0;
  objectsTotalSize:number = 0;
  constructor() {
    makeObservable(this, {
      get_obj_bytes: observable,
      put_obj_bytes: observable,
      get_obj_count: observable,
      put_obj_count: observable,
      bucketsCount: observable,
      objectsCount: observable,
      totalCaptivity: observable,
      objectsTotalSize: observable,
      top_20_get_obj_bytes:computed,
      top_20_put_obj_bytes:computed,
      top_20_get_obj_count:computed,
      top_20_put_obj_count:computed,
      fetchRequestOverview: action,
    });
  }

  get top_20_get_obj_bytes(){
    const _sort = sortArray(this.get_obj_bytes);
    const _top20 = top20(_sort);
    return _top20;
  }
  get top_20_put_obj_bytes(){
    const _sort = sortArray(this.put_obj_bytes);
    const _top20 = top20(_sort);
    return _top20;
  }
  get top_20_get_obj_count(){
    const _sort = sortArray(this.get_obj_count);
    const _top20 = top20(_sort);
    return _top20;
  }
  get top_20_put_obj_count(){
    const _sort = sortArray(this.put_obj_count);
    const _top20 = top20(_sort);
    return _top20;
  }
  

  fetchRequestOverview() {
    return new Promise(async (resolve) => {
      const params:SignModel = {
        service: 's3',
        body: '',
        protocol: 'http',
        method: HttpMethods.get,
        applyChecksum: true,
        path:`/admin/v1/request-overview`,
        query:{},
        region: '',
      }
      const res = await Axios.axiosJson(params);
      const get_obj_bytes = _.get(res,'Response.get_obj_bytes',[])??[];
      const put_obj_bytes = _.get(res,'Response.put_obj_bytes',[])??[];
      const get_obj_count = _.get(res,'Response.get_obj_count',[])??[];
      const put_obj_count = _.get(res,'Response.put_obj_count',[])??[];
      this.get_obj_bytes =  get_obj_bytes.map((n)=>{
        return {
          name:n.filetype,
          value:n.value
        }
      })
      this.put_obj_bytes =  put_obj_bytes.map((n)=>{
        return {
          name:n.filetype,
          value:n.value
        }
      })
      this.get_obj_count = get_obj_count.map((n)=>{
        return {
          name:n.filetype,
          value:n.value
        }
      })
      this.put_obj_count = put_obj_count.map((n)=>{
        return {
          name:n.filetype,
          value:n.value
        }
      })
      resolve(res);
    })
  }

  fetchStorePool(){
    return new Promise(async (resolve, reject) => {
      const params:SignModel = {
        service: 's3',
        body: '',
        protocol: 'http',
        method: HttpMethods.get,
        applyChecksum: true,
        path:`/admin/v1/store-pool-stats`,
        query:{},
        region: '',
      }
      const res = await Axios.axiosJson(params);
      const bucketsCount = _.get(res,'Response.bucketsCount',0);
      const objectsCount = _.get(res,'Response.objectsCount',0);
      const totalCaptivity = _.get(res,'Response.totalCaptivity',0);
      const objectsTotalSize = _.get(res,'Response.objectsTotalSize',0);
      this.bucketsCount = bucketsCount;
      this.objectsCount = objectsCount;
      this.totalCaptivity = totalCaptivity;
      this.objectsTotalSize = objectsTotalSize;
      resolve(res);
    })
  }
}

const sortArray = (data:Array<any>)=>{
  const res =  data.sort((a,b)=>{
    return  b.value - a.value;
  })
  return res;
}

const top20 = (data:Array<any>)=>{
  return data.length >= 20 ? data.splice(0,20):data;
}

const dashboardStore = new DashboardStore();

export default dashboardStore;

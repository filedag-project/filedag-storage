import { action, makeObservable, observable } from 'mobx';
import { SignModel } from '@/models/SignModel';
import _ from 'lodash';
import { HttpMethods, Axios } from '@/api/https';
import {  objType } from '@/models/DashboardModel';

class DashboardStore {
  get_obj_bytes: objType[] = [];
  put_obj_bytes: objType[] = [];
  get_obj_count: objType[] = [];
  put_obj_count: objType[] = [];

  constructor() {
    makeObservable(this, {
      get_obj_bytes: observable,
      put_obj_bytes: observable,
      get_obj_count: observable,
      put_obj_count: observable,
      fetchRequestOverview: action,
    });
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
}

const dashboardStore = new DashboardStore();

export default dashboardStore;

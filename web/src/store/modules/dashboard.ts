import { action, makeObservable, observable } from 'mobx';
import { SignModel } from '@/models/SignModel';
import _ from 'lodash';
import { HttpMethods, Axios } from '@/api/https';

interface overviewType {
  get_obj_bytes:string,
  get_obj_count:string,
  put_obj_bytes:string,
  put_obj_count:string,
}
class DashboardStore {
  userInfo:string = '';
  overview:overviewType = {
    get_obj_bytes:'',
    get_obj_count:'',
    put_obj_bytes:'',
    put_obj_count:''
  }
  json:any = {};
  constructor() {
    makeObservable(this, {
      userInfo: observable,
      json:observable,
      fetchRequestOverview: action,
      fetchUserInfos:action,
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
      const res = await Axios.axiosJsonAWS(params);
      const get_obj_bytes = _.get(res,'Response.get_obj_bytes','');
      const get_obj_count = _.get(res,'Response.get_obj_count','');
      const put_obj_bytes = _.get(res,'Response.put_obj_bytes','');
      const put_obj_count = _.get(res,'Response.get_obj_count','');
      this.overview = {
        get_obj_bytes,
        get_obj_count,
        put_obj_bytes,
        put_obj_count
      }
    })
  }

  fetchUserInfos(){
    return new Promise(async (resolve, reject) => {
      const params:SignModel = {
        service: 's3',
        body: '',
        protocol: 'http',
        method: HttpMethods.get,
        applyChecksum: true,
        path:`/admin/v1/user-infos`,
        query:{},
        region: '',
        contentType:'application/json;charset=UTF-8'
      }
      const res = await Axios.axiosXMLStream(params);
    })
  }
}

const dashboardStore = new DashboardStore();

export default dashboardStore;

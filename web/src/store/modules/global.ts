import { action, makeObservable, observable } from 'mobx';
import { SignModel } from '@/models/SignModel';
import { HttpMethods, Axios } from '@/api/https';
import _ from 'lodash';
import { ACCESS_KEY_ID, Cookies } from '@/utils/cookies';
import { formatBytes } from '@/utils';

interface userInfoType {
  total_storage_capacity:string;
  use_storage_capacity:string;
  buckets:number;
  objects:string;
  account_name:string;
}

class GlobalStore {
  isAdmin:boolean = false;
  userInfo:userInfoType = {
    total_storage_capacity:'0',
    use_storage_capacity:'0',
    buckets:0,
    objects:'0',
    account_name:''
  };
  constructor() {
    makeObservable(this, {
      isAdmin: observable,
      userInfo: observable,
      fetchIsAdmin: action,
    });
  }

  fetchIsAdmin() {
    return new Promise(async (resolve) => {
      const params:SignModel = {
        service: 's3',
        body: '',
        protocol: 'http',
        method: HttpMethods.get,
        applyChecksum: true,
        path:`/admin/v1/is-admin`,
        query:{},
        region: '',
      }
      const res = await Axios.axiosJsonAWS(params);
      const _response = _.get(res,'Response',false);
      this.isAdmin = _response;
    })
  }

  fetchUserInfo() {
    return new Promise(async (resolve) => {
      try {
        const accessKey = Cookies.getKey(ACCESS_KEY_ID);
        if(!accessKey) return;
        const params:SignModel = {
          service: 's3',
          body: '',
          protocol: 'http',
          method: HttpMethods.get,
          applyChecksum: true,
          path:'/console/v1/user-info',
          region: '',
          query:{
            "accessKey":Cookies.getKey(ACCESS_KEY_ID),
          },
          contentType:'application/json; charset=UTF-8'
        }
        const res = await Axios.axiosJson(params);
        const total:string = _.get(res,'total_storage_capacity',0);
        console.log(total,'total 90');
        
        const _total = formatBytes(total);
        console.log(_total,'total 90');
        const use = _.get(res,'use_storage_capacity',0);
        const _use:string = formatBytes(use);
        const buckets = _.get(res,'bucket_infos');
        const size = buckets.reduce((total,current)=>{
          const _current = _.get(current,'size','0');
          return Number(total) + Number(_current);
        },0);
        const name = _.get(res,'account_name')
        const _size = formatBytes(size);
        this.userInfo = {
          total_storage_capacity:_total,
          use_storage_capacity:_use,
          buckets: buckets.length,
          objects: _size,
          account_name: name
        };
      } catch (e) {
    
      }
    })
  }

}

const globalStore = new GlobalStore();

export default globalStore;

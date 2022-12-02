import { action, makeObservable, observable } from 'mobx';
import { SignModel } from '@/models/SignModel';
import _ from 'lodash';
import { HttpMethods, Axios } from '@/api/https';
import { ACCESS_KEY_ID, Cookies } from '@/utils/cookies';
import { formatBytes } from '@/utils';
interface userInfoType {
  total_storage_capacity:string;
  use_storage_capacity:string;
  buckets:number;
  objects:string;
}

class OverviewStore {
  userInfo:userInfoType = {
    total_storage_capacity:'0',
    use_storage_capacity:'0',
    buckets:0,
    objects:'0'
  };
  constructor() {
    makeObservable(this, {
      userInfo: observable,
      fetchUserInfo: action,
    });
  }

  fetchUserInfo() {
    return new Promise(async (resolve) => {
      try {
        const params:SignModel = {
          service: 's3',
          body: '',
          protocol: 'http',
          method: HttpMethods.get,
          applyChecksum: true,
          path:'/admin/v1/user-info',
          region: '',
          query:{
            "accessKey":Cookies.getKey(ACCESS_KEY_ID),
          },
          contentType:'application/json; charset=UTF-8'
        }
        const res = await Axios.axiosJson(params);
        const total:string = _.get(res,'total_storage_capacity');
        const _total = formatBytes(total);
        const use = _.get(res,'use_storage_capacity');
        const _use:string = formatBytes(use);
        const buckets = _.get(res,'bucket_infos');
        const size = buckets.reduce((total,current)=>{
          const _current = _.get(current,'size','0');
          return Number(total) + Number(_current);
        },0);
        
        const _size = formatBytes(size);
        this.userInfo = {
          total_storage_capacity:_total,
          use_storage_capacity:_use,
          buckets: buckets.length,
          objects: _size,
        };
      } catch (e) {
    
      }
    })
  }
}

const overviewStore = new OverviewStore();

export default overviewStore;

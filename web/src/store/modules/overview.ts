import { action, makeObservable, observable } from 'mobx';
import { SignModel } from '@/models/SignModel';
import { get } from 'lodash';
import { HttpMethods, Axios } from '@/api/https';
import { ACCESS_KEY_ID, Cookies } from '@/utils/cookies';
import { formatBytes } from '@/utils';
interface userInfoType {
  total_storage_capacity:string;
  use_storage_capacity:string;
  buckets:number;
  objects:string;
  account_name:string;
}

class OverviewStore {
  userInfo:userInfoType = {
    total_storage_capacity:'0',
    use_storage_capacity:'0',
    buckets:0,
    objects:'0',
    account_name:''
  };
  constructor() {
    makeObservable(this, {
      userInfo: observable,
      fetchUserInfo: action,
    });
  }

  SET_USER_INFO(data:userInfoType){
    this.userInfo = data
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
        const total:string = get(res,'Response.total_storage_capacity');
        const _total = formatBytes(total);
        const use = get(res,'Response.use_storage_capacity');
        const _use:string = formatBytes(use);
        const buckets = get(res,'Response.buckets_count',0)??0;
        const objects = get(res,'Response.objects_count',0)??0;
        const name = get(res,'account_name')
        
        this.SET_USER_INFO ({
          total_storage_capacity:_total,
          use_storage_capacity:_use,
          buckets: buckets,
          objects: objects,
          account_name: name
        });
      } catch (e) {
        this.SET_USER_INFO({
          total_storage_capacity:'0',
          use_storage_capacity:'0',
          buckets:0,
          objects:'0',
          account_name:''
        })
      }
    })
  }
}

const overviewStore = new OverviewStore();

export default overviewStore;

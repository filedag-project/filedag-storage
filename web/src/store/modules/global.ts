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
}

class GlobalStore {
  userInfo:userInfoType = {
    total_storage_capacity:'0',
    use_storage_capacity:'0',
    buckets:0,
    objects:'0',
  };
  constructor() {
    makeObservable(this, {
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
      resolve(_response);
    })
  }

}

const globalStore = new GlobalStore();

export default globalStore;

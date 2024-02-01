import { action, makeObservable, observable } from 'mobx';
import { SignModel } from '@/models/SignModel';
import { HttpMethods, Axios } from '@/api/https';
import _ from 'lodash';
class PowerStore {
  userInfo:string = '';
  json:any = {};
  powerName:string = '';
  constructor() {
    makeObservable(this, {
      userInfo: observable,
      json:observable,
      powerName:observable,
      fetchGetPower: action,
      fetchPutPower:action,
      fetchGetPowerName:action,
    });
  }

  SET_JSON(data:string){
    this.json = data
  }

  fetchGetPower(path:string) {
    return new Promise(async (resolve) => {
      const params:SignModel = {
        service: 's3',
        body: '',
        protocol: 'http',
        method: HttpMethods.get,
        applyChecksum: true,
        path:`${path}`,
        query:{
          policy:''
        },
        region: '',
      }
      const res = await Axios.axiosJsonAWS(params);
      const json = JSON.stringify(res);
      this.json = json
      resolve(json);
    })
  }

  fetchPutPower(path:string,json) {
    return new Promise(async (resolve) => {
      const params:SignModel = {
        service: 's3',
        body: json,
        protocol: 'http',
        method: HttpMethods.put,
        applyChecksum: true,
        path:`${path}`,
        query:{
          policy:'',
        },
        region: '',
      }
      const res = await Axios.axiosJsonAWS(params);
      resolve(res);
    })
  }

  fetchGetPowerName(bucket:string,json) {
    return new Promise(async (resolve) => {
      const params:SignModel = {
        service: 's3',
        body: json,
        protocol: 'http',
        method: HttpMethods.post,
        applyChecksum: true,
        path:`console/v1/get-policy-name`,
        query:{
          bucketName:bucket,
        },
        region: '',
      }
      const res = await Axios.axiosJsonAWS(params);
      const name = _.get(res,'Response','');
      this.powerName = name;
      resolve(name);
    })
  }
}

const powerStore = new PowerStore();

export default powerStore;

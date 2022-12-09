import { action, makeObservable, observable } from 'mobx';
import { SignModel } from '@/models/SignModel';
import { HttpMethods, Axios } from '@/api/https';

class PowerStore {
  userInfo:string = '';
  json:any = {};
  constructor() {
    makeObservable(this, {
      userInfo: observable,
      json:observable,
      fetchGetPower: action,
      fetchPutPower:action,
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
      this.json = JSON.stringify(res);
      resolve(res);
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
}

const powerStore = new PowerStore();

export default powerStore;

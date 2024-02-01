import { action, makeObservable, observable } from 'mobx';
import { SignModel } from '@/models/SignModel';
import _ from 'lodash';
import { HttpMethods, Axios } from '@/api/https';
import { formatBytes } from '@/utils';
import { addUserType, changePasswordType, userStatusType, userType } from '@/models/DashboardModel';

class UserStore {
  userInfos: userType[] = [];
  constructor() {
    makeObservable(this, {
      userInfos : observable,
      
      fetchUserInfos:action,
      fetchDeleteUser:action,
      fetchChangeUserPassword:action,
      fetchSetUserStatus:action,
      fetchAddUser:action,
      SET_USER_INFOS:action,
    });
  }

  SET_USER_INFOS(data:userType[]){
    this.userInfos = data
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
      const res = await Axios.axiosJson(params);
      const _list:any[] = _.get(res,'Response',[]);
      const data = _list.map((n)=>{
        const _total = formatBytes(n.total_storage_capacity??0);
        const _use = formatBytes(n.use_storage_capacity??0);
        const _buckets = _.get(n,'bucket_infos',[])??[];
        return {
          account_name:n.account_name,
          bucket_infos:_buckets.length,
          status:n.status,
          total_storage_capacity: _total,
          use_storage_capacity: _use
        }
      });

      this.SET_USER_INFOS(data);
      resolve(data);
    })
  }
  fetchDeleteUser(accessKey:string){
    return new Promise(async (resolve, reject) => {
      const params:SignModel = {
        service: 's3',
        body: '',
        protocol: 'http',
        method: HttpMethods.post,
        applyChecksum: true,
        path:`/admin/v1/remove-user`,
        query:{
          accessKey,
        },
        region: '',
        contentType:'application/json;charset=UTF-8'
      }
      const res = await Axios.axiosJson(params);
      resolve(res);
    })
  }
  fetchChangeUserPassword(data:changePasswordType){
    return new Promise(async (resolve, reject) => {
      const params:SignModel = {
        service: 's3',
        body: '',
        protocol: 'http',
        method: HttpMethods.post,
        applyChecksum: true,
        path:`/admin/v1/change-password`,
        query:{
          ...data
        },
        region: '',
        contentType:'application/json;charset=UTF-8'
      }
      const res = await Axios.axiosJson(params);
      resolve(res);
    })
  }
  fetchSetUserStatus(data:userStatusType){
    return new Promise(async (resolve, reject) => {
      const params:SignModel = {
        service: 's3',
        body: '',
        protocol: 'http',
        method: HttpMethods.post,
        applyChecksum: true,
        path:`/admin/v1/update-accessKey_status`,
        query:{
          ...data
        },
        region: '',
        contentType:'application/json;charset=UTF-8'
      }
      const res = await Axios.axiosJson(params);
      resolve(res);
    })
  }
  fetchAddUser(user:addUserType){
    return new Promise(async (resolve, reject) => {
      const params:SignModel = {
        service: 's3',
        body: '',
        protocol: 'http',
        method: HttpMethods.post,
        applyChecksum: true,
        path:`/admin/v1/add-user`,
        query:{
          ...user
        },
        region: '',
        contentType:'application/json;charset=UTF-8'
      }
      const res = await Axios.axiosJson(params);
      resolve(res);
    })
  }
}

const userStore = new UserStore();

export default userStore;

import { Axios, HttpMethods } from '@/api/https';
import { SignModel } from '@/models/SignModel';
import { observer } from 'mobx-react';
import { useEffect } from 'react';
import styles from './style.module.scss';
const Dashboard = (props:any) => {
  useEffect(()=>{
    get();
    getInfo();
  },[]);
  const get = async ()=>{
    const params:SignModel = {
      service: 's3',
      body: '',
      protocol: 'http',
      method: HttpMethods.get,
      applyChecksum: true,
      path:`/console/v1/request-overview`,
      query:{},
      region: '',
    }
    const res = await Axios.axiosXMLStream(params);
  }
  const getInfo = async ()=>{
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
  }
  return <div className={styles.dashboard}>
   Dashboard
  </div>;
};

export default observer(Dashboard);

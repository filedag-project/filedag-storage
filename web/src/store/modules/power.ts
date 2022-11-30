import { action, makeObservable, observable } from 'mobx';
import { SignModel } from '@/models/SignModel';
// import _ from 'lodash';
import { HttpMethods, Axios } from '@/api/https';
import { ACCESS_KEY_ID, Cookies } from '@/utils/cookies';

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
      console.log(res,'ssss');
      this.json = JSON.stringify(res);
    })
  }

  fetchPutPower(path:string,json) {
    console.log(json,'json2');
    
    return new Promise(async (resolve) => {
      const params:SignModel = {
        service: 's3',
        body: json,
        protocol: 'http',
        method: HttpMethods.put,
        applyChecksum: true,
        path:`${path}`,
        query:{
          policy:''
        },
        region: '',
      }
      const res = await Axios.axiosJsonAWS(params);
      console.log(res,'ssss');
      this.json = JSON.stringify(res);
    })
  }
}

const getPublic = (bucket:string)=>{
  const accessKey = Cookies.getKey(ACCESS_KEY_ID);
  const json = {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Action": [
          "s3:*"
        ],
        "Effect": "Allow",
        "Principal": {
          "AWS": [
            accessKey //////用户
          ]
        },
        "Resource": [
          `arn:aws:s3:::${bucket}/*` //////bucket
        ],
        "Sid": ""
      },
      {
        "Action": [
          "s3:GetBucketLocation",
          "s3:ListBucket",
          "s3:ListBucketMultipartUploads"
        ],
        "Effect": "Allow",
        "Principal": {
          "AWS": [
            "*"
          ]
        },
        "Resource": [
          `arn:aws:s3:::${bucket}`//bucket
        ],
        "Sid": ""
      },
      {
        "Action": [
          "s3:AbortMultipartUpload",
          "s3:DeleteObject",
          "s3:GetObject",
          "s3:ListMultipartUploadParts",
          "s3:PutObject"
        ],
        "Effect": "Allow",
        "Principal": {
          "AWS": [
            "*"
          ]
        },
        "Resource": [
          `arn:aws:s3:::${bucket}/*`//bucket
        ],
        "Sid": ""
      }
    ]
  }
}

const getDownload = (bucket:string)=>{

}



const powerStore = new PowerStore();

export default powerStore;

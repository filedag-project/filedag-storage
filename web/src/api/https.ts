import {FetchHttpHandler} from "@aws-sdk/fetch-http-handler";
import {HttpRequest} from '@aws-sdk/protocol-http';
import { SignModel } from '@/models/SignModel';
import _ from 'lodash';
import { notification } from 'antd';
import signV4 from './sign';
import { ACCESS_KEY_ID, Cookies, SECRET_ACCESS_KEY, SESSION_TOKEN } from '@/utils/cookies';
import { xmlStreamToJs,streamToJs } from "@/utils";

const INVALID_ACCESS_KEY_ID = 'InvalidAccessKeyId';

export enum HttpMethods {
    post = 'POST',
    get = 'GET',
    delete = 'DELETE',
    put = 'PUT'
}

export const Axios = {
  axiosUpload(params:SignModel,progressCb:Function){
    return new Promise(async (resolve, reject) => {
      const sign = await signV4(params);
      const xhr = new XMLHttpRequest();
      xhr.upload.addEventListener("progress", (event)=>{
        progressCb(event);
      });
      xhr.onload = (data)=>{
        resolve(data);
      }
      xhr.onerror = (error)=>{
        reject(error);
      }
      xhr.open(sign.method, process.env['REACT_APP_BASE_URL']+`${sign.path}`);
      xhr.setRequestHeader("authorization", sign.headers.authorization);
      xhr.setRequestHeader("x-amz-content-sha256", sign.headers['x-amz-content-sha256']);
      xhr.setRequestHeader("x-amz-security-token", sign.headers['x-amz-security-token']);
      xhr.setRequestHeader("x-amz-date", sign.headers['x-amz-date']);
      xhr.setRequestHeader("content-type", '');
      xhr.send(params.body);
      
    })
  },
  axiosXMLStream(params:SignModel){
    return new Promise(async (resolve, reject) => {
      const sign = await signV4(params);
      const nodeHttpHandler = new FetchHttpHandler();
      const request = new HttpRequest({
        ...sign,
      })
      nodeHttpHandler.handle(request)
      .then(async result=>{
        const _result = await this.handlerXMLStream(result);
        resolve(_result);
      }).catch(error=>{
        this.handlerError('network error');
        reject();
      })
    })
  },

  axiosStream(params:SignModel){
    return new Promise(async (resolve, reject) => {
      const sign = await signV4(params);
      const nodeHttpHandler = new FetchHttpHandler();
      const request = new HttpRequest({
        ...sign
      })
      nodeHttpHandler.handle(request)
      .then(async result=>{
        const _result = await this.handlerStream(result);
        resolve(_result);
      }).catch(error=>{
        this.handlerError('network error');
        reject();
      })
    })
  },

  axiosJson(params:SignModel){
    return new Promise(async (resolve, reject) => {
      const sign = await signV4(params);
      const nodeHttpHandler = new FetchHttpHandler();
      const request = new HttpRequest({
        ...sign
      })
      nodeHttpHandler.handle(request)
      .then(async result=>{
        const _result = await this.handlerJson(result);
        resolve(_result);
      }).catch(error=>{
        this.handlerError('network error');
        reject();
      })
    })
  },

  axiosJsonAWS(params:SignModel){
    return new Promise(async (resolve, reject) => {
      const sign = await signV4(params);
      const nodeHttpHandler = new FetchHttpHandler();
      const request = new HttpRequest({
        ...sign
      })
      nodeHttpHandler.handle(request)
      .then(async result=>{
        const _result = await this.handlerJsonAWS(result);
        resolve(_result);
      }).catch(error=>{
        this.handlerError('network error');
        reject();
      })
    })
  },

  handlerXMLStream(result){
    return new Promise(async (resolve, reject) => {
      const statusCode = _.get(result,'response.statusCode');
      const body = _.get(result,'response.body');
      const etag = _.get(result,'response.headers.etag','');
      const data = await xmlStreamToJs(body);
      if(statusCode === 200||statusCode===204){
       return resolve({ ...data, etag })
      }else{
        // object list
        const code = _.get(data,'Error.Code._text','');
        const message = _.get(data,'Error.Message._text','Error');
        this.handlerError(message);
        if(code === INVALID_ACCESS_KEY_ID){
          this.handlerLogout()
        }
      }
    })
  },

  handlerStream(result){
    return new Promise(async (resolve, reject) => {
      const statusCode = _.get(result,'response.statusCode');
      const body = _.get(result,'response.body');
      const headers = _.get(result,'response.headers');
      if(statusCode === 200||statusCode===204){
        resolve({
          headers,
          body,
        })
      }else{
        const data = await xmlStreamToJs(body);
        const code = _.get(data,'Error.Code._text','');
        const message = _.get(data,'Error.Message._text','Error');
        this.handlerError(message);
        if(code === INVALID_ACCESS_KEY_ID){
          this.handlerLogout()
        }
      }
    })
  },

  handlerJson(result){
    return new Promise(async (resolve, reject) => {
      const body = _.get(result,'response.body');
      const data = await streamToJs(body);
      const statusCode = _.get(data,'HTTPStatusCode');
      if(statusCode === 200||statusCode===204){
        resolve(data)
      }else{
        const code = _.get(data,'Code','');
        const message = _.get(data,'Message','Error');
        this.handlerError(message);
        if(code === INVALID_ACCESS_KEY_ID){
          this.handlerLogout()
        }
      }
    })
  },

  handlerJsonAWS(result){
    return new Promise(async (resolve, reject) => {
      try{
        const body = _.get(result,'response.body');
        const statusCode = _.get(result,'response.statusCode');
        if(statusCode === 200||statusCode===204){
          const data = await streamToJs(body);
          resolve(data)
        }else{
          const data = await xmlStreamToJs(body);
          const code = _.get(data,'Error.Code._text','');
          const message = _.get(data,'Error.Message._text','Error');
          this.handlerError(message);
          if(code === INVALID_ACCESS_KEY_ID){
            this.handlerLogout()
          }
        }
      }catch(error){

      }
    })
  },

  handlerError(description:string){
    notification.open({
      message: 'Error',
      description: description,
    });
  },

  handlerLogout(){
    Cookies.deleteKey(ACCESS_KEY_ID);
    Cookies.deleteKey(SECRET_ACCESS_KEY);
    Cookies.deleteKey(SESSION_TOKEN);
    setTimeout(()=>{
      window.location.href = `/login`;
    },3000)
  }
}


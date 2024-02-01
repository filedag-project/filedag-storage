import _Cookies from 'js-cookie';

export const ACCESS_KEY_ID = 'dag-accessKeyId';
export const SECRET_ACCESS_KEY = 'dag-secretAccessKey';
export const SESSION_TOKEN = 'dag-sessionToken';
export const USER_NAME = 'dag-username';

export class Cookies {
  static setKey(key:string,value:string, expires?:number):void {
    _Cookies.set(key, value,{ expires });
  }

  static getKey(key:string): string {
    const value = _Cookies.get(key);
    return value != null ? value : '';
  }

  static deleteKey(key:string):void {
    _Cookies.remove(key);
  }
}

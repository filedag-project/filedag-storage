export interface SignModel {
  service: string;
  body: string;
  protocol: string;
  method: string;
  region: string;
  path:string;
  accessKeyId?: string;
  secretAccessKey?: string;
  sessionToken?: string;
  contentType?:string;
  applyChecksum?:boolean;
  "X-Amz-Meta-File-Size"?:string,
  query?:{
    [key:string]:string
  }
}
export interface PreSignModel {
    region: string;
    path:string;
    expiresIn:number;
    accessKeyId?: string;
    secretAccessKey?: string;
    sessionToken?: string;
  }
export interface userType{
    account_name:string,
    status:string,
    total_storage_capacity:string,
    use_storage_capacity:string,
    bucket_infos:number
}

export interface objType{
    name:string,
    value:string
}

export interface addUserType{
    accessKey:string;
    secretKey:string;
    capacity:string;
}

export interface changePasswordType{
    newSecretKey:string;
    accessKey:string;
}

export interface userStatusType{
    accessKey:string;
    status: statusType.on|statusType.off;
}

export enum statusType{
    on = 'on',
    off = 'off'
}
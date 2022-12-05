import convert from 'xml-js';
import { ACCESS_KEY_ID, Cookies } from '@/utils/cookies';

const xmlStreamToJs = async (data) => {
  try{
    const res = await new Response(data, {
      headers: {'Content-Type': 'text/html'}
    })
    .text()
    .then((res) => {
      return convert.xml2js(res, {
        compact: true,
        ignoreDeclaration:true,
        ignoreAttributes:true
      });
    });
    return res;
  }catch(error){
    console.log(error,'streamToJs');
    throw new Error('error');
  }
};

const streamToJs = async (data) => {
  try{
    const res = await new Response(data, {
      headers: {'Content-Type': 'text/html'}
    })
    .json();
    return res;
  }catch(error){
    console.log(error,'streamToJs');
    throw new Error('error');
  }
};

const xmlToJs = async (data) => {
  try{
    return convert.xml2js(data, {
      compact: true,
      ignoreDeclaration:true,
      ignoreAttributes:true
    })
  }catch(error){
    console.log(error,'xmlToJs');
    throw new Error('error');
  }
};

const formatDate = (date:string):string=>{
  const _date = new Date(date);
  return _date.toISOString() ?? ''
}

const formatBytes = (bytes, decimals = 2) =>{
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}


const download = (blob:Blob,name:string)=>{
  let downloadElement = document.createElement('a');
  let href = window.URL.createObjectURL(blob);
  downloadElement.href = href;
  downloadElement.download = name;
  document.body.appendChild(downloadElement);
  downloadElement.click();
  document.body.removeChild(downloadElement);
  window.URL.revokeObjectURL(href);
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
  };
  return json;
}

const getDownload = (bucket:string)=>{
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
            accessKey
          ]
        },
        "Resource": [
          "arn:aws:s3:::books/*"//bucket
        ],
        "Sid": ""
      },
      {
        "Action": [
          "s3:GetBucketLocation",
          "s3:ListBucket"
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
          "s3:GetObject"
        ],
        "Effect": "Allow",
        "Principal": {
          "AWS": [
            "*"
          ]
        },
        "Resource": [
          "arn:aws:s3:::books/*"//bucket
        ],
        "Sid": ""
      }
    ]
  };
  return json;
}

const getUpload = (bucket:string)=>{
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
            accessKey//用户
          ]
        },
        "Resource": [
          "arn:aws:s3:::books/*"//bucket
        ],
        "Sid": ""
      },
      {
        "Action": [
          "s3:GetBucketLocation",
          "s3:ListBucketMultipartUploads"
        ],
        "Effect": "Allow",
        "Principal": {
          "AWS": [
            "*"
          ]
        },
        "Resource": [
          "arn:aws:s3:::books"//bucket
        ],
        "Sid": ""
      },
      {
        "Action": [
          "s3:AbortMultipartUpload",
          "s3:DeleteObject",// 还需商定 上传是否可以有删除权限
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
          `arn:aws:s3:::${bucket}`//
        ],
        "Sid": ""
      }
    ]
  };
  return json;
}

const getPrivate = (bucket:string)=>{
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
            accessKey //用户
          ]
        },
        "Resource": [
          "arn:aws:s3:::books/*"//bucket
        ],
        "Sid": ""
      },
      {
        "Action": [
          "s3:GetBucketLocation"
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
      }
    ]
  };
  return json;
}


export { 
  xmlStreamToJs,
  streamToJs,
  formatDate,
  formatBytes,
  download,
  xmlToJs,
  getPublic,
  getDownload,
  getUpload,
  getPrivate
};

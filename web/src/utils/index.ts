import convert from 'xml-js';
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

const escapeStr = (str:string)=>{
  return str;
};

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


export { 
  xmlStreamToJs,
  streamToJs,
  formatDate,
  formatBytes,
  escapeStr,
  download,
};

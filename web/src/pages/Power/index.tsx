import { observer } from "mobx-react";
import { Button, Radio,Input } from 'antd';
import styles from './style.module.scss';
import { useEffect, useState } from "react";
import powerStore from "@/store/modules/power";
import { useLocation } from "react-router";
import { getDownload, getPrivate, getPublic, getUpload } from "@/utils";
const { TextArea } = Input;

interface LocationParams {
  path: string;
}



const Power = (props:any) => {
  const {
    state: { path }
  } = useLocation<LocationParams>();
  const [selectValue,setSelectValue]=useState('');

  useEffect(()=>{
    console.log(path,'path 234');
    
    powerStore.fetchGetPower(path);
  },[]);

  const radioList = [
    {
      key:'Public',
      value: getPublic(path)
    },
    {
      key:'Download',
      value: getDownload(path)
    },
    {
      key:'Upload',
      value: getUpload(path)
    },
    {
      key:'Private',
      value: getPrivate(path)
    }
  ]

  const radioChange = (e)=>{
    const value = e.target.value;
    const obj = radioList.find(n=> n.key === value);
    console.log(obj,'ddd');
    const json = JSON.stringify(obj?.value);
    powerStore.SET_JSON(json);
    setSelectValue(value);
  };

  const save = ()=>{
    powerStore.fetchPutPower(path,powerStore.json);
  }

  return <div className={styles.power}>
          <Radio.Group onChange={radioChange} value={selectValue}>
            {
              radioList.map((item,index)=>{
                return <div className="radio-item" key={index}>
                <Radio value={item.key}>{item.key}</Radio>
              </div>
              })
            }
          </Radio.Group>
          <TextArea value={powerStore.json} rows={6} />
          <div className="btn-wrap" onClick={save}>
            <Button type="primary">Save</Button>
          </div>
  </div>;
};

export default observer(Power);



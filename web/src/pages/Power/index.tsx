import { observer } from "mobx-react";
import { Button, Radio,Input } from 'antd';
import styles from './style.module.scss';
import { useEffect, useState } from "react";
import powerStore from "@/store/modules/power";
import { useLocation } from "react-router";
import _ from 'lodash';
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
  useEffect(()=>{
    powerStore.fetchGetPower(path).then(res=>{
      radioList.forEach(n=>{
        const _value = n.value;
        if(n.key === 'Public'){
          console.log(_value,1);
          console.log(JSON.parse(powerStore.json),2);
          console.log(JSON.stringify(_value),'ziwei');
          console.log(powerStore.json,'yuguang');
          const ise = _.isEqual(res,_value)
          console.log(ise,3);
          console.log(JSON.stringify(_value) === powerStore.json,5);
          console.log((JSON.stringify(_value)).length,(powerStore.json).length,'6');
          
        }
        if(powerStore.json===_value){
          console.log(123);
          
        }
      })
      
    })
  },[]);

  

  const radioChange = (e)=>{
    const value = e.target.value;
    const obj = radioList.find(n=> n.key === value);
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
import { observer } from "mobx-react";
import { Button, Radio,Input } from 'antd';
import styles from './style.module.scss';
import { useEffect, useState } from "react";
import powerStore from "@/store/modules/power";
import { useLocation } from "react-router-dom";
import { getDownload, getPrivate, getPublic, getUpload } from "@/utils";
const { TextArea } = Input;

const Power = (props:any) => {
  const { state :{ bucket = '' }} = useLocation();
  const [selectValue,setSelectValue] = useState('');
  const radioList = [
    {
      key:'public',
      label:'Public',
      value: getPublic(bucket)
    },
    {
      key:'download',
      label:'Download',
      value: getDownload(bucket)
    },
    {
      key:'upload',
      label:'Upload',
      value: getUpload(bucket)
    },
    {
      key:'private',
      label:'Private',
      value: getPrivate(bucket)
    }
  ]
  useEffect(()=>{
    powerStore.fetchGetPower(bucket).then(res=>{
      powerStore.fetchGetPowerName(bucket,powerStore.json).then(name=>{
        const obj = radioList.find(n=>n.key === name);
        if(obj){
          setSelectValue(obj['key']);
        }
      });
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
    powerStore.fetchPutPower(bucket,powerStore.json);
  }

  return <div className={styles.power}>
          <Radio.Group onChange={radioChange} value={selectValue}>
            {
              radioList.map((item,index)=>{
                return <div className="radio-item" key={index}>
                <Radio value={item.key}>{item.label}</Radio>
              </div>
              })
            }
          </Radio.Group>
          <TextArea value={powerStore.json} rows={6} />
          <div className="btn-wrap" onClick={save}>
            <Button className="bg-btn" type="primary">Save</Button>
          </div>
  </div>;
};

export default observer(Power);
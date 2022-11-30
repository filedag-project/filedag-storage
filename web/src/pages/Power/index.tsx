import { observer } from "mobx-react";
import { Button, Radio, Tabs,Input } from 'antd';
import styles from './style.module.scss';
import { useEffect } from "react";
import powerStore from "@/store/modules/power";
import { useLocation } from "react-router";
const { TextArea } = Input;

interface LocationParams {
  path: string;
}

const Power = (props:any) => {
  const {
    state: { path }
  } = useLocation<LocationParams>();

  useEffect(()=>{
    powerStore.fetchGetPower(path);
  },[]);

  const commonChange = ()=>{

  };

  const save = ()=>{
    powerStore.fetchPutPower(path,JSON.parse(powerStore.json))
  }

  return <div className={styles.power}>
    <Tabs defaultActiveKey="1">
      <Tabs.TabPane tab="Common" key="1">
        <Radio.Group onChange={commonChange}>
          <div className="radio-item">
            <Radio value="public">Public</Radio>
          </div>
          <div className="radio-item">
            <Radio value="download">Download</Radio>
          </div>
          <div className="radio-item">
            <Radio value="upload">Upload</Radio>
          </div>
          <div className="radio-item">
            <Radio value="private">Private</Radio>
          </div>
          <div className="btn-wrap">
            <Button type="primary">Save</Button>
          </div>
        </Radio.Group>
      </Tabs.TabPane>
      <Tabs.TabPane tab="Custom" key="2">
        <TextArea value={powerStore.json} rows={6} />
        <div className="btn-wrap">
          <Button type="primary" onClick={save}>Save</Button>
        </div>
      </Tabs.TabPane>
    </Tabs>
  </div>;
};

export default observer(Power);

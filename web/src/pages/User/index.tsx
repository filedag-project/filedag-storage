import { Axios, HttpMethods } from '@/api/https';
import { SignModel } from '@/models/SignModel';
import { ACCESS_KEY_ID, Cookies } from '@/utils/cookies';
import { Form, Input } from 'antd';
import { observer } from 'mobx-react';
import { useEffect } from 'react';
import styles from './style.module.scss';
const Dashboard = (props:any) => {
  const [form] = Form.useForm();
  useEffect(()=>{
    // 
  },[]);
  
  const changePassword = async ()=>{
    const newSecretKey = form.getFieldValue('password');
    const accessKey = Cookies.getKey(ACCESS_KEY_ID);
    const params:SignModel = {
      service: 's3',
      body: '',
      protocol: 'http',
      method: HttpMethods.post,
      applyChecksum: true,
      path:`/console/v1/change-password`,
      query:{
        newSecretKey:newSecretKey,
        accessKey:accessKey
      },
      region: '',
    }
    const res = await Axios.axiosXMLStream(params);
    console.log(res,'change psd');
  };

  return <div className={styles.dashboard}>
    <Form form={form} autoComplete="off" style={{width:'400px'}}>
      <Form.Item
          name="password"
          rules={[
              {required: true, message: 'Please input your password!'},
          ]}
      >
          <Input.Password
              placeholder="please enter your password"
          />
      </Form.Item>
      <Form.Item
          name="password"
          rules={[
              {required: true, message: 'Please input your password!'},
          ]}
      >
          <Input.Password
              placeholder="please enter your password"
          />
      </Form.Item>
    </Form>
  </div>;
};

export default observer(Dashboard);

import { Axios, HttpMethods } from '@/api/https';
import { SignModel } from '@/models/SignModel';
import { RouterPath } from '@/router/RouterConfig';
import { ACCESS_KEY_ID, Cookies } from '@/utils/cookies';
import { Button, Form, Input } from 'antd';
import { observer } from 'mobx-react';
import { useNavigate } from 'react-router-dom';
import styles from './style.module.scss';
const CreateUser = (props:any) => {
  const [form] = Form.useForm();
  const navigate = useNavigate();
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
    Axios.axiosXMLStream(params).then(res=>{
      navigate(RouterPath.buckets)
    })
  };

  return <div className={styles.user}>
    <div className={styles.title}>Change Password</div>
      <Form form={form} autoComplete="off">
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
        
        <Form.Item>
          <Button type="primary" onClick={changePassword}>
              Confirm
          </Button>
        </Form.Item>
      </Form>
    </div>;
};

export default observer(CreateUser);

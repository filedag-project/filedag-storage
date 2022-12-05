import { Axios, HttpMethods } from '@/api/https';
import { SignModel } from '@/models/SignModel';
import { RouterPath } from '@/router/RouterConfig';
import { ACCESS_KEY_ID, Cookies } from '@/utils/cookies';
import { Button, Form, Input } from 'antd';
import { observer } from 'mobx-react';
import { useHistory } from 'react-router';
import styles from './style.module.scss';
const ChangePassword = (props:any) => {
  const [form] = Form.useForm();
  const history = useHistory();
  const confirm = async ()=>{
    const oldSecretKey = form.getFieldValue('oldPassword')
    const newSecretKey = form.getFieldValue('newPassword');
    const accessKey = Cookies.getKey(ACCESS_KEY_ID);
    const params:SignModel = {
      service: 's3',
      body: '',
      protocol: 'http',
      method: HttpMethods.post,
      applyChecksum: true,
      path:`/console/v1/change-password`,
      query:{
        oldSecretKey:oldSecretKey,
        newSecretKey:newSecretKey,
        accessKey:accessKey
      },
      region: '',
    }
    Axios.axiosJson(params).then(res=>{
      history.push(RouterPath.dashboard)
    })
  };

  return <div className={styles.user}>
    <div className={styles.title}>Change Password</div>
      <Form form={form} autoComplete="off">
        <Form.Item
            name="oldPassword"
            rules={[
                {required: true, message: 'Please input your old password'},
            ]}
        >
            <Input.Password
                placeholder="please enter your old password"
            />
        </Form.Item>
        <Form.Item
            name="newPassword"
            rules={[
                {required: true, message: 'Please input your new password'},
            ]}
        >
            <Input.Password
                placeholder="please enter your new password"
            />
        </Form.Item>
        <Form.Item>
          <Button type="primary" onClick={confirm}>
              Confirm
          </Button>
        </Form.Item>
      </Form>
    </div>;
};

export default observer(ChangePassword);

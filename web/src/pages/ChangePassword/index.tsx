import { Axios, HttpMethods } from '@/api/https';
import { SignModel } from '@/models/SignModel';
import { RouterPath } from '@/router/RouterConfig';
import { Cookies, SESSION_TOKEN } from '@/utils/cookies';
import { Button, Form, Input } from 'antd';
import { observer } from 'mobx-react';
import { useHistory } from 'react-router';
import jwt_decode from "jwt-decode";
import styles from './style.module.scss';
import { tokenType } from '@/models/RouteModel';
const ChangePassword = (props:any) => {
  const [form] = Form.useForm();
  const history = useHistory();
  const confirm = async ()=>{
    const oldSecretKey = form.getFieldValue('oldPassword')
    const newSecretKey = form.getFieldValue('newPassword');
    const _jwt = Cookies.getKey(SESSION_TOKEN);
    const _token:tokenType = jwt_decode(_jwt);
    const {parent}=_token;

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
        accessKey:parent
      },
      region: '',
    }
    Axios.axiosJson(params).then(res=>{
      history.push(RouterPath.buckets)
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

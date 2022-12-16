import { Axios, HttpMethods } from '@/api/https';
import classNames from 'classnames';
import { SignModel } from '@/models/SignModel';
import { RouterPath } from '@/router/RouterConfig';
import { Cookies, SESSION_TOKEN } from '@/utils/cookies';
import { Button, Form, Input } from 'antd';
import { observer } from 'mobx-react';
import { useNavigate } from 'react-router-dom';
import jwt_decode from "jwt-decode";
import styles from './style.module.scss';
import { tokenType } from '@/models/RouteModel';
import iconPassword from '@/assets/images/login/icon-password.png';
import iconShow from '@/assets/images/login/icon-show.png';
import iconHidden from '@/assets/images/login/icon-hidden.png';
import { useState } from 'react';

const ChangePassword = (props:any) => {
  const [active,setActive]=useState(false);
  const [form] = Form.useForm();
  const navigate = useNavigate();
  const confirm = async ()=>{
    try{
      await form.validateFields();
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
        navigate(RouterPath.buckets)
      })
    }catch(eror){

    }
  };

  const oldChange = (e)=>{
    const _old = e.target.value;
    const _new = form.getFieldValue('newPassword');
    const _bool = Boolean(_old) && Boolean(_new) 
    setActive(_bool);
}
const newChange = (e)=>{
    const _old = form.getFieldValue('oldPassword')
    const _new = e.target.value;
    const _bool = Boolean(_old) && Boolean(_new) 
    setActive(_bool);
}

  return <div className={styles.user}>
    <div className={styles.title}>Change Password</div>
      <Form form={form} autoComplete="off">
        <Form.Item
            name="oldPassword"
            rules={[
                {required: true, message: 'Please input your old password'},
                { min:8, message: 'password must be at least 8 characters'},
            ]}
        >
            <Input.Password
                placeholder="please enter your old password"
                prefix={<img src={iconPassword} alt=''/>}
                iconRender={(visible)=>{
                    return visible ? <img src={iconHidden} alt=''/>:<img src={iconShow} alt=''/>
                }}
                onChange={oldChange}
            />
        </Form.Item>
        <Form.Item
            name="newPassword"
            rules={[
                {required: true, message: 'Please input your new password'},
                { min:8, message: 'password must be at least 8 characters'},
            ]}
        >
            <Input.Password
                placeholder="please enter your new password"
                prefix={<img src={iconPassword} alt=''/>}
                iconRender={(visible)=>{
                    return visible ? <img src={iconHidden} alt=''/>:<img src={iconShow} alt=''/>
                }}
                onChange={newChange}
            />
        </Form.Item>
        <Form.Item className='btn-wrap'>
          <Button type="primary" className={classNames(active ? 'active':'')} onClick={confirm}>
              Confirm
          </Button>
        </Form.Item>
      </Form>
    </div>;
};

export default observer(ChangePassword);

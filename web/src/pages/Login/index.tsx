import classNames from 'classnames';
import styles from './style.module.scss';
import {Button, Checkbox, Form, Input} from 'antd';
import logo from '@/assets/images/common/logo.png';
import {useNavigate} from 'react-router-dom';
import {RouterPath} from '@/router/RouterConfig';
import _ from 'lodash';
import { SignModel } from '@/models/SignModel';
import { Cookies,ACCESS_KEY_ID,SECRET_ACCESS_KEY,SESSION_TOKEN, USER_NAME } from '@/utils/cookies';
import { HttpMethods, Axios } from '@/api/https';
import { useEffect, useState } from 'react';
import loginBg from '@/assets/images/login/login-bg.png';
import iconUser from '@/assets/images/login/icon-user.png';
import iconPassword from '@/assets/images/login/icon-password.png';
import iconShow from '@/assets/images/login/icon-show.png';
import iconHidden from '@/assets/images/login/icon-hidden.png';

const Login = () => {
    const navigate = useNavigate();
    const [active,setActive]=useState(false);
    const [form] = Form.useForm();
    useEffect(()=>{
        const username = Cookies.getKey(USER_NAME);
        if(username){
            form.setFieldValue('username',username);
        }
    },[]);
    const nameChange = (e)=>{
        const name = e.target.value;
        const password = form.getFieldValue('password');
        const _bool = Boolean(name) && Boolean(password) 
        setActive(_bool);
    }
    const passwordChange = (e)=>{
        const name = form.getFieldValue('username')
        const password = e.target.value;
        const _bool = Boolean(name) && Boolean(password) 
        setActive(_bool);
    }
    const submitLogin = async () => {
        try {
            await form.validateFields();
            const body = 'Action=AssumeRole&DurationSeconds=86400&Version=2011-06-15';
            const remember = form.getFieldValue('remember');
            const _username = form.getFieldValue('username');
            const _password = form.getFieldValue('password');
            const params:SignModel = {
                service: 'sts',
                accessKeyId: _username,
                secretAccessKey: _password,
                sessionToken: '',
                body: body,
                protocol: 'http',
                method: HttpMethods.post,
                applyChecksum: false,
                path:'/',
                region: '',
                contentType:'application/x-www-form-urlencoded'
            }
            const res = await Axios.axiosXMLStream(params);
            const AccessKeyId = _.get(res,'AssumeRoleResponse.AssumeRoleResult.Credentials.AccessKeyId._text','');
            const SecretAccessKey = _.get(res,'AssumeRoleResponse.AssumeRoleResult.Credentials.SecretAccessKey._text','');
            const SessionToken = _.get(res,'AssumeRoleResponse.AssumeRoleResult.Credentials.SessionToken._text','');
            Cookies.setKey(ACCESS_KEY_ID,AccessKeyId);
            Cookies.setKey(SECRET_ACCESS_KEY,SecretAccessKey);
            Cookies.setKey(SESSION_TOKEN,SessionToken);
            if(remember){
                Cookies.setKey(USER_NAME,_username,7); 
            }
            navigate(RouterPath.home);
        } catch (e) {
            
        }
    };

    return (
        <div className={classNames(styles.login)}>
            <div className={classNames(styles.left)}>
                <div className={classNames(styles['logo-group'])}>
                    <img className={classNames(styles.img)} src={logo} alt=""/>
                </div>
                <div className={classNames(styles['login-bg'])}>
                    <img src={loginBg} alt="" />
                </div>
            </div>
            <div className={classNames(styles.right)}>
                <div className={classNames(styles.form)}>
                    <div className={classNames(styles.title)}>Login to FileDAG</div>
                    <Form form={form} autoComplete="off">
                        <Form.Item
                            name="username"
                            rules={[
                                {required: true, message: 'Please input your username!'},
                            ]}
                        >
                            <Input
                                prefix={<img src={iconUser} alt=''/>}
                                placeholder="please enter your username"
                                onChange={nameChange}
                            />
                        </Form.Item>
                        <Form.Item
                            name="password"
                            rules={[
                                {required: true, message: 'Please input your password!'},
                            ]}
                        >
                            <Input.Password
                                prefix={<img src={iconPassword} alt=''/>}
                                placeholder="please enter your password"
                                iconRender={(visible)=>{
                                    return visible ? <img src={iconHidden} alt=''/>:<img src={iconShow} alt=''/>
                                }}
                                onChange={passwordChange}
                            />
                        </Form.Item>
                        <Form.Item className='login-btn'>
                            <Button type="primary" className={classNames(active ? 'active':'')} onClick={submitLogin}>
                                Login
                            </Button>
                        </Form.Item>
                        <Form.Item className='remember-me' name="remember" valuePropName="checked">
                            <Checkbox>Remember me</Checkbox>
                        </Form.Item>
                    </Form>
                </div>
            </div>
        </div>
    );
};

export default Login;

import classNames from 'classnames';
import styles from './style.module.scss';
import {Button, Checkbox, Form, Input} from 'antd';
import {LockOutlined, UserOutlined} from '@ant-design/icons';
import logo from '@/assets/images/common/logo.png';
import {useHistory} from 'react-router';
import {RouterPath} from '@/router/RouterConfig';
import _ from 'lodash';
import { SignModel } from '@/models/SignModel';
import { Cookies,ACCESS_KEY_ID,SECRET_ACCESS_KEY,SESSION_TOKEN, USER_NAME } from '@/utils/cookies';
import { HttpMethods, Axios } from '@/api/https';
import { useEffect } from 'react';

const Login = () => {
    const history = useHistory();
    const [form] = Form.useForm();
    useEffect(()=>{
        const username = Cookies.getKey(USER_NAME);
        if(username){
            form.setFieldValue('username',username);
        }
    },[]);
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
                Cookies.setKey(USER_NAME,_username); 
            }else{
                Cookies.deleteKey(USER_NAME); 
            }
            history.push(RouterPath.home);
        } catch (e) {
            
        }
    };

    return (
        <div className={classNames(styles.login)}>
            <div className={classNames(styles.left)}>
                <div className={classNames(styles['logo-group'])}>
                    <img className={classNames(styles.img)} src={logo} alt=""/>
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
                                prefix={<UserOutlined/>}
                                placeholder="please enter your username"
                            />
                        </Form.Item>
                        <Form.Item
                            name="password"
                            rules={[
                                {required: true, message: 'Please input your password!'},
                            ]}
                        >
                            <Input.Password
                                prefix={<LockOutlined/>}
                                placeholder="please enter your password"
                            />
                        </Form.Item>
                        <Form.Item name="remember" valuePropName="checked">
                            <Checkbox>Remember me</Checkbox>
                        </Form.Item>
                        <Form.Item>
                            <Button type="primary" onClick={submitLogin}>
                                Login
                            </Button>
                        </Form.Item>
                    </Form>
                </div>
            </div>
        </div>
    );
};

export default Login;

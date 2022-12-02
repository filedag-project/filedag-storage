import { Redirect, Route, useHistory } from 'react-router-dom';
import { RouterPath } from '@/router/RouterConfig';
import classNames from 'classnames';
import { useState } from 'react';
import { MenuFoldOutlined, MenuUnfoldOutlined,LoginOutlined, FormOutlined } from '@ant-design/icons';
import styles from './style.module.scss';
import { Form, Input, Layout, Modal } from 'antd';
import PageMenu from '@/components/Menu';
import { ACCESS_KEY_ID, Cookies, SECRET_ACCESS_KEY, SESSION_TOKEN } from '@/utils/cookies';
import { SignModel } from '@/models/SignModel';
import { Axios, HttpMethods } from '@/api/https';
const { Sider } = Layout;
const IS_LOGGED = false;

const noSliderPage = ['/login'];
const PageLayout = (props: any) => {
  const [changePasswordShow,setChangePasswordShow]=useState(false);
  const history = useHistory();
  const { component: Com, auth, path, ...rest } = props;
  const [collapsed, setCollapsed] = useState(false);
  const [form] = Form.useForm();
  const triggerSlider = () => {
    setCollapsed(!collapsed);
  };
  const logOut = ()=>{
    Cookies.deleteKey(ACCESS_KEY_ID);
    Cookies.deleteKey(SECRET_ACCESS_KEY);
    Cookies.deleteKey(SESSION_TOKEN);
    history.push(RouterPath.login);
  };
  const openChangePassword = ()=>{
    setChangePasswordShow(true);
  }

  const changePassword = async ()=>{
    const newSecretKey = form.getFieldValue('password');
    const accessKey = Cookies.getKey(ACCESS_KEY_ID);
    const params:SignModel = {
      service: 's3',
      body: '',
      protocol: 'http',
      method: HttpMethods.post,
      applyChecksum: true,
      path:`console/v1/change-password`,
      query:{
        newSecretKey:newSecretKey,
        accessKey:accessKey
      },
      region: '',
    }
    const res = await Axios.axiosXMLStream(params);
    console.log(res,'change psd');
  };

  const cancelChange = ()=>{
    setChangePasswordShow(false);
  }

  const triggerNode = () => {
    return collapsed ? (
      <MenuUnfoldOutlined></MenuUnfoldOutlined>
    ) : (
      <MenuFoldOutlined></MenuFoldOutlined>
    );
  };
  return (
    <Route
      {...rest}
      render={(props: any) => {
        if (!IS_LOGGED) {
          if (noSliderPage.includes(path)) {
            return <Com {...props} />;
          } else {
            return (
              <div className={styles['layout']}>
                <Layout hasSider>
                  <Sider collapsed={collapsed}>
                    <div
                      className={classNames('logo', collapsed ? 'small' : '')}
                    >
                      FileDAG
                    </div>
                    <PageMenu></PageMenu>
                  </Sider>
                  <div className="layout-content">
                    <div className="layout-header">
                      <div className="trigger" onClick={triggerSlider}>
                        {triggerNode()}
                      </div>
                      <div className={classNames('action')}>
                        <div className={classNames('login-out')} onClick={logOut}>
                          <LoginOutlined /><span>Log Out</span>
                        </div>
                        <div className={classNames('change-password')} onClick={openChangePassword}>
                          <FormOutlined /> <span>Change password</span>
                        </div>
                      </div>
                      <Modal
                        title="Change Password"
                        open={changePasswordShow}
                        onOk={changePassword}
                        onCancel={cancelChange}
                        okText="Confirm"
                        cancelText="Cancel"
                      >
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
                      </Modal>
                    </div>
                    <Com {...props} />
                  </div>
                </Layout>
              </div>
            );
          }
        } else {
          return <Redirect to={RouterPath.login} />;
        }
      }}
    />
  );
};

export default PageLayout;

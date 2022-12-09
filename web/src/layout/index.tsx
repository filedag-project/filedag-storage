
import { Redirect, Route, useHistory } from 'react-router-dom';
import { RouterPath } from '@/router/RouterConfig';
import classNames from 'classnames';
import { useEffect, useState } from 'react';
import { MenuFoldOutlined, MenuUnfoldOutlined, UserOutlined } from '@ant-design/icons';
import styles from './style.module.scss';
import { Dropdown, Layout, MenuProps } from 'antd';
import PageMenu from '@/components/Menu';
import { ACCESS_KEY_ID, Cookies, SECRET_ACCESS_KEY, SESSION_TOKEN } from '@/utils/cookies';
import { observer } from 'mobx-react';
import { tokenType } from '@/models/RouteModel';
import jwt_decode from "jwt-decode";
const { Sider } = Layout;
const IS_LOGGED = false;

const noSliderPage = ['/login'];
const PageLayout = (props: any) => {
  const history = useHistory();
  const [name,setName]=useState('');
  const { component: Com, auth, path, ...rest } = props;
  const [collapsed, setCollapsed] = useState(false);

  useEffect(()=>{
    const _jwt = Cookies.getKey(SESSION_TOKEN);
    if(_jwt){
      const _token:tokenType = jwt_decode(_jwt);
      const {parent}=_token;
      setName(parent);
    }
  },[]);

  const items: MenuProps['items'] = [
    {
      label: 'Change Password',
      key: 'changePassword',
      onClick:()=>{
        history.push(RouterPath.changePassword);
      }
    },
    {
      label: 'Log out',
      key: 'logout',
      onClick:()=>{
        Cookies.deleteKey(ACCESS_KEY_ID);
        Cookies.deleteKey(SECRET_ACCESS_KEY);
        Cookies.deleteKey(SESSION_TOKEN);
        history.push(RouterPath.login);
      }
    },
  ]

  const triggerSlider = () => {
    setCollapsed(!collapsed);
  };

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
        console.log(123);
        
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
                      <Dropdown menu={{ items }} placement="bottomLeft">
                        <div className="user">
                          <UserOutlined />
                          <span>{name}</span>
                        </div>
                      </Dropdown>
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

export default observer(PageLayout);


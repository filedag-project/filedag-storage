import { Redirect, Route, useHistory } from 'react-router-dom';
import { RouterPath } from '@/router/RouterConfig';
import classNames from 'classnames';
import { useState } from 'react';
import { MenuFoldOutlined, MenuUnfoldOutlined,LoginOutlined } from '@ant-design/icons';
import styles from './style.module.scss';
import { Layout } from 'antd';
import PageMenu from '@/components/Menu';
import { ACCESS_KEY_ID, Cookies, SECRET_ACCESS_KEY, SESSION_TOKEN } from '@/utils/cookies';
const { Sider } = Layout;
const IS_LOGGED = false;

const noSliderPage = ['/login'];
const PageLayout = (props: any) => {
  const history = useHistory();
  const { component: Com, auth, path, ...rest } = props;
  const [collapsed, setCollapsed] = useState(false);
  const triggerSlider = () => {
    setCollapsed(!collapsed);
  };
  const logOut = ()=>{
    Cookies.deleteKey(ACCESS_KEY_ID);
    Cookies.deleteKey(SECRET_ACCESS_KEY);
    Cookies.deleteKey(SESSION_TOKEN);
    history.push(RouterPath.login);
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
                      <div className={classNames('login-out')} onClick={logOut}>
                        <LoginOutlined /><span>Log Out</span>
                      </div>
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

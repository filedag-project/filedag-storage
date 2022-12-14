
import { useNavigate,Outlet, useLocation } from 'react-router-dom';
import { RouterPath, RouterToBreadcrumb } from '@/router/RouterConfig';
import classNames from 'classnames';
import { useEffect, useState } from 'react';
import { MenuFoldOutlined, MenuUnfoldOutlined, UserOutlined } from '@ant-design/icons';
import styles from './style.module.scss';
import { Dropdown, Layout as ALayout, MenuProps } from 'antd';
import Menu from '@/components/Menu';
import { ACCESS_KEY_ID, Cookies, SECRET_ACCESS_KEY, SESSION_TOKEN } from '@/utils/cookies';
import { observer } from 'mobx-react';
import { tokenType } from '@/models/RouteModel';
import jwt_decode from "jwt-decode";
import iconLogo from '@/assets/images/common/icon-logo.png';
const { Sider } = ALayout;

const Layout = (props: any) => {
  const location = useLocation();
  const navigate = useNavigate();
  const [name,setName]=useState('');
  const [breadcrumb,setBreadcrumb] = useState([])
  const [collapsed, setCollapsed] = useState(false);
  const [items,setItems] = useState([
    {
      label: 'Log out',
      key: 'logout',
      onClick:()=>{
        Cookies.deleteKey(ACCESS_KEY_ID);
        Cookies.deleteKey(SECRET_ACCESS_KEY);
        Cookies.deleteKey(SESSION_TOKEN);
        navigate(RouterPath.login);
      }
    },
  ])
  useEffect(()=>{
    const _jwt = Cookies.getKey(SESSION_TOKEN);
    if(_jwt){
      const _token:tokenType = jwt_decode(_jwt);
      const {parent,isAdmin = false}=_token;
      setName(parent);
      if(!isAdmin){
        setItems([
          {
            label: 'Change Password',
            key: 'change-password',
            onClick:()=>{
              navigate(RouterPath.changePassword);
            }
          },{
            label: 'Log out',
            key: 'logout',
            onClick:()=>{
              Cookies.deleteKey(ACCESS_KEY_ID);
              Cookies.deleteKey(SECRET_ACCESS_KEY);
              Cookies.deleteKey(SESSION_TOKEN);
              navigate(RouterPath.login);
            }
          }
        ])
      }
    }else{
      navigate(RouterPath.login)
    }
    getBreadcrumb();
  },[location]);

  const getBreadcrumb = ()=>{
    const pathname = location.pathname;
    const _breadcrumb = RouterToBreadcrumb[pathname]??[];
    setBreadcrumb(_breadcrumb);
  }


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
    <div className={styles['layout']}>
      <ALayout hasSider>
        <Sider collapsed={collapsed}>
          <div
            className={classNames('logo', collapsed ? 'small' : '')}
          >
            <img src={iconLogo} alt="" />
            <span>FileDAG</span>
          </div>
          <Menu></Menu>
        </Sider>
        <div className="layout-content">
          <div className="layout-header">
            <div className="trigger" onClick={triggerSlider}>
              {triggerNode()}
            </div>
            <div className="breadcrumb">
              {
                breadcrumb.map((n,index)=>{
                  return (
                    <span key={'breadcrumb'+index} className="breadcrumb-item">
                      <span className='breadcrumb-item-name' onClick={()=>{
                        if(index+1 === breadcrumb.length) return;
                        navigate(n['path'])
                      }}>
                        {n['label']}
                      </span>
                      <span className='breadcrumb-item-to'>{index+1 === breadcrumb.length ? '':'>'}</span>
                    </span>
                  )
                })
              }
            </div>
            <Dropdown menu={{ items }} placement="bottomLeft">
              <div className="user">
                <UserOutlined />
                <span>{name}</span>
              </div>
            </Dropdown>
          </div>
          <Outlet />
        </div>
      </ALayout>
    </div>
  );
};

export default observer(Layout);


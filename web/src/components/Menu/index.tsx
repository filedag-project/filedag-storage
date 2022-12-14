import { Menu as AMenu, MenuProps } from 'antd';
import {
  AppstoreOutlined,
  UserOutlined,
  DashboardOutlined
} from '@ant-design/icons';
import { RouterPath } from '@/router/RouterConfig';
import { useNavigate, useLocation } from 'react-router-dom';
import { useEffect, useState } from 'react';
import { Cookies, SESSION_TOKEN } from '@/utils/cookies';
import { tokenType } from '@/models/RouteModel';
import jwt_decode from "jwt-decode";

const items: MenuProps['items'] = [
  {
    icon: <DashboardOutlined />,
    label: 'Overview',
    key: 'overview',
  },
  { icon: <AppstoreOutlined />, label: 'Buckets', key: 'buckets'},
];

const adminItems: MenuProps['items'] = [
  {
    icon: <DashboardOutlined />,
    label: 'Dashboard',
    key: 'dashboard',
  },
  { icon: <AppstoreOutlined />, label: 'Buckets', key: 'buckets'},
  {
    icon: <UserOutlined />,
    label: 'User',
    key: 'user',
  },
]

const Menu = () => {
  const [selectedKey,setSelectedKey] = useState('')
  const navigate = useNavigate();
  const location = useLocation();
  const _default:MenuProps['items'] = []
  const [myMenu,setMyMenu] = useState(_default);
  useEffect(()=>{
    const { pathname } = location;
    setSelectedKey(pathname.replace("/",""));
    if(pathname.startsWith('/objects') || pathname.startsWith('/power')){
      setSelectedKey('buckets')
    }
    const _jwt = Cookies.getKey(SESSION_TOKEN);
    if(_jwt){
      const _token:tokenType = jwt_decode(_jwt);
      const { isAdmin } = _token;
      isAdmin ? setMyMenu(adminItems) : setMyMenu(items);
    }
  },[location])
  
  const menuClick = (e)=>{
    const key = e.key;
    navigate(RouterPath[key]);
  }
  return <AMenu selectedKeys={[selectedKey]} theme={'dark'} mode="inline" items={myMenu} onClick={menuClick}/>;
};

export default Menu;


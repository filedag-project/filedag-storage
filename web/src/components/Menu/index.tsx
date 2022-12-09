import { Menu, MenuProps } from 'antd';
import {
  AppstoreOutlined,
  UserOutlined,
  DashboardOutlined
} from '@ant-design/icons';
import { RouterPath } from '@/router/RouterConfig';
import { useHistory, useLocation } from 'react-router';
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

const PageMenu = () => {
  const [selectedKey,setSelectedKey] = useState('')
  const history = useHistory();
  const location = useLocation();
  const _default:MenuProps['items'] = []
  const [myMenu,setMyMenu] = useState(_default);
  useEffect(()=>{
    const { pathname } = location;
    setSelectedKey(pathname.replace("/",""));
    const _buckets = ['/power','/objects'];
    if(_buckets.includes(pathname)){
      setSelectedKey('buckets')
    }
    const _jwt = Cookies.getKey(SESSION_TOKEN);
    const _token:tokenType = jwt_decode(_jwt);
    const { isAdmin } = _token;
    isAdmin ? setMyMenu(adminItems) : setMyMenu(items);
  },[location])
  
  const menuClick = (e)=>{
    const key = e.key;
    history.push(RouterPath[key]);
  }
  return <Menu selectedKeys={[selectedKey]} theme={'dark'} mode="inline" items={myMenu} onClick={menuClick}/>;
};

export default PageMenu;


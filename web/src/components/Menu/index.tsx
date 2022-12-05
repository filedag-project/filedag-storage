import { Menu } from 'antd';
import {
  AppstoreOutlined,
  UserOutlined,
  DashboardOutlined
} from '@ant-design/icons';
import { RouterPath } from '@/router/RouterConfig';
import { useHistory, useLocation } from 'react-router';
import { useEffect, useState } from 'react';

const items = [
  {
    icon: <DashboardOutlined />,
    label: 'Dashboard',
    key: 'dashboard',
  },
  { icon: <AppstoreOutlined />, label: 'Buckets', key: 'buckets'},
  { icon: <UserOutlined />, label: 'Overview', key: 'overview' },
];



const PageMenu = () => {
  const [selectedKey,setSelectedKey] = useState('')
  const history = useHistory();
  const location = useLocation();
  useEffect(()=>{
    const { pathname } = location;
    setSelectedKey(pathname.replace("/",""));
    const _buckets = ['/power','/objects']
    if(_buckets.includes(pathname)){
      setSelectedKey('buckets')
    }
  },[location])
  
  const menuClick = (e)=>{
    const key = e.key;
    history.push(RouterPath[key]);
  }
  return <Menu selectedKeys={[selectedKey]} theme={'dark'} mode="inline" items={items} onClick={menuClick}/>;
};

export default PageMenu;

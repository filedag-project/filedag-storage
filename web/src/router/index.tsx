import Layout from '@/layout';
import { lazy, ReactNode, Suspense } from 'react';
const Login = lazy(()=> import('@/pages/Login'));
const Home = lazy(()=> import('@/pages/Home'));
const Dashboard = lazy(()=> import('@/pages/Dashboard'));
const Overview = lazy(()=> import('@/pages/Overview'));
const Buckets = lazy(()=> import('@/pages/Buckets'));
const CreateBucket = lazy(()=> import('@/pages/CreateBucket'));
const Objects = lazy(()=> import('@/pages/Objects'));
const User = lazy(()=> import('@/pages/User'));
const Power = lazy(()=> import('@/pages/Power'));
const ChangePassword = lazy(()=> import('@/pages/ChangePassword'));


const lazyLoad = (children:ReactNode):ReactNode => {
  return <Suspense fallback={<>loading...</>}>
    {children}
  </Suspense>
}

const router = [
  {
    path:'/',
    element:<Layout />,
    children:[
      { path:'', element:lazyLoad(<Home />) },
      { path:'home', element:lazyLoad(<Home />) },
      { path :'dashboard', element:lazyLoad(<Dashboard />) },
      { path :'overview', element:lazyLoad(<Overview />) },
      { path :'buckets', element:lazyLoad(<Buckets />) },
      { path :'power', element:lazyLoad(<Power />) },
      { path :'objects', element:lazyLoad(<Objects />) },
      { path :'create-bucket', element:lazyLoad(<CreateBucket />) },
      { path :'user', element:lazyLoad(<User />) },
      { path :'change-password', element:lazyLoad(<ChangePassword />) },
    ]
  },
  { path:'/login', element:lazyLoad(<Login />) },
]

export default router;
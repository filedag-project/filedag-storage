import { lazy } from 'react';
import { RouteModel } from '@/models/RouteModel';

enum RouterPath {
  objects = '/objects',
  error = '/error',
  login = '/login',
  buckets = '/buckets',
  createBucket = '/create-bucket',
  dashboard = '/dashboard',
  overview = '/overview',
  power = '/power',
  changePassword = '/changePassword',
  createUser = '/createUser',
  home = '/home',
  user = '/user'
}
// all roles can be accessed
const Routes: RouteModel[] = [
  {
    name: 'home',
    path: '/',
    auth: true,
    component: lazy(() => import('@/pages/Home')),
  },
  {
    name: 'home',
    path: RouterPath.home,
    auth: true,
    component: lazy(() => import('@/pages/Home')),
  },
  {
    name: 'dashboard',
    path: RouterPath.dashboard,
    auth: true,
    component: lazy(() => import('@/pages/Dashboard')),
  },
  {
    name: 'user',
    path: RouterPath.user,
    auth: true,
    component: lazy(() => import('@/pages/User')),
  },
  {
    name: 'overview',
    path: RouterPath.overview,
    auth: true,
    component: lazy(() => import('@/pages/Overview')),
  },
  {
    name: 'buckets',
    path: RouterPath.buckets,
    auth: true,
    component: lazy(() => import('@/pages/Buckets')),
  },
  {
    name: 'createBucket',
    path: RouterPath.createBucket,
    auth: true,
    component: lazy(() => import('@/pages/CreateBucket')),
  },
  {
    name: 'objects',
    path: RouterPath.objects,
    auth: true,
    component: lazy(() => import('@/pages/Objects')),
  },
  {
    name: 'power',
    path: RouterPath.power,
    auth: true,
    component: lazy(() => import('@/pages/Power')),
  },
  {
    name: 'changePassword',
    path: RouterPath.changePassword,
    auth: true,
    component: lazy(() => import('@/pages/ChangePassword')),
  },
  {
    name: 'createUser',
    path: RouterPath.createUser,
    auth: true,
    component: lazy(() => import('@/pages/CreateUser')),
  },
  {
    name: 'login',
    path: RouterPath.login,
    auth: false,
    component: lazy(() => import('@/pages/Login')),
  },
  {
    name: '404',
    path: RouterPath.error,
    auth: false,
    component: lazy(() => import('@/404')),
  },
];

export { RouterPath,Routes };

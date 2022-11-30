import { lazy } from 'react';
import { RouteModel } from '@/models/RouteModel';

enum RouterPath {
  objects = '/objects',
  error = '/error',
  login = '/login',
  buckets = '/buckets',
  createBucket = '/create-bucket',
  dashboard = '/dashboard',
  user = '/user',
  power = '/power'
}

const Routes: RouteModel[] = [
  {
    name: 'dashboard',
    path: '/',
    auth: false,
    component: lazy(() => import('@/pages/Dashboard')),
  },
  {
    name: 'buckets',
    path: RouterPath.buckets,
    auth: false,
    component: lazy(() => import('@/pages/Buckets')),
  },
  {
    name: 'createBucket',
    path: RouterPath.createBucket,
    auth: false,
    component: lazy(() => import('@/pages/CreateBucket')),
  },
  {
    name: 'objects',
    path: RouterPath.objects,
    auth: true,
    component: lazy(() => import('@/pages/Objects')),
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
    name: 'power',
    path: RouterPath.power,
    auth: true,
    component: lazy(() => import('@/pages/Power')),
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

export { RouterPath, Routes };

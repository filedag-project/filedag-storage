enum RouterPath {
  objects = '/objects',
  error = '/error',
  login = '/login',
  buckets = '/buckets',
  createBucket = '/create-bucket',
  dashboard = '/dashboard',
  overview = '/overview',
  power = '/power',
  changePassword = '/change-password',
  home = '/home',
  user = '/user'
}
const RouterToBreadcrumb = {
  '/buckets':[{ path: RouterPath.buckets,label:'Buckets'}],
  '/objects':[{ path: RouterPath.buckets,label:'Buckets'},{ path: RouterPath.objects,label:'Objects'}],
  '/dashboard':[{ path: RouterPath.dashboard,label:'Dashboard'}],
  '/change-password':[{ path: RouterPath.changePassword,label:'Change-password'}],
  '/user':[{ path: RouterPath.user,label:'User'}],
  '/power':[{ path: RouterPath.buckets,label:'Buckets'},{ path: RouterPath.power,label:'Power'}],
  '/create-bucket':[{ path: RouterPath.createBucket,label:'Create-bucket'}],
  '/overview':[{ path: RouterPath.overview,label:'Overview'}]
}

export { RouterPath,RouterToBreadcrumb };

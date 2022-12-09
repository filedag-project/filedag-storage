import { RouterPath } from '@/router/RouterConfig';
import { Cookies, SESSION_TOKEN } from '@/utils/cookies';
import { observer } from 'mobx-react';
import { Redirect } from 'react-router';
import jwt_decode from "jwt-decode";
import { tokenType } from '@/models/RouteModel';
const Home = (props:any) => {
  const _jwt = Cookies.getKey(SESSION_TOKEN);
  const _token:tokenType = jwt_decode(_jwt);
  const { isAdmin } = _token;
  return <Redirect to={isAdmin ? RouterPath.dashboard : RouterPath.overview} />;
};

export default observer(Home);

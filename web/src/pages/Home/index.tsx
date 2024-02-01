import { RouterPath } from '@/router/RouterConfig';
import { Cookies, SESSION_TOKEN } from '@/utils/cookies';
import { Navigate } from 'react-router-dom';
import jwt_decode from "jwt-decode";
import { tokenType } from '@/models/RouteModel';
const Home = (props:any) => {
  const _jwt = Cookies.getKey(SESSION_TOKEN);
  if(_jwt){
    const _token:tokenType = jwt_decode(_jwt);
    const { isAdmin } = _token;
    return <Navigate to={isAdmin ? RouterPath.dashboard : RouterPath.overview} />;
  }else{
    return <></>;
  }
};

export default Home;

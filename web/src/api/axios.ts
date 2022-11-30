import axios from 'axios';

const service = axios.create({
  withCredentials: false,
  timeout: 35000,
});

service.interceptors.request.use(
  (request) => {
    let headers = request.headers;
    request.headers = { ...headers };
    return request;
  },
  (error) => {},
);

service.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    return error;
  },
);

export default service;

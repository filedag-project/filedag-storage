import React from 'react';
import ReactDOM from 'react-dom';
import 'antd/dist/antd.min.css';
import './styles/var.scss';
import './styles/reset.scss';
import './index.css';
import { Provider } from 'mobx-react';
import * as stores from './store';
import App from './App';
import { BrowserRouter } from 'react-router-dom';
ReactDOM.render(
  <React.StrictMode>
    <Provider {...stores}>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </Provider>
  </React.StrictMode>,
  document.getElementById('root'),
);

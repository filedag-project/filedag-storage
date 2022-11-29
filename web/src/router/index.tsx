import { Suspense } from 'react';
import { BrowserRouter, Redirect, Route, Switch } from 'react-router-dom';

import PageLayout from '@/layout';
import { RouterPath, Routes } from '@/router/RouterConfig';

const MRouter = () => {
  return (
    <Suspense fallback={<div></div>}>
      <BrowserRouter>
        <Switch>
          {Routes.map((item, index) => {
            return (
              <PageLayout
                exact
                path={item.path}
                component={item.component}
                auth={item.auth}
                key={'layout' + index}
              ></PageLayout>
            );
          })}
          <Route path={'/*'}>
            <Redirect to={RouterPath.error}></Redirect>
          </Route>
        </Switch>
      </BrowserRouter>
    </Suspense>
  );
};
export default MRouter;

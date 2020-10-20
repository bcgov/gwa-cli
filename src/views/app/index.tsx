import * as React from 'react';
import { Box, useApp, useInput } from 'ink';
import { Route, Switch, useLocation, useHistory } from 'react-router';

import AppBar from './app-bar';
import { appState } from '../../state/app';
import AppContext from '../../services/context';
import ConfigOrg from '../config-org';
import { IAppContext } from '../../types';
import Plugins from '../plugins';
import Review from '../review';
import Start from '../start';

interface AppProps {
  args: IAppContext;
}

const App: React.FC<AppProps> = ({ args }) => {
  const [state, setAppState] = appState.use();
  const history = useHistory();
  const location = useLocation();
  const { exit } = useApp();

  useInput((input, key) => {
    if (input === 'q' && key.ctrl) {
      exit();
    }
    if (state.mode === 'view') {
      if (input === 'l' && key.ctrl) {
        history.push(location.pathname.replace(/\/([a-z0-9_-]*[\/]?)$/, ''));
      } else if (key.rightArrow) {
        history.goForward();
      } else if (key.leftArrow) {
        history.goBack();
      } else if (input === 'P') {
        history.push('/export');
      }
    } else {
      if (key.escape) {
        setAppState((prev) => ({
          ...prev,
          mode: 'view',
        }));
      }
    }
  });

  return (
    <AppContext.Provider value={args}>
      <Box flexDirection="column" width="100%">
        <Route path="/:any" component={AppBar} />
        <Box>
          <Switch>
            <Route exact path="/" component={Start} />
            <Route exact path="/org">
              <ConfigOrg />
            </Route>
            <Route path="/editor" component={Plugins} />
            <Route exact path="/export">
              <Review />
            </Route>
          </Switch>
        </Box>
      </Box>
    </AppContext.Provider>
  );
};

export default App;

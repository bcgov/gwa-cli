import * as React from 'react';
import { Box, Text, useApp, useInput } from 'ink';
import { Route, Switch, useLocation, useHistory } from 'react-router';

import AppBar from './app-bar';
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
  const history = useHistory();
  const location = useLocation();
  const { exit } = useApp();

  useInput((input, key) => {
    if (input === 'q' && key.ctrl) {
      exit();
    } else if (input === 'l' && key.ctrl) {
      history.push(location.pathname.replace(/\/([a-z0-9_-]*[\/]?)$/, ''));
    } else if (input === 'k' && key.ctrl) {
      history.goForward();
    } else if (input === 'j' && key.ctrl) {
      history.goBack();
    } else if (input === 'y' && key.ctrl) {
      history.push('/export');
    }
  });

  return (
    <AppContext.Provider value={args}>
      <Box flexDirection="column">
        <Route
          path="/:any"
          render={(props) => <AppBar {...props} file={args.file} />}
        />
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

import * as React from 'react';
import { Box, Text, useApp, useInput } from 'ink';
import { Route, Switch, useLocation, useHistory } from 'react-router';

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
    } else if (input === 'n' && key.ctrl) {
      history.goForward();
    } else if (input === 'p' && key.ctrl) {
      history.goBack();
    }
  });

  return (
    <AppContext.Provider value={args}>
      <Box flexDirection="column">
        <Box>
          <Switch>
            <Route exact path="/" component={Start} />
            <Route exact path="/org">
              <ConfigOrg />
            </Route>
            <Route path="/editor" component={Plugins} />
            <Route exact path="/review">
              <Review />
            </Route>
          </Switch>
        </Box>
      </Box>
    </AppContext.Provider>
  );
};

export default App;

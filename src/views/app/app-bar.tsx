import React, { useRef, useEffect, useState } from 'react';
import { Box, measureElement, Text, useStdout } from 'ink';
import { Route, Switch, RouteComponentProps } from 'react-router';

import { orgState } from '../../state/org';
import PluginStatus from '../plugins/status';

interface AppBarProps extends RouteComponentProps {}

const AppBar: React.FC<AppBarProps> = ({ match }) => {
  const { stdout } = useStdout();
  // const ref = useRef(null);
  const [fill, setFill] = useState<number>(0);
  const { name, file } = orgState.useValue();

  /* stdout.on('resize', () => {
   *   const { width } = measureElement(ref.current);
   *   setFill(width);
   * });

   * useEffect(() => {
   *   const { width } = measureElement(ref.current);
   *   setFill(width);
   * }, [match, name, setFill]); */

  return (
    <Box flexGrow={1} width="100%" justifyContent="space-between">
      <Box>
        <Text inverse bold color="cyan">
          {' GWA Config'}
        </Text>
        <Box>
          <Switch>
            <Route path="/org">
              <Text inverse color="cyan">
                {'/Settings '}
              </Text>
            </Route>
            <Route path="/editor">
              <Text inverse color="cyan">
                {'/Plugins '}
              </Text>
            </Route>
            <Route path="/export">
              <Text inverse color="cyan">
                {'/Export '}
              </Text>
            </Route>
          </Switch>
          <Route exact path="/editor/:plugin" component={PluginStatus} />
        </Box>
      </Box>
      <Box>
        <Text inverse>
          {` ${name || '! [Service not configured]'} ${
            file ? `[${file}] ` : ' [+] New configuration '
          }`}
        </Text>
        <Text inverse color="green">
          {'[P] Publish '}
        </Text>
      </Box>
    </Box>
  );
};

export default AppBar;

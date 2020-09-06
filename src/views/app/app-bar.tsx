import React, { useRef, useEffect, useState } from 'react';
import { Box, measureElement, Text, useStdout } from 'ink';
import { Route, Switch, RouteComponentProps } from 'react-router';

import { orgState } from '../../state/org';

interface AppBarProps extends RouteComponentProps {
  file: string | null;
}

const AppBar: React.FC<AppBarProps> = ({ file, match }) => {
  const { stdout } = useStdout();
  const ref = useRef(null);
  const [fill, setFill] = useState<number>(0);
  const { name } = orgState.useValue();

  stdout.on('resize', () => {
    const { width } = measureElement(ref.current);
    setFill(width);
  });

  useEffect(() => {
    const { width } = measureElement(ref.current);
    setFill(width);
  }, [match, name, setFill]);

  return (
    <Box width="100%">
      <Box>
        <Text inverse bold color="cyan">
          {' '}
          GWA Config{' '}
        </Text>
        <Box>
          <Switch>
            <Route path="/org">
              <Text inverse color="white">
                {' Team Settings '}
              </Text>
            </Route>
            <Route path="/editor">
              <Text inverse color="white">
                {' Plugins '}
              </Text>
            </Route>
            <Route path="/export">
              <Text inverse color="white">
                {' Export '}
              </Text>
            </Route>
          </Switch>
        </Box>
      </Box>
      <Box ref={ref} flexGrow={1}>
        <Text inverse color="grey">
          {' '.repeat(fill)}
        </Text>
      </Box>
      <Box>
        <Text inverse>{` ${name || '[Service not configured]'} `}</Text>
        <Text inverse color="gray">
          {`[${file}]` || 'New configuration'}
        </Text>
      </Box>
    </Box>
  );
};

export default AppBar;

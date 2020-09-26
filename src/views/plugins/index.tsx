import React, { useState } from 'react';
import { Box, Text, useInput } from 'ink';
import { Route } from 'react-router';
import { Tab, Tabs } from 'ink-tab';
import { match, useHistory, useLocation } from 'react-router';

import { appState } from '../../state/app';
import { IPlugin } from '../../types';
import { activePluginState, pluginsState } from '../../state/plugins';
import PluginsList from './list';
import PluginEditor from '../plugin-editor';
import StepHeader from '../../components/step-header';
import PluginItem from './item';

interface PluginsProps {
  match: match;
  onComplete: (item: any) => void;
}

const Plugins: React.FC<PluginsProps> = ({ match }) => {
  const [app, setAppState] = appState.use();
  const history = useHistory();
  const { pathname } = useLocation();
  const [index, setIndex] = useState<number>(0);
  const state = pluginsState.useValue();
  const plugins: IPlugin[] = Object.values(state);
  const urls = plugins.map((plugin) => `/editor/${plugin.id}`);

  useInput((input) => {
    const idx = urls.indexOf(pathname);

    if (app.mode === 'view') {
      if (input === 'n') {
        const nextIndex = idx + 1;
        const nextUrlCandidate = urls[nextIndex];
        const nextUrl = nextUrlCandidate || urls[0];
        history.push(nextUrl);
      } else if (input === 'p') {
        const prevIndex = idx - 1;
        const prevUrl = prevIndex < 0 ? urls.slice(-1)[0] : urls[prevIndex];
        history.push(prevUrl);
      }
    }
  });

  return (
    <Box flexDirection="column">
      <Box flexDirection="column" width="100%">
        <Route
          exact
          path={match.url}
          render={(props) => (
            <PluginsList
              {...props}
              data={plugins}
              index={index}
              onChange={setIndex}
            />
          )}
        />
      </Box>
      <Route exact path={`${match.url}/:plugin`} component={PluginEditor} />
    </Box>
  );
};

export default Plugins;

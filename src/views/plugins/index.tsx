import React, { useState } from 'react';
import { Box, Text, useInput } from 'ink';
import { Route } from 'react-router';
import { Tab, Tabs } from 'ink-tab';
import { match, useHistory } from 'react-router';

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
  const [index, setIndex] = useState<number>(0);
  const state = pluginsState.useValue();
  const plugins: IPlugin[] = Object.values(state);

  return (
    <Box flexDirection="column">
      <StepHeader title="Configure Your Organization" />
      <Box flexDirection="column">
        <Box marginBottom={1}>
          <Text bold underline>
            Plugins
          </Text>
        </Box>
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

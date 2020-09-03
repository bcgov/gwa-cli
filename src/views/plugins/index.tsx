import * as React from 'react';
import { Box, Text, useInput } from 'ink';
import { Route } from 'react-router';
import { Tab, Tabs } from 'ink-tab';
import { match, useHistory } from 'react-router';

import { activePluginState, pluginsState } from '../../state/plugins';
import PluginEditor from '../plugin-editor';
// import StepHeader from '../../components/step-header';
import PluginItem from './item';

interface PluginsProps {
  match: match;
  onComplete: (item: any) => void;
}

const Plugins: React.FC<PluginsProps> = ({ match }) => {
  const history = useHistory();
  const plugins = pluginsState.useValue();
  const [selectedPlugin, setSelectedPlugin] = activePluginState.use();
  const pluginNames: string[] = Object.keys(plugins);
  const tabs = pluginNames.map((plugin) => ({
    url: `/${plugin}`,
    name: plugin,
  }));
  const onChange = (name: string) => history.push(match.url + name);

  return (
    <Box>
      <Box marginRight={4} flexDirection="column">
        <Box marginBottom={1}>
          <Text bold underline>
            Plugins
          </Text>
        </Box>
        <Tabs flexDirection="column" onChange={onChange}>
          {tabs.map((tab: any) => (
            <Tab key={tab.url} name={tab.url}>
              {tab.name}
            </Tab>
          ))}
        </Tabs>
      </Box>
      <Box flexDirection="column" marginLeft={5}>
        <Box marginBottom={1}>
          <Text bold underline>
            Configure
          </Text>
        </Box>
        <Route exact path={`${match.url}/:plugin`} component={PluginEditor} />
      </Box>
    </Box>
  );
};

export default Plugins;

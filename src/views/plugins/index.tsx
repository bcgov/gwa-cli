import * as React from 'react';
import { Box, Text } from 'ink';
import SelectInput, { Item } from 'ink-select-input';

import { activePluginState, pluginsState } from '../../state/plugins';
import PluginEditor from '../plugin-editor';
import StepHeader from '../../components/step-header';
import PluginItem from './item';

interface PluginsProps {
  onComplete: (item: any) => void;
  step: number;
}

const Plugins = ({ step }: PluginsProps) => {
  const plugins = pluginsState.useValue();
  const [selectedPlugin, setSelectedPlugin] = activePluginState.use();
  const pluginNames = Object.keys(plugins);
  const activePlugin = pluginNames[selectedPlugin];
  const items = pluginNames.map((plugin) => ({
    id: plugin,
    label: plugin,
    value: plugin,
  }));
  const onHighlight = (item: any) => setSelectedPlugin(item.id);

  return (
    <Box flexDirection="column">
      <StepHeader step={step} title="Select & Configure Plugins" />
      <Box>
        <Box marginRight={4} flexDirection="column">
          <Box marginBottom={1}>
            <Text bold underline>
              Plugins
            </Text>
          </Box>
          <SelectInput
            items={items}
            itemComponent={PluginItem}
            onHighlight={onHighlight}
          />
        </Box>
        <Box flexDirection="column">
          <Box marginBottom={1}>
            <Text bold underline>
              {`Configure ${selectedPlugin} Plugin`}
            </Text>
          </Box>
          <PluginEditor selected={selectedPlugin} />
        </Box>
      </Box>
    </Box>
  );
};

export default Plugins;

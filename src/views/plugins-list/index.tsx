import React from 'react';
import { Box, Text, render } from 'ink';
import { uid } from 'react-uid';

import { usePluginsState } from '../../state/plugins';

const PluginsList: React.FC = () => {
  const plugins = usePluginsState();
  const elements = [];

  for (const plugin in plugins) {
    const data = plugins[plugin];

    elements.push(
      <Box key={uid(plugin)} flexDirection="column" paddingX={4}>
        <Box justifyContent="space-between">
          <Text>
            <Text bold>{data.meta.name}</Text>
            {data.meta.bcgov && '*'}
          </Text>
          <Text dimColor underline>
            {data.meta.url}
          </Text>
        </Box>
        <Box>
          <Text>{data.meta.description}</Text>
        </Box>
      </Box>
    );
  }

  return (
    <Box flexDirection="column" width="100%">
      <Box justifyContent="space-between" marginBottom={1}>
        <Text>GWA Plugins</Text>
        <Text italic>* denotes BC Gov plugin</Text>
      </Box>
      {elements}
    </Box>
  );
};

export default function () {
  render(<PluginsList />);
}

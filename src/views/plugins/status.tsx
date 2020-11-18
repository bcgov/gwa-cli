import React from 'react';
import { Box, Text } from 'ink';
import { RouteComponentProps } from 'react-router';

import { usePluginsState } from '../../state/plugins';

interface PluginStatusProps extends RouteComponentProps<{ plugin: string }> {}

const PluginStatus: React.FC<PluginStatusProps> = ({ match }) => {
  const plugins = usePluginsState();
  const plugin: any = plugins[match.params.plugin];

  return (
    <Box>
      <Text bold inverse color="blueBright">
        {` ${plugin.meta.name} `}
      </Text>
      <Text
        inverse
        color={plugin.meta.enabled ? 'greenBright' : 'magentaBright'}
      >
        {plugin.meta.enabled ? ' [e] Enabled ' : ' [e] Disabled '}
      </Text>
    </Box>
  );
};

export default PluginStatus;

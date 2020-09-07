import React from 'react';
import { Box, Text } from 'ink';
import { RouteComponentProps } from 'react-router';

import { IPlugin } from '../../types';
import { pluginsState } from '../../state/plugins';

interface PluginStatusProps extends RouteComponentProps<{ plugin: string }> {}

const PluginStatus: React.FC<PluginStatusProps> = ({ match }) => {
  const plugins = pluginsState.useValue();
  const plugin: IPlugin = plugins[match.params.plugin];

  return (
    <Box>
      <Text bold inverse color="blueBright">
        {` ${plugin.name} `}
      </Text>
      <Text
        inverse
        color={plugin.data.enabled ? 'greenBright' : 'magentaBright'}
      >
        {plugin.data.enabled ? ' Enabled ' : ' Disabled '}
      </Text>
    </Box>
  );
};

export default PluginStatus;

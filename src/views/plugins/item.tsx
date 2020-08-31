import React from 'react';
import { Box, Text } from 'ink';
import { ItemProps } from 'ink-select-input';

const PluginItem: React.FC<ItemProps> = ({ label }) => (
  <Box>
    <Text>{label}</Text>
  </Box>
);

export default PluginItem;

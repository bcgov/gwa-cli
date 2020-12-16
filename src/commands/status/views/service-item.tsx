import * as React from 'react';
import { Box, Text } from 'ink';

import type { StatusData } from '../types';

interface ServiceItemProps {
  data: StatusData;
}
const ServiceItem: React.FC<ServiceItemProps> = ({ data }) => {
  const isUp = data.status === 'UP';
  const textColor = isUp ? 'greenBright' : 'redBright';
  let hostText = `${data.envHost} [${data.upstream}]`;

  if (data.envHost === data.upstream) {
    hostText = data.envHost;
  }

  return (
    <Box width="100%">
      <Box width={2}>
        <Text color={textColor}>{isUp ? '▲' : '▼'}</Text>
      </Box>
      <Box width={18}>
        <Text color={textColor}>{data.name}</Text>
      </Box>
      <Box width={30}>
        <Text>{data.reason}</Text>
      </Box>
      <Box>
        <Text dimColor>{hostText}</Text>
      </Box>
    </Box>
  );
};

export default ServiceItem;

import * as React from 'react';
import { Box, Text } from 'ink';

import type { StatusData } from '../types';

interface ServiceItemProps {
  data: StatusData;
  sm: boolean;
}
const ServiceItem: React.FC<ServiceItemProps> = ({ data, sm }) => {
  const isUp = data.status === 'UP';
  const textColor = isUp ? 'greenBright' : 'redBright';

  return (
    <Box width="100%">
      <Box width={2}>
        <Text color={textColor}>{isUp ? '▲' : '▼'}</Text>
      </Box>
      <Box width={45}>
        <Text color={textColor} wrap="truncate">
          {data.name}
        </Text>
      </Box>
      <Box width={15}>
        <Text>{data.reason}</Text>
      </Box>
      {!sm && (
        <Box>
          <Text dimColor wrap="truncate-middle">
            {data.upstream}
          </Text>
        </Box>
      )}
    </Box>
  );
};

export default ServiceItem;

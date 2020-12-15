import React from 'react';
import { Box, Text } from 'ink';
import { uid } from 'react-uid';

import api from '../../services/api';
import config from '../../config';
import useAsync from '../../hooks/use-async';

interface StatusData {
  name: string;
  upstream: string;
  status: 'UP' | 'DOWN';
  reason: string;
  host: string;
  envHost: string;
}

interface StatusViewProps {}

const StatusView: React.FC<StatusViewProps> = () => {
  const { namespace } = config();
  const data = useAsync<StatusData[]>(api, '/namespaces/:namespace/services', {
    namespace: 'sampler',
  });
  const d: StatusData[] = [
    {
      name: 'sampler',
      upstream: 'api.host.com',
      status: 'UP',
      reason: '200 Response',
      host: 'sampler.api',
      envHost: 'sampler.api',
    },
    {
      name: 'sampler',
      upstream: 'api.host.com',
      status: 'DOWN',
      reason: '500 status code',
      host: 'sampler.api',
      envHost: 'sampler.api',
    },
  ];

  return (
    <Box flexDirection="column" width="100%">
      <Box marginY={1}>
        <Text>{`${namespace} Status`}</Text>
      </Box>
      {d.map((service) => (
        <Box key={uid(service)} width="100%">
          <Box width={2}>
            <Text color={service.status === 'UP' ? 'greenBright' : 'redBright'}>
              {service.status === 'UP' ? '▲' : '▼'}
            </Text>
          </Box>
          <Box width="20%">
            <Text color={service.status === 'UP' ? 'greenBright' : 'redBright'}>
              {service.name}
            </Text>
          </Box>
          <Box flexGrow={1}>
            <Text>{service.reason}</Text>
          </Box>
          <Box width="30%">
            <Text dimColor>{`${service.host} / ${service.upstream}`}</Text>
          </Box>
        </Box>
      ))}
    </Box>
  );
};

export default StatusView;

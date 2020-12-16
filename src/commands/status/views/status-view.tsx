import * as React from 'react';
import { Box, Text } from 'ink';
import { uid } from 'react-uid';

import api from '../../../services/api';
import config from '../../../config';
import useAsync from '../../../hooks/use-async';
import ServiceItem from './service-item';
import type { StatusData } from '../types';

interface StatusViewProps {}

const StatusView: React.FC<StatusViewProps> = () => {
  const { namespace } = config();
  const data = useAsync<StatusData[]>(api, '/namespaces/:namespace/services', {
    namespace,
  });

  React.useEffect(() => {
    const downServices = data.filter((d: StatusData) => d.status === 'DOWN');

    if (downServices.length > 0) {
      process.exitCode = 1;
    }
  }, [data]);

  return (
    <Box flexDirection="column" width="100%">
      <Box marginY={1}>
        <Text>{`${namespace} Status`}</Text>
      </Box>
      {!data.length && (
        <Box>
          <Text>You have no services yet.</Text>
        </Box>
      )}
      {data.map((service: StatusData) => (
        <ServiceItem key={uid(service)} data={service} />
      ))}
    </Box>
  );
};

export default StatusView;

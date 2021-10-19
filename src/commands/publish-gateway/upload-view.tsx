import React from 'react';
import { Box, Text } from 'ink';

import publish, { bundleFiles } from '../../services/publish';
import Success from '../../components/success';
import useAsync from '../../hooks/use-async';

interface UploadViewProps {
  options: {
    configFile?: string;
    dryRun: string;
  };
}

const UploadView: React.FC<UploadViewProps> = ({ options }) => {
  const formData = useAsync(async () => {
    const value = await bundleFiles(options.configFile);

    return {
      configFile: {
        value,
        options: {
          filename: options.configFile,
          contentType: null,
        },
      },
      dryRun: options.dryRun,
    };
  }, [bundleFiles, options]);
  const { message, results } = useAsync(
    publish,
    '/namespaces/:namespace/gateway',
    formData
  );

  return (
    <Box flexDirection="column">
      <Success>
        <Text color="green">Success</Text>
        {` Configuration ${options.configFile} Published`}
      </Success>
      <Box flexDirection="column" marginTop={1}>
        <Box>
          <Text>{message}</Text>
        </Box>
        <Box marginTop={1}>
          <Text>{results}</Text>
        </Box>
      </Box>
    </Box>
  );
};

export default UploadView;

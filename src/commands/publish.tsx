import React, { Suspense } from 'react';
import { Box, Text, render } from 'ink';
import { getToken, publish } from '../services/gwa';
import { ErrorBoundary } from 'react-error-boundary';

import Failed from '../components/failed';
import Loading from '../components/loading';
import Success from '../components/success';
import makeRequest from '../hooks/use-request';

type PublishResponse = {
  message: string;
  results: string;
};
const useApi = makeRequest<PublishResponse>();

interface PublishProps {
  input: string;
  options: any;
}

const Publish: React.FC<PublishProps> = ({ input, options }) => {
  const json = useApi(async () => {
    const token = await getToken();
    return await publish({
      configFile: input,
      env: options.env,
      dryRun: Boolean(options.dryRun).toString(),
      token,
    });
  });

  return (
    <Box flexDirection="column">
      <Success>
        <Text color="green">Success</Text> {`Configuration ${input} Published`}
      </Success>
      <Box flexDirection="column" marginTop={1}>
        <Box>
          <Text>{json.message}</Text>
        </Box>
        <Box marginTop={1}>
          <Text>{json.results}</Text>
        </Box>
      </Box>
    </Box>
  );
};

export default function app(input: string, options: any) {
  render(
    <ErrorBoundary FallbackComponent={Failed}>
      <Suspense fallback={<Loading>Uploading config...</Loading>}>
        <Publish input={input} options={options} />
      </Suspense>
    </ErrorBoundary>
  );
}

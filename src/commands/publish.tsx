import React from 'react';
import { Box, Text, render } from 'ink';
import { getToken, publish } from '../services/gwa';

import HttpRequest from '../components/http-request';

export default async function (input: string, options: any) {
  const { rerender } = render(
    <HttpRequest loading loadingText="Uploading config..." />
  );
  try {
    const token = await getToken();
    const json = await publish({
      token,
      configFile: input,
      dryRun: Boolean(options.dryRun).toString(),
    });

    rerender(
      <Box flexDirection="column">
        <HttpRequest successText="Config Published" />
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
  } catch (err) {
    console.error('Upload Failed');
    console.error(err);
  }
}

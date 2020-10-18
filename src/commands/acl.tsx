import React, { Suspense } from 'react';
import { Box, Text, render } from 'ink';
import { addMembers } from '../services/gwa';
import { ErrorBoundary } from 'react-error-boundary';

import Failed from '../components/failed';
import Loading from '../components/loading';
import Success from '../components/success';
import makeRequest from '../hooks/use-request';

const useApi = makeRequest<ACLResponse>();

type ACLResponse = {
  added: number;
  missing: number;
  removed: number;
};

interface ACLProps {
  options: {
    env: string;
    users: string;
  };
}

const ACL = ({ options }: ACLProps) => {
  const json = useApi(
    async () =>
      await addMembers({
        env: options.env,
        users: options.users,
      })
  );

  return (
    <Box>
      <Box flexDirection="column">
        <Box flexDirection="column">
          <Success>
            <Text color="green">Success</Text> Membership Updated
          </Success>
          <Box marginTop={1}>
            <Box marginRight={1}>
              <Text bold color="green">
                +
              </Text>
            </Box>
            <Text>{`${json.added} Added`}</Text>
          </Box>
          <Box>
            <Box marginRight={1}>
              <Text bold color="red">
                -
              </Text>
            </Box>
            <Text>{`${json.removed} Removed`}</Text>
          </Box>
          <Box>
            <Box marginRight={1}>
              <Text bold color="yellow">
                ?
              </Text>
            </Box>
            <Text>{`${json.missing} Missing`}</Text>
          </Box>
        </Box>
      </Box>
    </Box>
  );
};

export default function acl(input: string, options: any) {
  render(
    <ErrorBoundary FallbackComponent={Failed}>
      <Suspense fallback={<Loading>Publishing membership changes...</Loading>}>
        <ACL options={options} />
      </Suspense>
    </ErrorBoundary>
  );
}

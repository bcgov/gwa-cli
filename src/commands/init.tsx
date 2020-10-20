import React, { Suspense } from 'react';
import { Box, Text, render } from 'ink';
import { ErrorBoundary } from 'react-error-boundary';
import fs from 'fs';

import Failed from '../components/failed';
import Loading from '../components/loading';
import Success from '../components/success';
import makeRequest from '../hooks/use-request';

const useMakeEnv = makeRequest<string>();

type InitOptions = {
  namespace: string;
  devClientId: string | undefined;
  devClientSecret: string | undefined;
  testClientId: string | undefined;
  testClientSecret: string | undefined;
  prodClientId: string | undefined;
  prodClientSecret: string | undefined;
};

function makeEnvFile(options: InitOptions): Promise<string> {
  return new Promise((resolve, reject) => {
    fs.exists('.env', (exists) => {
      if (exists) {
        return reject(
          new Error(
            'You already have initiated a GWA workspace in this dir. You can edit the .env file to make changes'
          )
        );
      } else {
        const data = `GWA_NAMESPACE=${options.namespace}
DEV_CLIENT_ID=${options.devClientId || ''}
DEV_CLIENT_SECRET=${options.devClientSecret || ''}
TEST_CLIENT_ID=${options.testClientId || ''}
TEST_CLIENT_SECRET=${options.testClientSecret || ''}
PROD_CLIENT_ID=${options.prodClientId || ''}
PROD_CLIENT_SECRET=${options.prodClientSecret || ''}
`;
        fs.writeFile('.env', data, (err) => {
          if (err) {
            reject(new Error(`Unable to write file ${err}`));
          }
          resolve('.env file successfully generated');
        });
      }
    });
  });
}

interface InitProps {
  options: InitOptions;
}

const Init: React.FC<InitProps> = ({ options }) => {
  const text = useMakeEnv(async () => await makeEnvFile(options));

  return (
    <Success>
      <Text>{text}</Text>
    </Success>
  );
};

export default function init(input: string, options: InitOptions) {
  render(
    <ErrorBoundary FallbackComponent={Failed}>
      <Suspense fallback={<Loading>Uploading config...</Loading>}>
        <Init options={options} />
      </Suspense>
    </ErrorBoundary>
  );
}

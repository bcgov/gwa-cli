import React, { Suspense } from 'react';
import { ErrorBoundary } from 'react-error-boundary';
import fs from 'fs';
import isEmpty from 'lodash/isEmpty';
import pick from 'lodash/pick';
import { Text, render } from 'ink';

import Failed from '../components/failed';
import Loading from '../components/loading';
import Success from '../components/success';
import makeRequest from '../hooks/use-request';
import type { Envs } from '../types';
import PromptForm from '../components/prompt-form';
import { Prompt } from '../components/prompt-form/types';

const useMakeEnv = makeRequest<string>();

const prompts: Prompt[] = [
  {
    label: 'Namespace',
    key: 'namespace',
    constraint: {
      presence: { allowEmpty: false },
      length: { minimum: 5, maximum: 10 },
      format: {
        pattern: '[a-z0-9]+',
        flags: 'i',
        message: 'can only contain a-z and 0-9',
      },
    },
  },
  {
    label: 'Client ID',
    key: 'clientId',
    constraint: {
      presence: { allowEmpty: false },
    },
  },
  {
    label: 'Client Secret',
    key: 'clientSecret',
    secret: true,
    constraint: {
      presence: { allowEmpty: false },
    },
  },
];

type InitOptions = {
  namespace: string;
  clientId: string;
  clientSecret: string;
  env: Envs;
  debug: boolean;
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
        if (!options.namespace) {
          return reject(new Error('--namespace is required'));
        }
        const envArgs = pick(options, ['dev', 'test', 'prod']);
        const env = Object.keys(envArgs)[0] ?? 'test';
        const data = `GWA_NAMESPACE=${options.namespace}
CLIENT_ID=${options.clientId ?? ''}
CLIENT_SECRET=${options.clientSecret ?? ''}
GWA_ENV=${env}
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

export default function init(_: string, options: InitOptions) {
  if (isEmpty(pick(options, ['namespace', 'clientId', 'clientSecret']))) {
    render(
      <PromptForm
        options={prompts}
        title="Configure this folder's environment variables"
      >
        {({ data }) => (
          <ErrorBoundary
            fallbackRender={({ error }) => (
              <Failed error={error} verbose={options.debug} />
            )}
          >
            <Suspense fallback={<Loading>Writing .env file</Loading>}>
              <Init options={data} />
            </Suspense>
          </ErrorBoundary>
        )}
      </PromptForm>
    );
  } else {
    render(
      <ErrorBoundary
        fallbackRender={({ error }) => (
          <Failed error={error} verbose={options.debug} />
        )}
      >
        <Suspense fallback={<Loading>Writing .env file</Loading>}>
          <Init options={options} />
        </Suspense>
      </ErrorBoundary>
    );
  }
}

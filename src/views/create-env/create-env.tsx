import React, { Suspense, useCallback, useState } from 'react';
import { Box, render } from 'ink';

import { ErrorBoundary } from 'react-error-boundary';
import Loading from '../../components/loading';
import Failed from '../../components/failed';
import PromptForm, { Prompt } from '../../components/prompt-form';
import WriteEnvAction from './write-env-action';
import type { InitOptions } from '../../types';

interface CreateEnvViewProps {
  env: string;
  prompts: Prompt[];
}

const CreateEnvView: React.FC<CreateEnvViewProps> = ({ env, prompts }) => {
  const [data, setData] = useState<InitOptions | null>(null);
  const onSubmit = useCallback((formData: Omit<InitOptions, 'env'>) => {
    setData({
      ...formData,
      env,
    });
  }, []);

  return (
    <Box flexDirection="column">
      <PromptForm
        options={prompts}
        onSubmit={onSubmit}
        title="Configure this folder's environment variables"
      />
      {data && (
        <ErrorBoundary
          fallbackRender={({ error }) => (
            <Failed error={error} verbose={false} />
          )}
        >
          <Suspense fallback={<Loading>Writing .env file</Loading>}>
            <WriteEnvAction data={data} />
          </Suspense>
        </ErrorBoundary>
      )}
    </Box>
  );
};

export const renderView = (prompts: Prompt[], env: string) => {
  return render(<CreateEnvView env={env} prompts={prompts} />);
};

export default CreateEnvView;

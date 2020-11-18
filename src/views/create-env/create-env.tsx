import React, { useCallback, useState } from 'react';
import { Box } from 'ink';

import PromptForm, { Prompt } from '../../components/prompt-form';
import WriteEnvAction from './write-env-action';
import type { InitOptions } from '../../types';
import AsyncAction from '../../components/async-action';

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
        <AsyncAction loadingText="Writing file...">
          <WriteEnvAction data={data} />
        </AsyncAction>
      )}
    </Box>
  );
};

export default CreateEnvView;

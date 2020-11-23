import React, { useCallback, useState } from 'react';
import { Box } from 'ink';

import AsyncAction from '../../../components/async-action';
import PromptForm, { Prompt } from '../../../components/prompt-form';
import type { ConfigFormData } from '../types';
import WriteConfigAction from './write-config-action';

interface CreateEnvViewProps {
  prompts: Prompt[];
}

const CreateEnvView: React.FC<CreateEnvViewProps> = ({ prompts }) => {
  const [data, setData] = useState<ConfigFormData | null>(null);
  const onSubmit = useCallback((formData: ConfigFormData) => {
    setData(formData);
  }, []);

  return (
    <Box flexDirection="column">
      <PromptForm
        options={prompts}
        onSubmit={onSubmit}
        title="Create a new configuration file"
      />
      {data && (
        <AsyncAction loadingText="Writing config file...">
          <WriteConfigAction data={data} />
        </AsyncAction>
      )}
    </Box>
  );
};

export default CreateEnvView;

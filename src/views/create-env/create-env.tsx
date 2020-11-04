import React, { useState } from 'react';
import { Box, render } from 'ink';

import PromptForm, { Prompt } from '../../components/prompt-form';
import Success from '../../components/success';
import { makeEnvFile } from '../../services/app';

interface CreateEnvViewProps {
  env: string;
  prompts: Prompt[];
}

const CreateEnvView: React.FC<CreateEnvViewProps> = ({ env, prompts }) => {
  const [doneText, setDoneText] = useState<string>('');
  const onSubmit = async (data: any) => {
    try {
      const result = await makeEnvFile({
        ...data,
        env,
      });
      setDoneText(doneText);
    } catch (err) {
      throw err;
    }
  };

  return (
    <Box>
      <PromptForm
        options={prompts}
        onSubmit={onSubmit}
        title="Configure this folder's environment variables"
      />
      {Boolean(doneText) && <Success>{doneText}</Success>}
    </Box>
  );
};

export const renderView = (prompts: Prompt[], env: string) => {
  return render(<CreateEnvView env={env} prompts={prompts} />);
};

export default CreateEnvView;

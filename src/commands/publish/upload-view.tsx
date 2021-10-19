import React from 'react';
import { Box, Text } from 'ink';
import { resolve } from 'path';
import fs from 'fs';
import YAML from 'yaml';

import { publish } from '../../services/publish';
import Success from '../../components/success';
import useAsync from '../../hooks/use-async';

interface UploadViewProps {
  action: string;
  options: {
    body: string;
    content?: string;
  };
}

const UploadView: React.FC<UploadViewProps> = ({ action, options }) => {
  const payload = useAsync(async () => {
    try {
      const filePath = resolve(process.cwd(), options.body);
      const file = await fs.promises.readFile(filePath, 'utf8');
      const result = YAML.parse(file);

      if (options.content) {
        const contentPath = resolve(process.cwd(), options.content);
        const content = await fs.promises.readFile(contentPath, 'utf8');
        result.content = content;
      }
      return result;
    } catch (err) {
      throw new Error(err);
    }
  });
  const { result, status } = useAsync(
    publish,
    `/ds/api/namespaces/:namespace/${action}s`,
    payload
  );

  if (status !== 200) {
    return (
      <Box>
        <Box flexDirection="column">
          <Box>
            <Text bold color="red">
              x Error
            </Text>
            <Box marginLeft={1}>
              <Text>{result}</Text>
            </Box>
          </Box>
        </Box>
      </Box>
    );
  }

  return (
    <Box flexDirection="column">
      <Success>
        <Text color="green">Success</Text>
        {` ${action} published`}
      </Success>
      <Box flexDirection="column" marginTop={1}>
        <Box>
          <Text>{`Result: ${result}`}</Text>
        </Box>
      </Box>
    </Box>
  );
};

export default UploadView;

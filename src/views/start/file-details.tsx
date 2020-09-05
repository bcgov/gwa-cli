import React from 'react';
import { Box, Text, Newline } from 'ink';

import { orgState } from '../../state/org';

interface FileDetailsProps {
  file: string;
}

const FileDetails: React.FC<FileDetailsProps> = ({ file }) => {
  const org = orgState.useValue();

  return (
    <Box
      paddingX={3}
      paddingY={1}
      flexDirection="column"
      marginBottom={1}
      borderColor="green"
      borderStyle="bold"
    >
      <Box>
        <Box width={10} marginRight={2} justifyContent="flex-end">
          <Text bold>Group:</Text>
        </Box>
        <Text>{org.name}</Text>
      </Box>
      <Box>
        <Box width={10} marginRight={2} justifyContent="flex-end">
          <Text bold>Spec:</Text>
        </Box>
        <Text>{org.specUrl}</Text>
      </Box>
      <Box>
        <Box width={10} marginRight={2} justifyContent="flex-end">
          <Text bold>Config:</Text>
        </Box>
        <Box>
          <Text>
            <Text bold color="green">
              {`✓ ${file} `}
            </Text>
            <Newline />
            Hit [Enter] to edit.
          </Text>
        </Box>
      </Box>
    </Box>
  );
};

export default FileDetails;

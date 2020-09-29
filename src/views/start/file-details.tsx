import React from 'react';
import { Box, Text, Newline } from 'ink';

import { useTeamState } from '../../state/team';

interface FileDetailsProps {
  file: string | unknown;
}

const FileDetails: React.FC<FileDetailsProps> = ({ file }) => {
  const { name, team, host } = useTeamState();

  return (
    <Box paddingX={3} flexDirection="column" marginBottom={1}>
      <Box>
        <Box width={10} marginRight={2} justifyContent="flex-end">
          <Text bold>Name:</Text>
        </Box>
        <Text>{name}</Text>
      </Box>
      <Box>
        <Box width={10} marginRight={2} justifyContent="flex-end">
          <Text bold>Team:</Text>
        </Box>
        <Text>{team}</Text>
      </Box>
      <Box>
        <Box width={10} marginRight={2} justifyContent="flex-end">
          <Text bold>Host:</Text>
        </Box>
        <Text>{host}</Text>
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
          </Text>
        </Box>
      </Box>
    </Box>
  );
};

export default FileDetails;

import React from 'react';
import { render, Static, Box, Text } from 'ink';
import { uid } from 'react-uid';

const Help: React.FC = () => {
  const commands = ['init', 'validate', 'edit'];
  return (
    <Box flexDirection="column">
      <Box marginBottom={1}>
        <Text>GWA Help</Text>
      </Box>
      <Box marginBottom={1}>
        <Text>Usage: </Text>
        <Text>
          {'gwa <command>'} <Text italic>{'<file> '}</Text>
        </Text>
        <Text>--options</Text>
      </Box>
      <Box marginBottom={1}>
        <Text bold underline>
          Commands:
        </Text>
      </Box>
      {commands.map((d) => (
        <Box key={uid(d)} marginBottom={3}>
          <Box>
            <Box marginRight={2} width={20}>
              <Text>{d}</Text>
            </Box>
            <Box>
              <Text>
                Lorem ipsum dolor sit, amet consectetur adipisicing elit. Illum
                rerum sit vitae voluptatum nemo quibusdam molestias odio,
                maiores deserunt, tempora neque ipsum minima itaque nobis hic
                nam odit reiciendis. Beatae.
              </Text>
            </Box>
          </Box>
        </Box>
      ))}
    </Box>
  );
};

export default () => render(<Help />);

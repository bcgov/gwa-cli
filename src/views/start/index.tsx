import React from 'react';
import BigText from 'ink-big-text';
import { Box, Text } from 'ink';
import SelectInput, { Item } from 'ink-select-input';
import { useHistory } from 'react-router';

const StartView: React.FC<{}> = () => {
  const { push } = useHistory();
  const items = [
    {
      label: 'Create organization',
      value: '/org',
    },
    {
      label: 'Plugin Editor',
      value: '/editor',
    },
  ];
  const onSelect = (item: any) => {
    push(item.value);
  };

  return (
    <Box width="100%" justifyContent="center">
      <Box flexDirection="column" justifyContent="center">
        <Box
          borderStyle="round"
          borderColor="cyan"
          justifyContent="center"
          paddingX={1}
          marginBottom={2}
        >
          <BigText
            colors={['cyanBright', 'cyan', 'yellow']}
            font="chrome"
            text="GWA Config"
          />
        </Box>
        <Box alignItems="center" flexDirection="column">
          <Text bold>API Gateway Config</Text>
          <Text>Version 1.0.0</Text>
        </Box>
        <Box justifyContent="center" marginY={2}>
          <SelectInput items={items} onSelect={onSelect} />
        </Box>
        <Box justifyContent="center">
          <Text>Help/Legend coming soon</Text>
        </Box>
      </Box>
    </Box>
  );
};

export default StartView;

import React, { useContext } from 'react';
import BigText from 'ink-big-text';
import { Box, Text } from 'ink';
import SelectInput, { Item } from 'ink-select-input';
import { useHistory } from 'react-router';

import AppContext from '../../services/context';
import FileDetails from './file-details';

const StartView: React.FC<{}> = () => {
  const { file } = useContext(AppContext);
  const { push } = useHistory();
  const items = [
    {
      label: 'Configure Group',
      value: '/org',
      enabled: !file,
    },
    {
      label: 'Plugin Editor',
      value: '/editor',
      enabled: true,
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
        <Box alignItems="center" flexDirection="column" marginY={2}>
          {file && <FileDetails file={file} />}
          <SelectInput
            items={items.filter((d) => d.enabled)}
            onSelect={onSelect}
          />
        </Box>
        <Box justifyContent="center">
          <Text>Help/Legend coming soon</Text>
        </Box>
      </Box>
    </Box>
  );
};

export default StartView;

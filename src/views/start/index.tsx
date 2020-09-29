import React, { useContext } from 'react';
import BigText from 'ink-big-text';
import { Box, Text } from 'ink';
import isEmpty from 'lodash/isEmpty';
import SelectInput, { Item } from 'ink-select-input';
import { useHistory } from 'react-router';

import { useAppState } from '../../state/app';
import FileDetails from './file-details';

const StartView: React.FC<{}> = () => {
  const input = useAppState((state) => state.input);
  const { push } = useHistory();
  const items = [
    {
      label: 'Configure Group',
      value: '/org',
      enabled: !input,
    },
    {
      label: 'Plugin Editor',
      value: '/editor',
      enabled: true,
    },
    {
      label: 'Export',
      value: '/export',
      enabled: !!input,
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
          {!isEmpty(input) && <FileDetails file={input} />}
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

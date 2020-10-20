import React from 'react';
import { Box, Text, useInput } from 'ink';
import { IPlugin } from '../../types';
import { match, useHistory } from 'react-router';

interface PluginsListProps {
  data: IPlugin[];
  index: number;
  match: match;
  onChange: (index: number) => void;
}

const PluginsList: React.FC<PluginsListProps> = ({
  data,
  index,
  match,
  onChange,
}) => {
  const history = useHistory();

  useInput((input, key) => {
    if (input === 'j' || key.downArrow) {
      onChange(Math.min(data.length - 1, index + 1));
    } else if (input === 'k' || key.upArrow) {
      onChange(Math.max(0, index - 1));
    }

    if (key.return) {
      history.push(`${match.url}/${data[index].meta.id}`);
    }
  });

  return (
    <Box margin={1} flexDirection="column">
      {data.map((plugin: IPlugin, pluginIdx: number) => (
        <Box key={plugin.meta.name} flexDirection="column" marginBottom={1}>
          <Box>
            <Box width={1}>
              {index === pluginIdx && <Text color="yellowBright">▋</Text>}
            </Box>
            <Box marginRight={1} width={4} justifyContent="flex-end">
              <Text>{`${(pluginIdx + 1).toString()}.`}</Text>
            </Box>
            <Box>
              <Text
                bold
                inverse={pluginIdx === index}
                color={
                  plugin.meta.enabled || pluginIdx === index ? 'white' : 'grey'
                }
              >
                {plugin.meta.name}
              </Text>
            </Box>
            <Box marginLeft={2}>
              {plugin.meta.enabled ? (
                <Text color="greenBright">[Enabled]</Text>
              ) : (
                <Text dimColor color="red">
                  [Disabled]
                </Text>
              )}
            </Box>
          </Box>
          <Box marginLeft={6}>
            <Text
              bold={false}
              color={
                plugin.meta.enabled || pluginIdx === index ? 'white' : 'grey'
              }
            >
              {plugin.meta.description}
            </Text>
          </Box>
        </Box>
      ))}
    </Box>
  );
};

export default PluginsList;

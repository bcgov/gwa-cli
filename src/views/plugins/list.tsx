import React from 'react';
import { Box, Text, useInput } from 'ink';
import { match, useHistory } from 'react-router';
import type { PluginObject } from '../../types';

interface PluginsListProps {
  data: PluginObject[];
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
    <Box marginY={1} flexDirection="column">
      <Box marginBottom={1}>
        <Box width={3} />
        <Box width={4} justifyContent="center">
          <Text underline>✓</Text>
        </Box>
        <Box width={40}>
          <Text bold underline>
            Plugin
          </Text>
        </Box>
        <Box width={20}>
          <Text bold underline>
            Author
          </Text>
        </Box>
      </Box>
      {data.map((plugin: PluginObject, pluginIdx: number) => (
        <Box key={plugin.meta.name}>
          <Box width={3} justifyContent="center">
            <Text>{pluginIdx === index ? '>' : ''}</Text>
          </Box>
          <Box width={3} marginRight={1} justifyContent="center">
            {plugin.meta.enabled ? (
              <Text color="greenBright">✓</Text>
            ) : (
              <Text dimColor color="yellow">
                -
              </Text>
            )}
          </Box>
          <Box width={40}>
            <Text
              inverse={pluginIdx === index}
              color={plugin.meta.enabled ? 'green' : 'white'}
            >
              {plugin.meta.name}
            </Text>
          </Box>
          <Box width={20}>
            <Text color="yellow">{plugin.meta.author}</Text>
          </Box>
        </Box>
      ))}
    </Box>
  );
};

export default PluginsList;

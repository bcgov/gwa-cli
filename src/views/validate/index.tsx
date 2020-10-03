import React from 'react';
import { Box, Text, render } from 'ink';
import { uid } from 'react-uid';

interface ValidateProps {
  errors: any[];
}

const Validate: React.FC<ValidateProps> = ({ errors }) => {
  if (errors.length === 0) {
    return (
      <Box>
        <Box>
          <Text bold color="green">
            {'✓ '}
          </Text>
          <Text>File is valid</Text>
        </Box>
      </Box>
    );
  }

  return (
    <Box flexDirection="column">
      {errors.map((err) => (
        <Box key={uid(err)} marginY={2} flexDirection="column">
          <Box marginBottom={1}>
            <Text dimColor underline>
              {`Plugin: ${err.plugin} [${
                Object.keys(err.error).length
              } incorrect fields]`}
            </Text>
          </Box>
          <Box flexDirection="column">
            {Object.keys(err.error).map((key) => (
              <Box key={uid(key)}>
                <Box marginRight={2}>
                  <Text bold color="red">
                    {`× ${key}`}
                  </Text>
                </Box>
                <Box>
                  {err.error[key].map((text: string) => (
                    <Text key={uid(text)}>{text}</Text>
                  ))}
                </Box>
              </Box>
            ))}
          </Box>
        </Box>
      ))}
      <Box>
        <Text
          bold
          color="red"
        >{`There were ${errors.length} error(s) found`}</Text>
      </Box>
    </Box>
  );
};

export default (errors: any[]) => render(<Validate errors={errors} />);

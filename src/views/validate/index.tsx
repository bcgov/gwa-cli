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
              {`Plugin: ${err.plugin} [${err.error.details.length} incorrect fields]`}
            </Text>
          </Box>
          <Box flexDirection="column">
            {err.error.details.map((detail: any) => (
              <Box key={uid(detail)}>
                <Box marginRight={2}>
                  <Text bold color="red">
                    {`× ${detail.message}`}
                  </Text>
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

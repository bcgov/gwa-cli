import React, { Fragment } from 'react';
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
              {err.plugin}
            </Text>
          </Box>
          <Box>
            {Object.keys(err.error).map((key) => (
              <Fragment key={uid(key)}>
                <Box marginX={4}>
                  <Text color="red">{key}</Text>
                </Box>
                <Box>
                  {err.error[key].map((text) => (
                    <Text key={uid(text)}>{text}</Text>
                  ))}
                </Box>
              </Fragment>
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

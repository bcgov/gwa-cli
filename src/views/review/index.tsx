import React, { useContext, useState } from 'react';
import { Box, Newline, Text, useInput } from 'ink';
import { useHistory } from 'react-router';

import AppContext from '../../services/context';
import { buildSpec } from '../../services/kong';
import { orgState } from '../../state/org';

const Review: React.FC = () => {
  const history = useHistory();
  const [done, setDone] = useState<boolean>(false);
  const [error, setError] = useState<boolean>(false);
  const { dir, file } = useContext(AppContext);
  const org = orgState.useValue();
  const fileName = file || org.file;

  useInput((input, key) => {
    if (key.return) {
      if (error) {
        history.goBack();
      } else {
        try {
          buildSpec(dir, fileName);
          setDone(true);
        } catch (err) {
          setError(true);
        }
      }
    }
  });

  return (
    <Box
      marginY={3}
      flexDirection="column"
      justifyContent="center"
      alignItems="center"
      width="100%"
    >
      {error && (
        <Box flexDirection="column" alignItems="center">
          <Text bold color="red">
            Export Failed
          </Text>
          <Text>[ Enter ] to go back</Text>
        </Box>
      )}
      {!error && (
        <>
          <Text bold>Export</Text>
          <Newline />
          <Text>
            Hit [ENTER] to export your changes to{' '}
            <Text inverse color="green">{`${fileName || 'spec.yaml'}`}</Text>.
            {done && (
              <Text bold color="green">
                {' '}
                Saved!
              </Text>
            )}
          </Text>
          <Newline />
          <Text>Press [ ctrl + c ] to exit</Text>
        </>
      )}
    </Box>
  );
};

export default Review;

import React, { useContext, useEffect, useState } from 'react';
import { Box, Newline, Text, useInput } from 'ink';
import { useHistory } from 'react-router';

import AppContext from '../../services/context';
import { buildSpec } from '../../services/kong';
import { orgState } from '../../state/org';

const Review: React.FC = () => {
  const history = useHistory();
  const [error, setError] = useState<boolean>(false);
  const { dir, file } = useContext(AppContext);
  const org = orgState.useValue();
  const fileName = file || org.file;

  useInput((input, key) => {
    if (error && key.return) {
      history.goBack();
    }
  });

  useEffect(() => {
    try {
      buildSpec(dir, fileName);
    } catch (err) {
      setError(true);
    }
  }, [dir, file]);

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
          <Text bold>All Done!</Text>
          <Newline />
          <Text>
            Your Kong config file has been generated in{' '}
            <Text inverse color="green">{`${fileName || 'spec.yaml'}`}</Text>.
          </Text>
          <Newline />
          <Text>Commit and push this branch and make a PR to complete.</Text>
          <Newline />
          <Text>Press [ ctrl + c ] to exit</Text>
        </>
      )}
    </Box>
  );
};

export default Review;

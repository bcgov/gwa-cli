import * as React from 'react';
import { Box, Text, useApp, useInput } from 'ink';

import AppContext from '../../services/context';
import ConfigOrg from '../config-org';
import { IAppContext } from '../../types';
import Plugins from '../plugins';
import Review from '../review';
import Stepper from '../../components/stepper';

interface AppProps {
  args: IAppContext;
}

const App: React.FC<AppProps> = ({ args }) => {
  const [step, advance] = React.useState<number>(0);
  const [plugin, setPlugin] = React.useState<any | null>(null);
  const nextStep = () => advance((prev) => prev + 1);
  const { exit } = useApp();

  useInput((input, key) => {
    if (input === 'q' && key.ctrl) {
      exit();
    } else if (input === 'n' && key.ctrl) {
      advance((s) => Math.min(s + 1, 4));
    } else if (input === 'p' && key.ctrl) {
      advance((s) => Math.max(0, s - 1));
    }
  });

  return (
    <AppContext.Provider value={args}>
      <Box flexDirection="column">
        <Box
          borderStyle="round"
          borderColor="green"
          justifyContent="center"
          paddingX={1}
          marginBottom={1}
        >
          <Text>APS Gateway Configuration Tool</Text>
        </Box>
        <Box>
          <Stepper step={step}>
            <ConfigOrg onComplete={nextStep} step={step} />
            <Plugins
              onComplete={(item) => {
                setPlugin(item);
                nextStep();
              }}
              step={step}
            />
            <Review />
          </Stepper>
        </Box>
      </Box>
    </AppContext.Provider>
  );
};

export default App;

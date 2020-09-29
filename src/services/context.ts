import * as React from 'react';

import { AppContext } from '../types';

const AppContext = React.createContext<AppContext>({
  dir: '',
  file: '',
  version: '',
});

export default AppContext;

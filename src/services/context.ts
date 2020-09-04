import * as React from 'react';

import { IAppContext } from '../types';

const AppContext = React.createContext<IAppContext>({
  dir: '',
  file: null,
});

export default AppContext;

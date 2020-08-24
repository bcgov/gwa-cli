import * as React from 'react';

import { IAppContext } from '../types';

const AppContext = React.createContext<IAppContext>({
  dir: '',
});

export default AppContext;

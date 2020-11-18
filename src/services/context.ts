import * as React from 'react';

import { AppContext } from '../types';

export default React.createContext<AppContext>({
  dir: '',
  file: '',
  version: '',
});

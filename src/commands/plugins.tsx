//import { render } from 'ink';
import pluginsState from '../state/plugins';

export default function (...args) {
  console.log('plugins loaded', pluginsState.getState());
}

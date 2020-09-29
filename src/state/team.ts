import create from 'zustand/vanilla';
import createStore from 'zustand';
// import produce from 'immer';

type TeamState = {
  name: string;
  team: string;
  host: string;
};

const store = create<TeamState>(() => ({
  name: '',
  team: '',
  host: '',
}));

export const initTeamState = (data: TeamState) => store.setState(data);

export const useTeamState = createStore(store);

export default store;

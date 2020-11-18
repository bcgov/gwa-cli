import create from 'zustand/vanilla';

type SpecState = {};

const store = create<SpecState>(() => ({}));

export const loadSpec = (data: SpecState) => store.setState(data);

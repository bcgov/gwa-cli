export interface Prompt {
  label: string;
  key: string;
  constraint: any;
}

export interface SetupState {
  step: number;
  error: any | undefined;
  status: 'idle' | 'loading' | 'success' | 'failed';
  value: string;
  data: any;
  done: boolean;
  prompts: Prompt[];
  specError: string;
}
export type SetupAction =
  | { type: 'change'; payload: string }
  | { type: 'error'; payload: any | undefined }
  | { type: 'next'; payload: any }
  | { type: 'reset' }
  | { type: 'spec/loading' }
  | { type: 'spec/success'; payload: any }
  | { type: 'spec/failed'; payload: string }
  | { type: 'spec/written'; payload: string };

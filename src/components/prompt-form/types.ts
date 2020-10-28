export interface Prompt {
  label: string;
  key: string;
  secret?: boolean;
  constraint: any;
}

export interface SetupState {
  step: number;
  error: any | undefined;
  value: string;
  data: any;
  done: boolean;
  prompts: Prompt[];
}
export type SetupAction =
  | { type: 'change'; payload: string }
  | { type: 'error'; payload: any | undefined }
  | { type: 'next'; payload: any }
  | { type: 'done' }
  | { type: 'reset' };

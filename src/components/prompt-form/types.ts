export interface Prompt {
  label: string;
  key: string;
  secret?: boolean;
  constraint: any;
}

export interface PromptFormState {
  step: number;
  error: string[] | undefined;
  value: string;
  data: any;
  prompts: Prompt[];
}

export type PromptAction =
  | { type: 'change'; payload: string }
  | { type: 'next'; payload: string }
  | { type: 'reset' };

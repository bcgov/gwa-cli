import React from 'react';
import chalk from 'chalk';
import { render } from 'ink-testing-library';

import PromptForm from '../prompt-form';

const options = [
  {
    label: 'Prompt 1',
    key: 'prompt_1',
    constraint: {
      presence: { allowEmpty: false },
    },
  },
  {
    label: 'Prompt 1',
    key: 'prompt_1',
    constraint: {},
  },
];

describe('components/prompt-form', () => {
  const title = 'PromptForm Test';

  it('should render a title', () => {
    const { lastFrame } = render(<PromptForm title={title} options={[]} />);
    expect(lastFrame()).toEqual(expect.stringContaining(title));
  });

  it('should render the first option on render', () => {
    const { lastFrame } = render(
      <PromptForm title={title} options={options} />
    );
    expect(lastFrame()).toEqual(
      expect.stringContaining(`
${title}

${chalk.bold.green('?')} ${chalk.bold('Prompt 1')} ${chalk.inverse(' ')}`)
    );
  });
});

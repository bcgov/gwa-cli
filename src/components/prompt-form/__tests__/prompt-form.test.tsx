import React from 'react';
import chalk from 'chalk';
import delay from 'delay';
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
    label: 'Prompt 2',
    key: 'prompt_2',
    secret: true,
    constraint: {},
  },
];

describe('components/prompt-form', () => {
  const title = 'PromptForm Test';

  it('should render a title', () => {
    const onSubmit = jest.fn();
    const { lastFrame } = render(
      <PromptForm onSubmit={onSubmit} title={title} options={[]} />
    );
    expect(lastFrame()).toEqual(expect.stringContaining(title));
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('should show the first option on render', () => {
    const { lastFrame } = render(
      <PromptForm onSubmit={jest.fn()} title={title} options={options} />
    );
    expect(lastFrame()).toEqual(
      expect.stringContaining(`
${title}

${chalk.bold.green('?')} ${chalk.bold('Prompt 1')} ${chalk.inverse(' ')}`)
    );
  });

  it('should render the next option on successful submit', async () => {
    const { lastFrame, stdin } = render(
      <PromptForm onSubmit={jest.fn()} title={title} options={options} />
    );
    stdin.write('not empty');
    await delay(100);
    stdin.write('\r');

    expect(lastFrame()).toEqual(
      expect.stringContaining(`
${title}

${chalk.bold.green('✓')} ${chalk.bold('Prompt 1')} not empty
${chalk.bold.green('?')} ${chalk.bold('Prompt 2')} ${chalk.inverse(' ')}`)
    );
  });

  it('should render an error', async () => {
    const onSubmit = jest.fn();
    const { lastFrame, stdin } = render(
      <PromptForm onSubmit={onSubmit} title={title} options={[options[0]]} />
    );
    stdin.write('\r');

    expect(lastFrame()).toEqual(
      expect.stringContaining(`
${title}

${chalk.bold.green('?')} ${chalk.bold('Prompt 1')} ${chalk.inverse(
        ' '
      )}${chalk.red("<-- can't be blank")}`)
    );
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('should hide secrets', async () => {
    // Reverse the options so we can test the output of a secret
    const { lastFrame, stdin } = render(
      <PromptForm
        onSubmit={jest.fn()}
        title={title}
        options={[...options].reverse()}
      />
    );
    await delay(100);
    stdin.write('$ecret');
    await delay(100);
    stdin.write('\r');
    await delay(100);

    expect(lastFrame()).toEqual(
      expect.stringContaining(`
${title}

${chalk.bold.green('✓')} ${chalk.bold('Prompt 2')} **********
${chalk.bold.green('?')} ${chalk.bold('Prompt 1')} ${chalk.inverse(' ')}`)
    );
  });

  it('should call onSubmit', async () => {
    const onSubmit = jest.fn();
    const { stdin } = render(
      <PromptForm onSubmit={onSubmit} title={title} options={[options[0]]} />
    );
    await delay(100);
    stdin.write('value');
    await delay(100);
    stdin.write('\r');
    await delay(100);

    expect(onSubmit).toHaveBeenCalledWith({
      prompt_1: 'value',
    });
  });
});

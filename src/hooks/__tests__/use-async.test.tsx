import * as React from 'react';
import delay from 'delay';
import { Text } from 'ink';
import { render } from 'ink-testing-library';
import useAsync, { cache } from '../use-async';

const Fallback = <Text>Loading...</Text>;

describe('hooks/useAsync', () => {
  beforeEach(() => {
    cache.clear();
  });

  it('should cache requests', () => {
    const promise = jest.fn().mockResolvedValue(true);

    const Tester: React.FC<{}> = () => {
      const result = useAsync(promise, 'arg1');
      return <Text>{result}</Text>;
    };

    expect(cache.size).toEqual(0);
    render(<Tester />);
    expect(promise).toHaveBeenCalledWith('arg1');
    expect(cache.size).toEqual(1);
  });

  it('should return cached results', async () => {
    let inc = 1;
    const promise = jest.fn().mockResolvedValueOnce(inc);

    const Tester: React.FC<{}> = () => {
      const result = useAsync(promise, 'arg1');
      return <Text>{result}</Text>;
    };

    const { lastFrame, rerender } = render(
      <React.Suspense fallback={Fallback}>
        <Tester />
      </React.Suspense>
    );

    rerender(
      <React.Suspense fallback={Fallback}>
        <Tester />
      </React.Suspense>
    );
    await delay(100);
    expect(lastFrame()).toEqual('1');
  });
});

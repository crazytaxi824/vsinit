import React from 'react';

// factory mode, return a React hook function to share state value.
export function createSharedStates<T>(
  initState: T,
  resetState?: boolean // reset state value when on one is watching the state.
  // return a Readonly Type to prevent modification of state's value when use.
): () => [Readonly<T>, (state: T) => void] {
  let setStatesPool: ((state: T) => void)[] = [];
  let sharedState: T = initState;

  function setSharedState(newState: T) {
    // store new state value
    sharedState = newState;

    // setState() only notifies Components to refresh and get new state value.
    setStatesPool.forEach((setState) => {
      setState(newState);
    });
  }

  // return a Reack hook function
  return function useSharedState(): [Readonly<T>, (state: T) => void] {
    // gen a setState() for every Component which uses this hook.
    const [_, setState] = React.useState(initState);

    // useEffect() makes sure the following contents only excute once when Component init.
    React.useEffect(() => {
      // put setState() in the pool.
      setStatesPool.push(setState);

      return () => {
        // remove setState() from the pool when Component Unmount.
        setStatesPool = setStatesPool.filter((v) => v !== setState);

        // set state to init value, when no one is watch the state value.
        if (resetState && setStatesPool.length === 0) {
          sharedState = initState;
        }
      };
    }, []);

    return [sharedState, setSharedState];
  };
}

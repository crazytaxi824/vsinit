import React from 'react';

// factory mode, return a React hook function to share state value.
function createSharedReducer<S, A>(
  reducer: (prevState: S, action: A) => S,
  initState: S,
  resetState?: boolean // reset state value when on one is watching the state.
): () => [Readonly<S>, (action: A) => void] {
  let dispatchPool: ((action: A) => void)[] = [];
  let sharedState = initState;

  function sharedDispatch(action: A): void {
    // using reducer to calculate new state value.
    sharedState = reducer(sharedState, action);

    // dispatch() only notifies Components to refresh and get new state value.
    dispatchPool.forEach((dispatch) => {
      dispatch(action);
    });
  }

  // return a Reack hook function, which return a Readonly Type to prevent modification of state's value when use.
  return function useSharedReducer(): [Readonly<S>, (action: A) => void] {
    // gen a dispatch() for every Component which uses this hook.
    const [_, dispatch] = React.useReducer(reducer, initState);

    React.useEffect(() => {
      // put dispatch() in the pool.
      dispatchPool.push(dispatch);

      return () => {
        // remove dispatch() from the pool when Component Unmount.
        dispatchPool = dispatchPool.filter((v) => v !== dispatch);

        // set state to init value, when no one is watch the state value.
        if (resetState && dispatchPool.length === 0) {
          sharedState = initState;
        }
      };
    }, []);

    return [sharedState, sharedDispatch];
  };
}

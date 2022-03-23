// 使用拦截器来进行 retry 时, retryCount 对于每一个 AxiosInstance 都是全局的,
// 如果使用同一个 AxiosInstance 进行不同的 http 请求, retryCount 不会重置为 0
// 所以每次需要进行 retry 请求的时候, 一定要重新生成新的 AxiosInstance.

import axios, { AxiosInstance, AxiosError } from 'axios';

// 生成 AxiosInstance 实例
export function httpRetryInstance(retry = 3, delay = 1000): AxiosInstance {
  const httpInstance = axios.create();

  // 使用 retry 拦截器
  retryInterceptor(httpInstance, retry, delay);

  return httpInstance;
}

// 延迟请求用
export function sleep(duration: number): Promise<unknown> {
  return new Promise((resolve, reject) => {
    setTimeout(resolve, duration);
  });
}

// retry 拦截器
function retryInterceptor(
  httpInstance: AxiosInstance, // 传入 AxiosInstance
  retry = 3,
  delay = 1000
): void {
  let retryCount = 0;

  httpInstance.interceptors.response.use(
    undefined, // 返回成功的情况, 不需要任何处理
    async (error: AxiosError) => {
      const reqConfig = error.config; // 获取 http 请求参数

      if (retryCount >= retry) {
        return Promise.reject(error); // 返回错误
      }

      retryCount += 1;

      await sleep(delay); // 延迟 n 之后再执行

      return httpInstance(reqConfig); // ⭐️ 重新发起请求, 递归
    }
  );
}

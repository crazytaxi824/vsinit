import { AxiosInstance, AxiosError } from 'axios';

// 延迟请求用
export function sleep(duration: number): Promise<unknown> {
  return new Promise((resolve, reject) => {
    setTimeout(resolve, duration);
  });
}

// retry 拦截器
export function retryInterceptor(
  httpInstance: AxiosInstance, // 传入 AxiosInstance
  retry = 3,
  delay = 1000
): void {
  let retryCount = 0;

  httpInstance.interceptors.response.use(
    undefined,
    async (error: AxiosError) => {
      const reqConfig = error.config; // 获取 http 请求参数

      if (retryCount >= retry) {
        return Promise.reject(error); // 返回错误
      }

      retryCount += 1;

      await sleep(delay); // 延迟 n 之后再执行

      return httpInstance(reqConfig); // 重新发起请求
    }
  );
}

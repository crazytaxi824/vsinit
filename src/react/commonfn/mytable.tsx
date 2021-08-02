import React from 'react';
import axios, { AxiosError, AxiosInstance, AxiosResponse } from 'axios';

// 排序字段
type Order = 'asc' | 'desc' | null | undefined;

// 用于数据请求
interface ITableReqOption {
  page?: number;
  pageSize?: number;
  sortField?: string;
  sortOrder?: Order;
  // TODO 本示例未实现 filter，根据项目实际情况实现 filter
  filter?: string; 
}

// 数据结构，用于接受 http 数据
interface ITableResp<T> {
  error?: string;
  data: {
    rowCount: number; // 总共有多少条数据
    rows: T[]; // 返回的数据
  };
}

// 用于控制表格的操作
interface ITableControl {
  loading: boolean; // 是否显示 loading 界面
  editable: boolean; // data 是否可以编辑
  selection: boolean; // 是否显示 select 勾选框
}

// 定义类 table 用于驱动 DataGrid
export class Table<T> {
  // tableName 一旦实例化不可更改，用于 localStorage 浏览器储存数据
  private readonly tableName: string;

  // columns 一旦实例化不可修改
  private readonly columns?: string;

  // 默认请求数据
  private req: ITableReqOption = {
    page: 0, // page 从 0 开始计算，否则 <DataGrid> 中 rowCount 属性工作不正常。
    pageSize: 10, // 默认值
  };

  // 初始化返回数据
  private resp: ITableResp<T> = {
    data: {
      rowCount: 0,
      rows: [],
    },
  };

  // 初始化表格操作
  private action: ITableControl = {
    loading: false, // 是否显示 loading 界面
    editable: true, // data cell 是否可以编辑
    selection: true, // 是否显示 select 勾选框
  };

  // axios 请求实例
  private httpInstance: AxiosInstance;

  // ⚠️ delay 请求，如果限时时间内重新发起请求用 clearTimeout() 方法终止之前的 http 请求。
  private delayHttpReq?: NodeJS.Timeout;

  // 记录需要更新的字段，结构是 map[id]map[field]value，服务端用 map 解析 json
  // ⚠️ 记得 reset 的时候要清空该数据
  private cacheUpdateItems: Record<string, Record<string, unknown>> = {};

  // refresh 是一个 useState() 的闭包函数，用于刷新组件
  private refresh: () => void;

  // 构造器
  constructor(
    tableName: string,
    baseURL: string,
    refreshFn: () => void,
    columns?: (keyof T)[]
  ) {
    // DEBUG 检测 table 是否只会实例化一次
    console.log('constructor');

    this.tableName = tableName;

    // 创建 http 请求实例，// TODO 本示例中请求未使用身份凭证 auth，根据情况自行添加
    this.httpInstance = axios.create({ baseURL, timeout: 5000 });

    if (columns) {
      // ⚠️ columns 传入 url.query 时用, 字段之间用 ',' 分隔开。
      // columns=id,name,age...
      this.columns = columns.join(',');
    }

    // 刷新组件用
    this.refresh = refreshFn;

    // ⚠️ 获取 localstorage 中的 pageSize，如果已经存在则直接读取。如果不存在则使用默认值。
    const pageSize = localStorage.getItem(`${this.tableName}:pageSize`);
    if (pageSize) {
      this.req.pageSize = parseInt(pageSize, 10);
    }
  }

  // 将 pageSize 缓存到浏览器，使用 ${tableName}:pageSize 字段。
  storePageSizeLocally(pageSize: number): void {
    localStorage.setItem(`${this.tableName}:pageSize`, pageSize.toString());
  }

  // ⚠️ 请求 list 数据, 使用延迟发送, 减少请求次数。
  getList(
    config: { api?: string; reqOpt?: ITableReqOption; delay?: number } = {
      api: '',
      delay: 300,
    }
  ): void {
    // ⚠️赋值给 table req 放在 loading 界面前处理
    // 使表格组件可以快速刷新状态。
    this.req = { ...this.req, ...config.reqOpt };

    // 显示 loading 界面
    this.action.loading = true;
    this.refresh();

    // ⚠️ delay 时间内发起多次请求时，清除之前的 http 请求，只执行最后一个请求
    if (this.delayHttpReq) {
      clearTimeout(this.delayHttpReq);
    }

    // ⚠️ 延迟发送，如果请求变更了，则重新计时!!
    this.delayHttpReq = setTimeout(() => {
      this.httpInstance
        .get(config.api || '', {
          params: {
            ...this.req,
            ...config.reqOpt,
            ...{ columns: this.columns },
          },
        })
        .then((resp: AxiosResponse<ITableResp<T>>) => {
          // 确保 error 为 undefined，否则 <DataGrid> 报错。
          this.resp.error = resp.data.error || undefined;
          // 确保 rows 不为 null / undefined，否则 <DataGrid> 报错。
          this.resp.data.rows = resp.data.data.rows || [];
          // 确保 rowCount 不为 null / undefined，否则 <DataGrid> 报错。
          this.resp.data.rowCount = resp.data.data.rowCount || 0;

          // 关闭 loading 界面
          this.action.loading = false;
          this.refresh();

          // 返回 Promise
          return resp;
        })
        .catch((err: AxiosError) => {
          this.resp.error = handleAxiosError(err);

          // 关闭 loading 界面
          this.action.loading = false;
          this.refresh();
        });
    }, config.delay || 300); // 默认延迟时间，默认 300ms
  }

  // 通过 id 删除所选的 items
  deleteItems(ids: string[], api = ''): void {
    // 显示 loading 界面
    this.action.loading = true;
    this.refresh();

    this.httpInstance
      .delete(api, {
        params: { ids: ids.join(',') },
      })
      .then((resp: AxiosResponse<ITableResp<T>>) => {
        // ⚠️ 删除成功后，请求新数据，getList() 会关闭 loading 界面。
        this.getList({ delay: 0 });

        // 返回 Promise
        return resp;
      })
      .catch((err: AxiosError) => {
        this.resp.error = handleAxiosError(err);

        // 关闭 loading 界面
        this.action.loading = false;
        this.refresh();
      });
  }

  // update Item, 使用 cacheUpdateItems 缓存需要更新的数据
  addUpdateItems(item: { id: string; field: string; value: unknown }): void {
    const editItem = {};

    // 反射赋值给 object
    Reflect.set(editItem, item.field, item.value);
    // 记录到 cacheUpdateItems
    this.cacheUpdateItems[item.id] = {
      ...this.cacheUpdateItems[item.id],
      ...editItem,
    };

    // DEBUG
    console.log(JSON.stringify(this.cacheUpdateItems));
  }

  // 发送 update 更新数据到 server
  commitUpdate(api = ''): void {
    if (Object.keys(this.cacheUpdateItems).length === 0) {
      return;
    }

    // 显示 loading 界面
    this.action.loading = true;
    this.refresh();

    // 发送数据
    this.httpInstance
      .patch(api, this.cacheUpdateItems)
      .then((resp) => {
        // 关闭 loading 界面
        this.action.loading = false;
        this.refresh();

        return resp;
      })
      .catch((err: AxiosError) => {
        // 返回错误信息
        this.resp.error = handleAxiosError(err);

        // 关闭 loading 界面
        this.action.loading = false;
        this.refresh();
      });

    // 清空 updateItems
    this.cacheUpdateItems = {};
  }

  // 已经 edit 过，但是还没有发送到服务器的数据，使用该方法会重置。
  // ⚠️ 方法：强制重新渲染已经储存到 table.resp 中的数据。
  resetData(): void {
    // 清空需要 update 的数据
    this.cacheUpdateItems = {};

    // ⚠️ 要刷新表格数据必须生成新的 row 对象，否则表格会认为数据没有变更，而不会刷新表格数据。
    this.resp.data.rows = [...this.resp.data.rows];
    this.refresh();
  }

  // 是否显示选择框
  toggleSelection(): void {
    this.action.selection = !this.action.selection;
    this.refresh();
  }

  // 是否可以修改数据
  toggleEditable(): void {
    this.action.editable = !this.action.editable;
    this.refresh();
  }

  // ⚠️ 使用 Readonly 返回 table 数据拷贝，确保 table 属性不会被外部更改。
  get view(): Readonly<ITableReqOption & ITableResp<T>> {
    return { ...this.req, ...this.resp };
  }

  // ⚠️ 同上，返回控制表格样式的数据拷贝
  get control(): Readonly<ITableControl> {
    return { ...this.action };
  }
}

// ⚠️ 生成 table 对象
export function useTableView<T>(
  tableName: string, // 表格名字
  baseURL: string, // 请求地址
  columns?: (keyof T)[] // 表格要渲染的字段，需要和 <DataGrid> columns 对应。
  // 返回 Readonly Table
): Readonly<Table<T>> {
  // 只为刷新组件用
  const refreshFn = useRefreshComponent();

  // ⚠️ 使用 useRef 生成 initParams，防止 table 每次都重新实例化。
  const initParams = React.useRef({
    tableName,
    baseURL,
    columns,
  });

  // table 只会在组件初始化的时候生成一次。refreshFn 不会变。
  const table = React.useMemo(
    () =>
      new Table(
        initParams.current.tableName,
        initParams.current.baseURL,
        refreshFn,
        initParams.current.columns
      ),
    [refreshFn]
  );

  // 初始化 table 时，执行 getList() 获取数据
  // table 不变 useEffect 只会在初始化的时候执行一次。
  React.useEffect(() => {
    table.getList({ delay: 0 });
  }, [table]);

  return table;
}

// 使用 useState() 生成一个闭包函数，用来刷新组件。
function useRefreshComponent(): () => void {
  const [initMarker, refreshFn] = React.useState(true);

  // ⚠️ useRef 保证 initParams 不会变，使 useCallback()
  // 不会每次刷新时重新生成新的 function
  const initParams = React.useRef(initMarker);

  // useCallBack() 没有任何对比的对象，所以只会在组件初始化的时候执行一次。
  return React.useCallback(() => {
    initParams.current = !initParams.current;
    refreshFn(initParams.current);
  }, []);
}

// 处理 axios http 错误，返回错误信息。
export function handleAxiosError(err: AxiosError): string {
  if (err.response) {
    if (err.response?.status >= 500) {
      return `internal server error: ${err.message}`;
    }

    switch (err.response.status) {
      case 400:
        return 'request params error, please contact Admin';

      case 401:
        // history.push('/login')
        return 'need to login first';

      case 403:
        return 'no authorization to this action';

      default:
        return `error status code: [${err.response.status}], please contact Admin`;
    }
  } else {
    console.log(err);
    return `error: ${err.message}, please contact Admin`;
  }
}

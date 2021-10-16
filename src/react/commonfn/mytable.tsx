import React from 'react';
import axios, { AxiosError, AxiosInstance, AxiosResponse } from 'axios';

// 排序字段, 定义和 GridSortDirection 一样。
type Order = 'asc' | 'desc' | null | undefined;

// 用于数据请求
interface ITableReqOption {
  page: number;
  pageSize: number;
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
  // 默认请求数据
  req: ITableReqOption = {
    page: 0, // page 从 0 开始计算，否则 <DataGrid> 中 rowCount 属性工作不正常。
    pageSize: 10, // 默认值
  };

  // 初始化返回数据
  resp: ITableResp<T> = {
    data: {
      rowCount: 0,
      rows: [],
    },
  };

  // 初始化表格操作
  control: ITableControl = {
    loading: false, // 是否显示 loading 界面
    editable: false, // data cell 是否可以编辑
    selection: true, // 是否显示 select 勾选框
  };

  // ⭐️ 用于合并请求用，延迟发送请求。
  timeOut?: NodeJS.Timeout;

  // axios 请求实例, 一旦实例化不可修改
  private readonly httpInstance: AxiosInstance;

  // columns 一旦实例化不可修改
  private readonly columns?: string;

  // 记录需要更新的字段，⭐️ 记得 reset 的时候要清空该数据
  // 结构是 map[id:string]map[field:string]value:interface{}
  private cacheUpdateItems: Record<string, Record<string, unknown>> = {};

  // 构造器
  constructor(
    private readonly tableName: string, // tableName 一旦实例化不可更改，用于 localStorage 浏览器储存数据
    baseURL: string,
    public refresh: () => void, // refresh 是一个 useState() 的闭包函数，用于刷新组件
    columns?: (keyof T)[] // columns 只能是返回数据的字段名。
  ) {
    // DEBUG 检测 table 是否只会实例化一次
    console.log('table class constructor');

    // 创建 http 请求实例，⭐️ 本示例中请求未使用 auth，根据情况自行添加
    this.httpInstance = axios.create({ baseURL, timeout: 5000 });

    // columns 传入 url.query 时用, 字段之间用 ',' 分隔开。
    this.columns = columns ? columns.join(',') : undefined;

    // 获取 localstorage 中的 pageSize，如果已经存在则直接读取。如果不存在则使用默认值。
    const pageSize = localStorage.getItem(`${this.tableName}:pageSize`);
    this.req.pageSize = pageSize ? parseInt(pageSize, 10) : this.req.pageSize;

    console.log(this.req.pageSize);
  }

  // 将 pageSize 缓存到浏览器，使用 ${tableName}:pageSize 字段。
  storePageSizeLocally(pageSize: number): void {
    localStorage.setItem(`${this.tableName}:pageSize`, pageSize.toString());
    // localStorage.setItem(this.tableName + ':pageSize', pageSize.toString());
  }

  // ⭐️ 请求 list 数据, 使用延迟发送, 减少请求次数。
  getList(api = ''): void {
    console.log(api);

    // 显示 loading 界面
    this.showLoading(true);

    // 发送请求
    this.httpInstance
      .get(api, {
        params: {
          columns: this.columns || undefined, // 如果 columns 不存在这里不会有 column
          ...this.req,
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
        this.showLoading(false);

        // 返回 Promise
        return resp;
      })
      .catch((err: AxiosError) => {
        this.resp.error = handleAxiosError(err);

        // 关闭 loading 界面
        this.showLoading(false);
      });
  }

  // 通过 id 删除所选的 items
  deleteItems(ids: string[], api = ''): void {
    // 如果没有 delete item 则直接终止
    if (ids.length === 0) {
      return;
    }

    // 显示 loading 界面
    this.showLoading(true);

    this.httpInstance
      .delete(api, {
        params: { ids: ids.join(',') },
      })
      .then((resp: AxiosResponse<ITableResp<T>>) => {
        // ⭐️ 删除成功后，请求新数据，getList() 会关闭 loading 界面。
        this.getList();

        // 返回 Promise
        return resp;
      })
      .catch((err: AxiosError) => {
        this.resp.error = handleAxiosError(err);

        // 关闭 loading 界面
        this.showLoading(false);
      });
  }

  // update Item, 使用 cacheUpdateItems 缓存需要更新的数据
  addUpdateItems(item: { id: string; field: string; value: unknown }): void {
    // 记录到 cacheUpdateItems
    this.cacheUpdateItems[item.id] = {
      ...this.cacheUpdateItems[item.id],

      // ⭐️ 使用 item.field 的值作为 key 的写法
      ...{ [item.field]: item.value },
    };

    // DEBUG
    console.log('cacheUpdateItems:', JSON.stringify(this.cacheUpdateItems));
  }

  // 发送 update 更新数据到 server
  commitUpdate(api = ''): void {
    // DEBUG
    console.log('updateItems:', JSON.stringify(this.cacheUpdateItems));

    // 如果没有 update item 直接终止。
    if (Object.keys(this.cacheUpdateItems).length === 0) {
      return;
    }

    // 显示 loading 界面
    this.showLoading(true);

    // 发送数据
    this.httpInstance
      .patch(api, this.cacheUpdateItems)
      .then((resp) => {
        // 关闭 loading 界面
        this.showLoading(false);

        return resp;
      })
      .catch((err: AxiosError) => {
        // 返回错误信息
        this.resp.error = handleAxiosError(err);

        // 关闭 loading 界面
        this.showLoading(false);
      });

    // 清空 updateItems
    this.cacheUpdateItems = {};
  }

  // 已经 edit 过，但是还没有发送到服务器的数据，使用该方法会重置。
  // 方法：强制重新渲染已经储存到 table.resp 中的数据。
  resetData(): void {
    // 清空需要 update 的数据
    this.cacheUpdateItems = {};

    // ⭐️ 要强制刷新表格数据必须生成新的 row 对象，
    // 否则表格会认为数据没有变更，而不会刷新表格数据。
    this.resp.data.rows = [...this.resp.data.rows];
    this.refresh();
  }

  // 是否显示 loading 界面
  showLoading(b: boolean): void {
    this.control.loading = b;
    this.refresh();
  }

  // 是否显示选择框
  allowSelect(b: boolean): void {
    this.control.selection = b;
    this.refresh();
  }

  // 是否可以修改数据
  allowEdit(b: boolean): void {
    this.control.editable = b;
    this.refresh();
  }
}

// ⭐️ 生成 table 对象
export function useTableView<T>(
  tableName: string, // 表格名字
  baseURL: string, // 请求地址
  columns?: (keyof T)[] // 表格要渲染的字段，需要和 <DataGrid> columns 对应。
  // 返回 Readonly Table
): Table<T> {
  // 只为刷新组件用
  const refreshFn = useRefreshComponent();

  // ⭐️ 使用 useRef 生成 initParams，防止 table 每次都重新实例化。
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
    table.getList();
  }, [table]);

  return table;
}

// 使用 useState() 生成一个闭包函数，用来刷新组件。
export function useRefreshComponent(): () => void {
  const [_, refreshFn] = React.useState({});

  // useCallBack() 没有任何对比的对象，所以只会在组件初始化的时候执行一次。
  return React.useCallback(() => {
    refreshFn({});
  }, []);
}

// 处理 axios http 错误，返回错误信息。
export function handleAxiosError(err: AxiosError): string {
  if (err.response) {
    if (err.response.status >= 500) {
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

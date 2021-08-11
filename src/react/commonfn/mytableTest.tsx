import React from 'react';
import {
  DataGrid,
  GridColDef,
  GridPageChangeParams,
  GridValueGetterParams,
  GridToolbar,
  GridColumnMenuContainer,
  GridToolbarContainer,
  GridToolbarColumnsButton,
  GridToolbarDensitySelector,
  GridToolbarFilterButton,
  GridToolbarExport,
  useGridSlotComponentProps,
  GridRowId,
  GridOverlay,
  GridFilterModelState,
  GridFilterModel,
  GridEditCellPropsParams,
  GridEditRowModelParams,
  GridRowData,
  GridSortModel,
} from '@material-ui/data-grid';
import { Pagination } from '@material-ui/lab';
import {
  Box,
  Button,
  Select,
  MenuItem,
  Grid,
  LinearProgress,
} from '@material-ui/core';
import { useTableView, Table } from './mytable';

// 顶部控制组件，默认提供 filter，columns... 等功能
function CustomToolBar(props: { table: Table<ITableData> }) {
  // ⚠️⚠️⚠️ useGridSlotComponentProps() 可以获取 DataGrid 内置的方法和属性。
  const { state, apiRef } = useGridSlotComponentProps();

  return (
    <Box display="flex" alignItems="center" width="100%" height="60px">
      {state.selection.length !== 0 ? (
        <Button
          variant="outlined"
          onClick={() => {
            props.table.deleteItems(
              // 把 ids 都转成 string，方便传输。
              state.selection.map((id) => id.toString())
            );
          }}
        >
          Delete
        </Button>
      ) : (
        <Button variant="outlined">edit</Button>
      )}
      <Box flexGrow={1} />
      <GridToolbarColumnsButton /> {/* 控制显示那些字段 */}
      <GridToolbarFilterButton /> {/* 控制 filter */}
      <GridToolbarDensitySelector /> {/* 控制表格大小 */}
      {/* <GridToolbarExport /> */} {/* 导出表格到文件 */}
    </Box>
  );
}

// 官方示例提供的 loading 组件，不用改。
const CustomLoadingOverlay = React.memo(() => (
  <GridOverlay>
    <div style={{ position: 'absolute', top: 0, width: '100%' }}>
      <LinearProgress />
    </div>
  </GridOverlay>
));

// 显示错误信息的组件
function CustomErrorOverlay() {
  const { state } = useGridSlotComponentProps();

  return (
    <div
      style={{
        height: '300px',
        width: '100%',
        textAlign: 'center',
        lineHeight: '300px',
      }}
    >
      {state.error}
    </div>
  );
}

// footer 显示 page 和 pageSize
function CustomFooter(props: { table: Table<ITableData> }) {
  // ⚠️⚠️⚠️ useGridSlotComponentProps() 可以获取 DataGrid 内置的方法和属性。
  const { state, apiRef } = useGridSlotComponentProps();

  return (
    <div style={{ height: '60px' }}>
      <Grid container style={{ height: '100%' }}>
        {/* 左侧部分空出来做其他的 */}
        <Grid
          item
          xs={3}
          container
          justifyContent="flex-start"
          alignItems="center"
        >
          <Button
            variant="outlined"
            onClick={() => props.table.toggleSelection()}
          >
            Selection
          </Button>
        </Grid>

        {/* 中间部分显示 page */}
        <Grid item xs={6} container justifyContent="center" alignItems="center">
          <Pagination
            color="primary"
            // ⚠️ 页码总数。这里由 rowCount 自动计算。
            count={state.pagination.pageCount}
            // ⚠️ 使用 rowCount page 从 0 开始计算，所以显示的时候需要 +1，返回的时候需要 -1。
            page={state.pagination.page + 1}
            onChange={(event, page) => apiRef.current.setPage(page - 1)}
          />
        </Grid>
        {/* 右侧部分显示 page size */}
        <Grid
          item
          xs={3}
          container
          alignItems="center"
          justifyContent="flex-end"
        >
          Page Size:
          <Select
            // ⚠️ 这里要设置控制，否则会报 warning，DataGrid 默认 pageSize 100 上限。
            value={
              state.pagination.pageSize > 50 ? 50 : state.pagination.pageSize
            }
            style={{
              minWidth: 40,
              marginLeft: '0.5rem',
              marginRight: '1rem',
            }}
            onChange={(ev) => {
              apiRef.current.setPageSize(ev.target.value as number);
            }}
          >
            <MenuItem value={10}>10</MenuItem>
            <MenuItem value={20}>20</MenuItem>
            <MenuItem value={50}>50</MenuItem>
          </Select>
        </Grid>
      </Grid>
    </div>
  );
}

// ⚠️ 返回表头字段和操作，editable 表示是否可以编辑 Data Cell
function tableHeaders(editable: boolean): GridColDef[] {
  return [
    { field: 'id', headerName: 'ID', width: 90 },
    {
      field: 'firstName',
      headerName: 'First name',
      // headerAlign: 'center', // 需要关闭 sortable 和 disableColumnMenu
      // align: 'right', // 数据显示的位置
      width: 150,
      editable,
      sortable: false,
    },
    {
      // row 数据中的字段名
      field: 'lastName',

      // Table header 中显示的名字，如果不填则会直接显示 field 名字。
      headerName: 'Last name',

      // 数据类型，默认 string
      type: 'string',

      // column 宽度，可以设置 flex 来根据屏幕宽度动态改变大小。
      // 如果设置 flex 可以设置 minWidth。
      width: 150,

      // 允许直接编辑 row 中的数据，双击编辑
      editable,

      // 不允许 sort，默认为 true。
      // 如果为 true 的话，header 中会出现 sort 的图标
      sortable: false,

      // 不允许 filter
      filterable: false,
    },
    {
      field: 'age',
      headerName: 'Age',
      type: 'number',
      width: 110,
      editable,
    },
    {
      field: 'fullName',
      headerName: 'Full name',
      description: 'This column has a value getter and is not sortable.',

      sortable: false,
      width: 160,

      // ⚠️ valueGetter 函数根据已有数据生成新的数据。
      valueGetter: (params: GridValueGetterParams) => {
        // 获取 firstName
        const fn = params.getValue(params.id, 'firstName');
        // 获取 lastName
        const ln = params.getValue(params.id, 'lastName');

        // 返回 firstName + lastName
        return `${fn ? fn.toString() : ''} ${ln ? ln.toString() : ''}`;
      },
    },
  ];
}

// ⚠️ 自定义的 table 需要获取的数据结构，Table<T> 用
interface ITableData {
  id: number;
  firstName: string;
  lastName: string;
  age: number;
}

// ⚠️ 所有数据都必须有唯一 id，如果没有 id 字段，DataGrid 会报错!!
export function MyDataTable(): JSX.Element | null {
  console.log('mytable');

  // const table = useTableView<ITableData>('/api', ['id', 'age', 'firstName']);
  const table = useTableView<ITableData>('testTable', '/api');

  // 每次重渲染的时候 sort 都会重新发送请求。
  // 记录 sortModel 状态，判断是否需要停止 sort 发送请求。
  const sortRef = React.useRef<GridSortModel>([]);

  return (
    <>
      {/*  这里最重要的是宽度必须设置。DataGrid 无法撑开 div. */}
      <div style={{ width: '100%' }}>
        <DataGrid
          // 数据设置 --------------------------------------------------------
          // column 数据，这里是设置表格头是否可以编辑等...
          columns={tableHeaders(table.control.editable)}
          // ⚠️ row 数据，这里是 body 数据, 确保 row 不能为 null / undefined
          rows={table.view.data.rows}
          // DateGrid 属性设置 -----------------------------------------------
          autoHeight // 自动高度
          checkboxSelection={table.control.selection} // 允许勾选数据
          disableSelectionOnClick // 禁止点击 row 的时候选取 checkbox，只能手动点击 checkbox 才行。
          disableColumnMenu // 禁止在 column 中显示 menu icon
          // component 设置 -------------------------------------------------
          components={{
            // Toolbar: GridToolbar, // 官方提供的 toolbar 组件
            Toolbar: CustomToolBar,
            Footer: CustomFooter,
            LoadingOverlay: CustomLoadingOverlay,
            ErrorOverlay: CustomErrorOverlay,
          }}
          // 给不同的 component 传递 props
          componentsProps={{ footer: { table }, toolbar: { table } }}
          // 是否显示 loading
          loading={table.control.loading}
          // 是否显示 error 界面，error 不为 null 则会显示。
          error={table.view.error || null}
          // 以下全是 pagination 设置 ------------------------------------------
          paginationMode="server"
          // 使用内置 rowCount 方法自动计算一共有多少页数据
          // ⚠️⚠️⚠️ page 必须从 0 开始做为第一页，否则最后2页数据长度不对。
          rowCount={table.view.data.rowCount}
          // 显示第几页数据
          page={table.view.page}
          // 每页显示多少条数据
          pageSize={table.view.pageSize}
          // 切换页面时加载数据
          onPageChange={(params: GridPageChangeParams) => {
            table.getList({ reqOpt: { page: params.page } });
          }}
          // pageSize 改变时运行
          onPageSizeChange={(params: GridPageChangeParams) => {
            // 使用 localStorage 将 pageSize 储存到浏览器。
            table.storePageSizeLocally(params.pageSize);
            // ⚠️⚠️⚠️ 改变 page size 的时候要注意同时接受 page 和 pagesize，
            // 因为超出 page*pageSize 的数据为空。框架会自动返回最大的 page 数。
            table.getList({
              reqOpt: { page: params.page, pageSize: params.pageSize },
            });
          }}
          // update 数据 -------------------------------------------------
          // item 数据变更时缓存起来，后面统一更新。
          onEditCellChange={(params: GridEditCellPropsParams) => {
            const editItem = {
              id: params.id.toString(),
              field: params.field,
              value: params.props.value,
            };
            table.addUpdateItems(editItem);
          }}
          // DEBUG 模拟发送请求，不推荐使用 onEditCellChangeCommitted 记录 update 数据。
          // ⚠️ 使用 onEditCellChangeCommitted() 的时候只有点回车才会返回有效 props.value 数据。
          onEditCellChangeCommitted={(params: GridEditCellPropsParams) => {
            if (params.props) {
              table.commitUpdate(); // 发送请求
            } else {
              table.resetData(); // 重置数据，不发送 update 请求。
            }
          }}
          // sorting -----------------------------------------------------
          // ⚠️⚠️⚠️ sorting 巨坑，只要sorting 状态不是 undefined，
          // 每次重渲染页面都会触发 onSortModelChange()，很容易造成无限循环。
          sortingMode="server"
          onSortModelChange={(newModel: GridSortModel) => {
            // ⚠️⚠️⚠️ 必须在 sort 数据相同的情况下停下来，否则会造成无限循环
            if (JSON.stringify(sortRef.current) === JSON.stringify(newModel)) {
              return;
            }

            // 通过 useRef() 记录 sortModel 状态，如果状态相同则停止向下执行。
            sortRef.current = newModel;

            // 请求数据
            table.getList({
              reqOpt: {
                sortField: newModel[0] ? newModel[0].field : undefined,
                sortOrder: newModel[0] ? newModel[0].sort : undefined,
              },
            });
          }}
          // filtering ----------------------------------------------------
          // filter 是实时发送的，所以采用和 edit 相同的办法。先缓存，然后统一发送。
          // DataGrid 不提供 multi-filter，需要自己实现 component。
          filterMode="server"
          // DEBUG
          onFilterModelChange={(filter: GridFilterModelState) => {
            console.log(JSON.stringify(filter));
          }}
        />
      </div>
    </>
  );
}

import React from "react";
import { DataGrid, GridRowId, GridRowData, GridColDef, GridValueGetterParams, GridToolbar, GridColumnMenuContainer, GridToolbarContainer, GridToolbarExport, GridToolbarColumnsButton, GridToolbarDensitySelector, GridToolbarFilterButton, useGridSlotComponentProps, GridOverlay, GridFilterModel, GridSortModel, GridFooter, GridRowCount, GridSelectedRowCount, GridCellParams, GridCellEditCommitParams } from "@mui/x-data-grid";
import { Pagination, Box, Button, Select, MenuItem, Grid, LinearProgress } from "@mui/material";

import { Table, useTableView, useRefreshComponent } from "./table";

// 顶部控制组件，默认提供 filter，columns... 等功能
function CustomToolBar(props: { table: Table<ITableData> }) {
	// ⭐️ useGridSlotComponentProps() 可以获取 DataGrid 内置的方法和属性。
	const { state, apiRef } = useGridSlotComponentProps();

	return (
		<Grid container sx={{ height: 60 }}>
			{/* 左侧部分 delete selected item */}
			<Grid item xs={2} container justifyContent="flex-start" alignItems="center">
				{state.selection.length !== 0 && (
					<Button
						variant="outlined"
						onClick={() => {
							props.table.deleteItems(
								// 把 ids 都转成 string，方便传输。
								state.selection.map((id) => id.toString())
							);

							// 清空已选的 rows, 这里依靠 loading 的 refresh 刷新了组件。
							state.selection = [];
						}}
					>
						Delete
					</Button>
				)}
			</Grid>

			{/* 中间部分 edit item */}
			<Grid item xs={4} container justifyContent="center" alignItems="center">
				{props.table.control.editable ? (
					<div>
						<Button
							variant="outlined"
							onClick={() => {
								props.table.commitUpdate();
								props.table.allowEdit(false);
							}}
						>
							commit update
						</Button>
						<Button
							variant="outlined"
							onClick={() => {
								props.table.resetData();
								props.table.allowEdit(false);
							}}
						>
							cancel edit
						</Button>
					</div>
				) : (
					<Button
						variant="outlined"
						onClick={() => {
							props.table.allowEdit(true);
						}}
					>
						start edit
					</Button>
				)}
			</Grid>

			{/* 右侧部分 - 官方提供的组件 */}
			<Grid item xs={6} container justifyContent="flex-end" alignItems="center">
				<GridToolbarColumnsButton /> {/* 控制显示那些字段 */}
				<GridToolbarFilterButton /> {/* 控制 filter */}
				<GridToolbarDensitySelector /> {/* 控制表格大小 */}
				{/* <GridToolbarExport /> */} {/* 导出表格到文件 */}
			</Grid>
		</Grid>
	);
}

// 官方示例提供的 loading 组件，不用改。
const CustomLoadingOverlay = React.memo(() => (
	<GridOverlay>
		<div style={{ position: "absolute", top: 0, width: "100%" }}>
			<LinearProgress />
		</div>
	</GridOverlay>
));

// 显示错误信息的组件
const CustomErrorOverlay = React.memo(() => {
	const { state } = useGridSlotComponentProps();

	return <div style={{ height: "300px", textAlign: "center", lineHeight: "300px" }}>{state.error}</div>;
});

// footer 显示 page 和 pageSize
function CustomFooter(props: { table: Table<ITableData> }) {
	const { state, apiRef } = useGridSlotComponentProps();

	// 生成一个组件刷新函数。
	const refresh = useRefreshComponent();

	return (
		<Grid container sx={{ height: 60 }}>
			{/* 左侧部分显示选择多少行 */}
			<Grid item xs={2} container justifyContent="flex-start" alignItems="center">
				{state.selection.length > 0 && <div style={{ marginLeft: "20px" }}>{state.selection.length} row(s) selected</div>}
			</Grid>

			{/* 中间部分显示 page */}
			<Grid item xs={8} container justifyContent="center" alignItems="center">
				<Pagination
					color="primary"
					// ⭐️ 页码总数。这里由 rowCount 自动计算。
					count={Math.ceil(props.table.resp.data.rowCount / props.table.req.pageSize)}
					// ⭐️ 使用 page 从 0 开始计算，所以显示的时候需要 +1，请求的时候需要 -1。
					page={props.table.req.page + 1}
					onChange={(event, newPage) => {
						// 先刷新 pagination 显示
						props.table.req.page = newPage - 1;
						refresh();

						// ⭐️ 延迟发送请求，便于快速翻页。
						if (props.table.timeOut) {
							clearTimeout(props.table.timeOut);
						}

						props.table.timeOut = setTimeout(() => {
							props.table.getList();
						}, 300);
					}}
				/>
			</Grid>
			{/* 右侧部分显示 page size */}
			<Grid item xs={2} container justifyContent="flex-end" alignItems="center">
				Page Size:
				<Select
					value={props.table.req.pageSize}
					variant="standard"
					style={{
						minWidth: 40,
						marginLeft: "0.5rem",
						marginRight: "1rem",
					}}
					onChange={(ev) => {
						const newPageSize = ev.target.value as number;

						// 使用 localstorage 储存 pageSize
						props.table.storePageSizeLocally(newPageSize);

						// 发送请求
						props.table.req.page = 0;
						props.table.req.pageSize = newPageSize;
						props.table.getList();
					}}
				>
					<MenuItem value={10}>10</MenuItem>
					<MenuItem value={20}>20</MenuItem>
					<MenuItem value={50}>50</MenuItem>
				</Select>
			</Grid>
		</Grid>
	);
}

// 返回表头字段和操作，editable 表示是否可以编辑 Data Cell
function tableHeaders(editable: boolean): GridColDef[] {
	return [
		{ field: "id", headerName: "ID", width: 100, editable: false },
		{
			field: "firstName",
			headerName: "First name",
			// headerAlign: 'center', // 需要关闭 sortable 和 disableColumnMenu
			// align: 'right', // 数据显示的位置
			width: 150,
			editable: editable,
			sortable: true,
		},
		{
			// row 数据中的字段名
			field: "lastName",

			// Table header 中显示的名字，如果不填则会直接显示 field 名字。
			headerName: "Last name",

			// 数据类型，默认 string
			type: "string",

			// column 宽度，可以设置 flex 来根据屏幕宽度动态改变大小。
			// 如果设置 flex 可以设置 minWidth。
			width: 150,

			// 允许直接编辑 row 中的数据，双击编辑
			editable: editable,

			// 不允许 sort，默认为 true。
			// 如果为 true 的话，header 中会出现 sort 的图标
			sortable: true,

			// 不允许 filter
			filterable: true,
		},
		{
			field: "age",
			headerName: "Age",
			type: "number",
			width: 100,
			editable: true,
		},
		{
			field: "fullName",
			headerName: "Full name",
			description: "This column has a value getter and is not sortable.",

			sortable: false,
			width: 300,

			// ⭐️ valueGetter 函数根据已有数据生成新的数据。
			valueGetter: (params: GridValueGetterParams) => {
				// 获取 firstName
				const fn = params.getValue(params.id, "firstName");
				// 获取 lastName
				const ln = params.getValue(params.id, "lastName");
				// 返回 firstName + lastName
				return `${fn ? fn.toString() : ""} ${ln ? ln.toString() : ""}`;
			},
		},
	];
}

// 自定义的 table 需要获取的数据结构，Table<T> 用
interface ITableData {
	id: number;
	firstName: string;
	lastName: string;
	age: number;
}

// 这里没有使用 pagination 的相关属性，所有的 pagination 数据通过 `table` 直接传递给 `CustomFooter`
export function MyDataTable2(): JSX.Element {
	console.log("mytable");

	// ⭐️ 生成 table 对象，用于驱动表格。
	// const table = useTableView<ITableData>('testTable', '/api', ['id', 'age', 'firstName']);
	const table = useTableView<ITableData>("testTable", "/api");

	return (
		<>
			{/*  这里最重要的是宽度必须设置。DataGrid 无法撑开 div. */}
			<div style={{ width: "100%" }}>
				<DataGrid
					// 数据设置 -----------------------------------------------------
					// column 数据，设置表格头是否可以编辑。
					columns={tableHeaders(table.control.editable)}
					// row 数据，确保 row 不能为 null / undefined
					rows={table.resp.data.rows}
					// DateGrid 属性设置 --------------------------------------------
					autoHeight // 自动高度
					density="compact" // 行间距最窄的模式
					checkboxSelection={table.control.selection} // 允许勾选数据
					disableSelectionOnClick // 禁止点击 row 的时候选取 checkbox，只能手动点击 checkbox 才行。
					// disableColumnMenu // 禁止在 column 中显示 menu icon
					loading={table.control.loading} // 是否显示 loading
					error={table.resp.error || undefined} // 是否显示 error 界面，error 不为 undefined | null 则会显示。
					// component 设置 -----------------------------------------------
					components={{
						// Toolbar: GridToolbar, // 官方提供的 toolbar 组件
						Toolbar: CustomToolBar,
						Footer: CustomFooter,
						LoadingOverlay: CustomLoadingOverlay,
						ErrorOverlay: CustomErrorOverlay,
					}}
					// 给不同的 component 传递 props
					componentsProps={{ footer: { table }, toolbar: { table } }}
					// pagination -------------------------------------------------
					// 这里没有使用 pagination 的相关属性，所有的 pagination 数据通过 `table` 直接传递给 `CustomFooter`

					// update 数据 -------------------------------------------------
					// item 数据变更时缓存起来，后面统一更新。
					// ⭐️ onCellEditCommit 会返回两种类型：
					// 如果 ev.type 是 'keydown' - 'enter', 'tab' 则返回 GridCellEditCommitParams 类型。
					// 如果 ev.type 是 'click' - 点击其他的地方则返回 GridCellParams 类型。
					onCellEditCommit={(params, ev) => {
						const p = params as GridCellParams;

						// DEBUG
						// console.log(ev.type, p.value, p.formattedValue);

						if (ev.type === "keydown") {
							const editItem = {
								id: params.id.toString(),
								field: params.field,
								value: params.value,
							};
							table.addUpdateItems(editItem);
							return;
						} else if (p.formattedValue !== p.value) {
							// 重置数据
							table.resetData();
						}
					}}
					// sorting -----------------------------------------------------
					// x-data-grid 5.0.0-beta3 之前这里有错误。
					sortingMode="server"
					onSortModelChange={(newModel: GridSortModel) => {
						if (table.timeOut) {
							clearTimeout(table.timeOut);
						}
						table.timeOut = setTimeout(() => {
							// 请求数据
							table.req.sortField = newModel[0] && newModel[0].field;
							table.req.sortOrder = newModel[0] && newModel[0].sort;
							table.getList();
						}, 300);
					}}
					// filtering ---------------------------------------------------
					// filter 是实时发送的，所以采用和 edit 相同的办法。先缓存，然后统一发送。
					// DataGrid 不提供 multi-filter，需要自己实现 component。
					filterMode="server"
					// DEBUG
					onFilterModelChange={(filter: GridFilterModel) => {
						console.log(JSON.stringify(filter.items));
					}}

					// selection ---------------------------------------------------
					// selection 一般不需要自己处理，表格会自动处理。
					// selectionModel={15}
					// onSelectionModelChange={(model) => {
					//   console.log(model); // 返回一个 number[],eg: [11,12,13]
					// }}
				/>
			</div>
		</>
	);
}

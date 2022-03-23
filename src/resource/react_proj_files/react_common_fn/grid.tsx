import React from 'react';
import { Grid, Paper, GridSize } from '@mui/material';

// 设置 Paper 的 style
const PaperItem = (props: { content: string }): JSX.Element => {
  return (
    <Paper
      sx={{
        bgcolor: 'primary.main',
        color: 'primary.contrastText',
        height: '60px',
        lineHeight: '60px',
        textAlign: 'center',
      }}
    >
      {props.content ? props.content : 'xs'}
    </Paper>
  );
};

// 基本用法
export function GridTest(): JSX.Element {
  const xsValues: GridSize[] = [8, 4, 6, 6, 4, 4, 4, 2, 2, 2, 2, 2, 2, 6, 4, 2];

  const gridItems = xsValues.map((v, i) => {
    return (
      // 这里使用 xs={} 来设置 item 占多少个 grid 格子
      <Grid item xs={v}>
        <PaperItem content={'xs=' + v.toString()} />
      </Grid>
    );
  });

  return (
    // 外面套一个 div 来约束 Grid
    <div style={{ margin: '20px', borderStyle: 'dashed', borderWidth: '1px' }}>
      <Grid container spacing={1}>
        {gridItems}
      </Grid>
    </div>
  );
}

export function GridAutoLayout(): JSX.Element {
  const arry = [1, 2, 3, 4, 5];

  const gridItems = arry.map(() => {
    return (
      // 这里使用 xs=true 来自动布局, 也可以写作 <Grid item xs={true}>
      <Grid item xs>
        <PaperItem content="xs" />
      </Grid>
    );
  });

  return (
    // 外面套一个 div 来约束 Grid
    <div style={{ margin: '20px', borderStyle: 'dashed', borderWidth: '1px' }}>
      <Grid container spacing={1}>
        {gridItems}
      </Grid>
    </div>
  );
}

// 多行自动布局
export function GridAutoLayoutMulti(): JSX.Element {
  const arry = [1, 2, 3, 4, 5];

  const gridItems = arry.map(() => {
    return (
      <Grid item xs>
        <PaperItem content="xs" />
      </Grid>
    );
  });

  const arry2 = [1, 2, 3, 4, 5];

  const gridItems2 = arry2.map(() => {
    return (
      // 这里使用 xs=true 来自动布局
      <Grid item xs>
        <PaperItem content="xs" />
      </Grid>
    );
  });

  return (
    <div style={{ margin: '20px', borderStyle: 'dashed', borderWidth: '1px' }}>
      <Grid container spacing={1}>
        {/* 多行的自动布局需要使用多个 container 嵌套, 否则会排列在一行中 */}
        {/* xs 是 Grid item 的属性, xs={12} 表示占满一整行 */}
        <Grid item container xs={12} spacing={2}>
          {gridItems}
        </Grid>
        <Grid item container xs={12} spacing={4}>
          {gridItems2}
        </Grid>
      </Grid>
    </div>
  );
}

export function GridAutoLayoutMulti2(): JSX.Element {
  const arry = [1, 2, 3, 4, 5];

  const gridItems = arry.map(() => {
    return (
      // 这里使用 xs=true 来自动布局
      <Grid item xs={true}>
        <PaperItem content="xs" />
      </Grid>
    );
  });

  const arry2 = [1, 2, 3, 4, 5];

  const gridItems2 = arry2.map(() => {
    return (
      // 这里使用 xs=true 来自动布局
      <Grid item xs>
        <PaperItem content="xs" />
      </Grid>
    );
  });

  return (
    <div style={{ margin: '20px', borderStyle: 'dashed', borderWidth: '1px' }}>
      <Grid container spacing={1}>
        {/* 这里不使用 Grid item container */}
        {gridItems}
        {gridItems2}
      </Grid>
    </div>
  );
}

export function GridItemWidth(): JSX.Element {
  return (
    <div style={{ margin: '20px', borderStyle: 'dashed', borderWidth: '1px' }}>
      {/* Grid container 不要设置 columnGap */}
      <Grid container justifyContent="space-between" spacing={1}>
        <Grid item width="30%">
          <PaperItem content="width=30%" />
        </Grid>
        <Grid item width="40%">
          <PaperItem content="width=40%" />
        </Grid>
        <Grid item width="30%">
          <PaperItem content="width=30%" />
        </Grid>
        <Grid item width="50%">
          <PaperItem content="width=50%" />
        </Grid>
        <Grid item width="40%">
          <PaperItem content="width=40%" />
        </Grid>
        <Grid item width="10%">
          <PaperItem content="width=10%" />
        </Grid>
      </Grid>
    </div>
  );
}

export function GridItemWidth2(): JSX.Element {
  return (
    <div
      style={{
        minWidth: '720px', // 最小宽度 720 > 140*5
        margin: '20px',
        borderStyle: 'dashed',
        borderWidth: '1px',
      }}
    >
      {/* Grid container 不要设置 columnGap */}
      <Grid container justifyContent="space-between">
        <Grid item width="140px">
          <PaperItem content="140px" />
        </Grid>
        <Grid item width="140px">
          <PaperItem content="140px" />
        </Grid>
        <Grid item width="140px">
          <PaperItem content="140px" />
        </Grid>
        <Grid item width="140px">
          <PaperItem content="140px" />
        </Grid>
        <Grid item width="140px">
          <PaperItem content="140px" />
        </Grid>
      </Grid>
    </div>
  );
}

export function GridComplexLayout(): JSX.Element {
  return (
    <div style={{ margin: '20px', borderStyle: 'dashed', borderWidth: '1px' }}>
      <Grid container rowSpacing={10} columnSpacing={1}>
        {/* 自动布局 */}
        <Grid container item xs={12} spacing={3}>
          <Grid item xs>
            <PaperItem content="xs" />
          </Grid>
          <Grid item xs>
            <PaperItem content="xs" />
          </Grid>
          <Grid item xs>
            <PaperItem content="xs" />
          </Grid>
          <Grid item xs>
            <PaperItem content="xs" />
          </Grid>
          <Grid item xs>
            <PaperItem content="xs" />
          </Grid>
        </Grid>

        {/* 横向布局 */}
        <Grid container item xs={8} columnSpacing={3} rowSpacing={1}>
          <Grid item xs={4}>
            <PaperItem content="xs" />
          </Grid>
          <Grid item xs={4}>
            <PaperItem content="xs" />
          </Grid>
          <Grid item xs={4}>
            <PaperItem content="xs" />
          </Grid>
          <Grid item xs={4}>
            <PaperItem content="xs" />
          </Grid>
          <Grid item xs={4}>
            <PaperItem content="xs" />
          </Grid>
          <Grid item xs={4}>
            <PaperItem content="xs" />
          </Grid>
        </Grid>

        {/* 纵向布局 */}
        <Grid container item xs={4} rowSpacing={3}>
          <Grid item xs={12}>
            <PaperItem content="xs" />
          </Grid>
          <Grid item xs={12}>
            <PaperItem content="xs" />
          </Grid>
          <Grid item xs={12}>
            <PaperItem content="xs" />
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}

export function GridBreakPoint(): JSX.Element {
  return (
    <div
      style={{
        width: '300px', // 这里将 grid 宽度定死了
        margin: '20px',
        borderStyle: 'dashed',
        borderWidth: '1px',
      }}
    >
      <Grid container rowSpacing={3} columnSpacing={1}>
        <Grid item xs={6} sm={4} md={3} lg={2}>
          <PaperItem content="xs={6} sm={4} md={3} lg={2}" />
        </Grid>
        <Grid item xs={6} sm={4} md={3} lg={2}>
          <PaperItem content="xs={6} sm={4} md={3} lg={2}" />
        </Grid>
        <Grid item xs={6} sm={4} md={3} lg={2}>
          <PaperItem content="xs={6} sm={4} md={3} lg={2}" />
        </Grid>
        <Grid item xs={6} sm={4} md={3} lg={2}>
          <PaperItem content="xs={6} sm={4} md={3} lg={2}" />
        </Grid>
        <Grid item xs={6} sm={4} md={3} lg={2}>
          <PaperItem content="xs={6} sm={4} md={3} lg={2}" />
        </Grid>
        <Grid item xs={6} sm={4} md={3} lg={2}>
          <PaperItem content="xs={6} sm={4} md={3} lg={2}" />
        </Grid>
      </Grid>
    </div>
  );
}

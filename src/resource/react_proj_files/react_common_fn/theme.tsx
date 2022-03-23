import React from 'react';
import { createTheme, Button, Box, Grid, ThemeProvider, CssBaseline, PaletteMode, Paper } from '@mui/material';

export function MyTheme(): JSX.Element {
  const [mode, setMode] = React.useState<PaletteMode>('light');

  const myTheme = createTheme({
    palette: {
      mode: mode, // 'light' | 'dark'
    },
  });

  return (
    <>
      {/* 启用自定义主题颜色 */}
      <ThemeProvider theme={myTheme}>
        {/* dark mode 启用深色背景需要用到 CssBaseline */}
        <CssBaseline />
        <ColorPalette />
        <Box display="flex" justifyContent="center">
          <Button
            variant="outlined"
            onClick={() => {
              if (mode === 'dark') {
                setMode('light');
                return;
              }
              setMode('dark');
            }}
          >
            Toggle Mode
          </Button>
        </Box>
      </ThemeProvider>
    </>
  );
}

function ColorPalette(): JSX.Element {
  const colorTypes = ['primary', 'secondary', 'error', 'warning', 'info', 'success'];

  const colorStyle = ['.main', '.light', '.dark'];

  // 生成主题颜色
  const colorList = colorTypes.map((typ) =>
    colorStyle.map((style) => (
      <Grid item xs={4} key={typ + style}>
        <Box bgcolor={typ + style} color={`${typ}.contrastText`} p={4}>
          {typ + style}
        </Box>
      </Grid>
    ))
  );

  // 生成 background 颜色
  const bgColorList = ['default', 'paper'].map((bg) => (
    <Grid item xs={6} key={bg}>
      <Box bgcolor={`background.${bg}`} p={4} border={1}>
        <div>{`background.${bg}`}</div>
      </Box>
    </Grid>
  ));

  // 生成 text 颜色
  const textColorList = ['primary', 'secondary', 'disabled'].map((txt) => (
    <Grid item xs={4} key={txt}>
      <Box color={`text.${txt}`} p={4} border={1}>
        <div>{`text.${txt}`}</div>
      </Box>
    </Grid>
  ));

  return (
    <div style={{ padding: '8px' }}>
      <Grid container spacing={1}>
        {colorList}
        {bgColorList}
        {textColorList}
      </Grid>
    </div>
  );
}

import React from 'react';
import CssBaseline from '@material-ui/core/CssBaseline';
import {
  createTheme,
  Button,
  PaletteType,
  ThemeProvider,
  Box,
  Grid,
} from '@material-ui/core';

export function TestPalette(): JSX.Element {
  const [mode, setMode] = React.useState<PaletteType>('light');

  const newTheme = createTheme({
    palette: {
      type: mode,
    },
  });

  return (
    <>
      <ThemeProvider theme={newTheme}>
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

function ColorPalette() {
  const colorTypes = [
    'primary',
    'secondary',
    'error',
    'warning',
    'info',
    'success',
  ];

  const colorStyle = ['.main', '.light', '.dark'];

  const colorList = colorTypes.map((typ) =>
    colorStyle.map((style) => (
      <Grid item xs={4} key={typ + style}>
        <Box bgcolor={typ + style} color={`${typ}.contrastText`} p={4}>
          {typ + style}
        </Box>
      </Grid>
    ))
  );

  const bgColorList = ['default', 'paper'].map((bg) => (
    <Grid item xs={6} key={bg}>
      <Box bgcolor={`background.${bg}`} p={4} border={1}>
        <div>{`background.${bg}`}</div>
      </Box>
    </Grid>
  ));

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

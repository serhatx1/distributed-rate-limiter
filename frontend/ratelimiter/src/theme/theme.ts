import { createTheme } from '@mui/material';

export const theme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#5E6AD2',
      light: '#8E96E3',
      dark: '#4A54C6',
    },
    background: {
      default: '#111111',
      paper: '#1C1C1F',
    },
    text: {
      primary: '#FFFFFF',
      secondary: 'rgba(255, 255, 255, 0.6)',
    },
    divider: 'rgba(255, 255, 255, 0.04)',
  },
  shape: {
    borderRadius: 8,
  },
  typography: {
    fontFamily: 'Inter, -apple-system, BlinkMacSystemFont, sans-serif',
    button: {
      textTransform: 'none',
      fontWeight: 500,
    },
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 6,
          textTransform: 'none',
          fontWeight: 500,
          letterSpacing: '-0.01em',
        },
      },
    },
  },
});

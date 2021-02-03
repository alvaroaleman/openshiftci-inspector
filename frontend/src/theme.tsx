import red from '@material-ui/core/colors/red';
import { createMuiTheme } from '@material-ui/core/styles';
import {green, orange} from "@material-ui/core/colors";

// A custom theme for this app
const theme = createMuiTheme({
  palette: {
    background: {
      default: '#f5f5f5',
    },
    success: {
      main: green.A400,
    },
    error: {
      main: red.A400,
    },
    warning: {
      main: orange.A400,
    },
    primary: {
      main: '#bb0000',
    },
    secondary: {
      main: '#151515',
    },
  },
});

export default theme;

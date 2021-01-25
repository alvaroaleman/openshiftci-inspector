import red from '@material-ui/core/colors/red';
import { createMuiTheme } from '@material-ui/core/styles';
import {green, grey, orange} from "@material-ui/core/colors";

// A custom theme for this app
const theme = createMuiTheme({
  palette: {
    background: {
      default: '#fff',
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
      main: '#556cd6',
    },
    secondary: {
      main: grey.A100,
    },
  },
});

export default theme;

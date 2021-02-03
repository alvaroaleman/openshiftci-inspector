import {Box, createMuiTheme, List, ListItem, ListItemIcon, ListItemText, Typography} from "@material-ui/core";
import DashboardIcon from '@material-ui/icons/Dashboard';
import * as React from "react";
import { useHistory } from "react-router-dom";
import {ThemeProvider} from "@material-ui/core/styles";
import {green, orange} from "@material-ui/core/colors";
import red from "@material-ui/core/colors/red";
import './sidebar.css';

export interface ISidebarProps {
}

export default function Sidebar(props: ISidebarProps) {
    const history = useHistory();
    const theme = createMuiTheme({
        palette: {
            background: {
                default: '#262626',
            },
            text: {
                primary: '#ffffff',
                secondary: '#ffffff',
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
    return <ThemeProvider theme={theme}>
        <Box className="app__sidebar">
            <Box p={2} className="app__logo">
                <img src={"/logo.svg"} alt={"A yellow fedora hat with a yellow propeller on top"} />
                <Typography variant={"h5"} component={"h1"} gutterBottom={true} align={"center"}>Openshift CI Inspector</Typography>
            </Box>
            <List>
                <ListItem button={true}>
                    <ListItemIcon><DashboardIcon color={"primary"} /></ListItemIcon>
                    <ListItemText onClick={function() {history.push("/")}}>
                        Home
                    </ListItemText>
                </ListItem>
            </List>
        </Box>
    </ThemeProvider>
}

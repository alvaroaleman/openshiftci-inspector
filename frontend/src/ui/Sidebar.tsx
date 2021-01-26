import { Divider } from "@blueprintjs/core";
import {Box, List, ListItem, ListItemIcon, ListItemText, Typography} from "@material-ui/core";
import DashboardIcon from '@material-ui/icons/Dashboard';
import * as React from "react";
import { useHistory } from "react-router-dom";

export interface ISidebarProps {
}

export default function Sidebar(props: ISidebarProps) {
    const history = useHistory();
    return <div>
        <Box m={2}>
            <Typography variant={"h5"} component={"h1"} gutterBottom={true}>Openshift CI Inspector</Typography>
        </Box>
        <Divider />
        <List>
            <ListItem button={true}>
                <ListItemIcon><DashboardIcon /></ListItemIcon>
                <ListItemText onClick={function() {history.push("/")}}>
                    Home
                </ListItemText>
            </ListItem>
        </List>
    </div>
}

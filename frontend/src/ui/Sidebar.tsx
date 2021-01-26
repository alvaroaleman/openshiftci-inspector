import { Divider } from "@blueprintjs/core";
import {Box, Link, List, ListItem, ListItemIcon, ListItemText, Typography} from "@material-ui/core";
import DashboardIcon from '@material-ui/icons/Dashboard';
import * as React from "react";
import {Link as RouterLink} from "react-router-dom";

export interface ISidebarProps {
}

class Sidebar extends React.Component<ISidebarProps> {
    render() {
        return <div>
            <Box m={2}>
                <Typography variant={"h5"} component={"h1"} gutterBottom={true}>Openshift CI Inspector</Typography>
            </Box>
            <Divider />
            <List>
                <ListItem button={true}>
                    <ListItemIcon><DashboardIcon /></ListItemIcon>
                    <ListItemText primary={
                        <Link component={RouterLink} to={"/"}>
                            Home
                        </Link>
                    } />
                </ListItem>
            </List>
        </div>
    }
}

export default Sidebar;
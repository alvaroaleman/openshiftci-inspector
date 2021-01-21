import { Divider } from "@blueprintjs/core";
import {Box, List, ListItem, ListItemIcon, ListItemText, Typography} from "@material-ui/core";
import DashboardIcon from '@material-ui/icons/Dashboard';
import * as React from "react";
import LinkFactory from "../router/ui/LinkFactory";

export interface ISidebarProps {
    linkFactory: LinkFactory
}

class Sidebar extends React.Component<ISidebarProps> {
    public render() {
        return <div>
            <Box m={2}>
                <Typography variant={"h5"} component={"h1"} gutterBottom={true}>Openshift CI Inspector</Typography>
            </Box>
            <Divider />
            <List>
                <ListItem button={true}>
                    <ListItemIcon><DashboardIcon /></ListItemIcon>
                    <ListItemText primary={
                        this.props.linkFactory.create(
                            "Home",
                            "sidebar__link",
                            "/"
                        )
                    } />
                </ListItem>
            </List>
        </div>
    }
}

export default Sidebar;
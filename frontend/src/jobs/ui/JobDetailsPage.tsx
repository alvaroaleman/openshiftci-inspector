import React from "react";
import JobsGetService from "../get";
import JobInfoCard from "./components/JobInfoCard";
import {
    AppBar,
    Badge,
    Box,
    Grid,
    List,
    ListItem,
    ListItemIcon,
    ListItemText, Paper,
    Tab,
    Tabs,
    Typography
} from "@material-ui/core";
import JobPreviousCard from "./components/JobPreviousCard";
import JobRelatedCard from "./components/JobRelatedCard";
import JobsGetRelatedService from "../related";
import JobsGetPreviousService from "../previous";
import JobMetricsCard from "./components/JobMetricsCard";
import JobsMetricsService from "../metrics";
import ReportProblemIcon from '@material-ui/icons/ReportProblem';
import CheckIcon from '@material-ui/icons/Check';

interface IJobDetailsState {
    tab: number
}

interface IJobDetailsProps {
    id: string
    jobsGetService: JobsGetService
    jobsGetRelatedService: JobsGetRelatedService
    jobsGetPreviousService: JobsGetPreviousService
    metricsService: JobsMetricsService
}

interface TabPanelProps {
    children?: React.ReactNode;
    dir?: string;
    index: any;
    value: any;
}

function TabPanel(props: TabPanelProps) {
    const { children, value, index, ...other } = props;

    return (
        <div
            role="tabpanel"
            hidden={value !== index}
            id={`full-width-tabpanel-${index}`}
            aria-labelledby={`full-width-tab-${index}`}
            {...other}
        >
            {value === index && (
                <Box m={2}>
                    {children}
                </Box>
            )}
        </div>
    );
}

export default class JobDetailsPage extends React.Component<IJobDetailsProps, IJobDetailsState> {
    constructor(props: IJobDetailsProps) {
        super(props);
        this.state = {
            tab: 0
        }
    }

    render = () => {
        return <React.Fragment>
            <AppBar position="static" color={"secondary"}>
                <Tabs value={this.state.tab} onChange={this.changeTab} indicatorColor={"primary"}>
                    <Tab label="Info"  />
                    <Tab label="Metrics"  />
                    <Tab label={<Badge badgeContent={5} color="error">Operators</Badge>}  />
                </Tabs>
            </AppBar>
            <TabPanel value={this.state.tab} index={0}>
                <Grid container spacing={2}>
                    <Grid item xs={12}>
                        <JobInfoCard id={this.props.id} jobsGetService={this.props.jobsGetService} />
                    </Grid>
                    <Grid item xs={6}>
                        <JobRelatedCard id={this.props.id}  jobsGetRelatedService={this.props.jobsGetRelatedService}/>
                    </Grid>
                    <Grid item xs={6}>
                        <JobPreviousCard id={this.props.id}  jobsGetPreviousService={this.props.jobsGetPreviousService}/>
                    </Grid>
                </Grid>
            </TabPanel>
            <TabPanel value={this.state.tab} index={1}>
                <Grid item xs={12}>
                    <JobMetricsCard id={this.props.id} metricsService={this.props.metricsService} />
                </Grid>
            </TabPanel>
            <TabPanel value={this.state.tab} index={2}>
                <Box display={"flex"} flexDirection={"row"}>
                    <Box mr={2}>
                        <Paper elevation={2}>
                            <List>
                                <ListItem button={true}>
                                    <ListItemIcon color={"error"}><Typography color={"error"}><ReportProblemIcon /></Typography></ListItemIcon>
                                    <ListItemText>machine-api</ListItemText>
                                </ListItem>
                                <ListItem button={true}>
                                    <ListItemIcon color={"error"}><Typography color={"error"}><ReportProblemIcon /></Typography></ListItemIcon>
                                    <ListItemText>authentication</ListItemText>
                                </ListItem>
                                <ListItem button={true}>
                                    <ListItemIcon color={"success"}><CheckIcon /></ListItemIcon>
                                    <ListItemText>cloud-credential</ListItemText>
                                </ListItem>
                                <ListItem button={true}>
                                    <ListItemIcon color={"success"}><CheckIcon /></ListItemIcon>
                                    <ListItemText>cluster-autoscaler</ListItemText>
                                </ListItem>
                                <ListItem button={true}>
                                    <ListItemIcon color={"success"}><CheckIcon /></ListItemIcon>
                                    <ListItemText>config-operator</ListItemText>
                                </ListItem>
                                <ListItem button={true}>
                                    <ListItemIcon color={"success"}><CheckIcon /></ListItemIcon>
                                    <ListItemText>console</ListItemText>
                                </ListItem>
                                <ListItem button={true}>
                                    <ListItemIcon color={"success"}><CheckIcon /></ListItemIcon>
                                    <ListItemText>csi-snapshot-controller</ListItemText>
                                </ListItem>
                            </List>
                        </Paper>
                    </Box>
                    <Box flex={1}>
                        <Paper elevation={2}>
                            <Box p={2}>
                        <pre>{`{
    "apiVersion": "v1",
    "items": [],
    "kind": "List",
    "metadata": {
    "resourceVersion": "",
    "selfLink": ""
}`}</pre>
                            </Box>
                        </Paper>
                    </Box>
                </Box>
            </TabPanel>
        </React.Fragment>
    }

    changeTab = (event: React.ChangeEvent<{}>, newValue: number) => {
        this.setState({
            tab: newValue
        })
    }
}
import React from "react";
import JobsGetService from "../get";
import {Job} from "../../api-client";
import {
    Box,
    Card,
    CardContent, CardHeader,
    CircularProgress,
    Grid,
    Table,
    TableBody, TableCell,
    TableRow,
    Typography
} from "@material-ui/core";
import LaunchIcon from '@material-ui/icons/Launch';
import JobStatus from "./JobStatus";
import JobTime from "./JobTime";

interface IJobDetailsState {
    isLoaded: boolean,
    isRefreshing: boolean,
    job?: Job,
}

interface IJobDetailsProps {
    id: string
    jobsGetService: JobsGetService
}

export default class JobDetails extends React.Component<IJobDetailsProps, IJobDetailsState> {
    constructor(props: IJobDetailsProps) {
        super(props);
        this.state = {
            isLoaded: false,
            isRefreshing: true,
        }
    }

    componentDidMount = async () => {
        this.setState({
            isRefreshing: true
        })
        const job = await this.props.jobsGetService.getJob(this.props.id)
        this.setState({
            isRefreshing: false,
            isLoaded: true,
            job: job,
        })
    }

    componentWillUnmount = () => {
    }

    update = () => {
        this.setState({
        })
    }

    render = () => {
        if (this.state.isRefreshing) {
            return <Box m={4} display="flex" alignItems="center" justifyContent="center"><CircularProgress /></Box>
        }
        if (!this.state.isLoaded) {
            return <div />
        }
        if (!this.state.job) {
            return <Typography component={"h1"}>Job not found</Typography>
        }
        const job = this.state.job
        return <Box m={2}>
            <Card>
                <CardHeader avatar={<JobStatus status={job.status} fontSize={"inherit"} />} title={this.state.job.job} />
                <CardContent>
                    <Grid container spacing={2}>
                        <Grid item xs={6}>
                            <Table size="small">
                                <TableBody>
                                    <TableRow>
                                        <TableCell><strong>Status:</strong></TableCell>
                                        <TableCell>{job.status}</TableCell>
                                    </TableRow>
                                    <TableRow>
                                        <TableCell><strong>Start time:</strong></TableCell>
                                        <TableCell><JobTime time={job.startTime} /></TableCell>
                                    </TableRow>
                                    <TableRow>
                                        <TableCell><strong>Pending time:</strong></TableCell>
                                        <TableCell><JobTime time={job.pendingTime} /></TableCell>
                                    </TableRow>
                                    <TableRow>
                                        <TableCell><strong>Completion time:</strong></TableCell>
                                        <TableCell><JobTime time={job.completionTime} /></TableCell>
                                    </TableRow>
                                </TableBody>
                            </Table>
                        </Grid>
                        <Grid item xs={6}>
                            <Table size="small">
                                <TableBody>
                                    <TableRow>
                                        <TableCell><strong>Repository:</strong></TableCell>
                                        <TableCell>{job.gitOrg != null && job.gitRepo != null?<a href={job.gitRepoLink} target="_blank" rel={"noreferrer noopener"}>{job.gitOrg}/{job.gitRepo} <LaunchIcon style={{fontSize:"1em"}} /></a>:null}</TableCell>
                                    </TableRow>
                                    <TableRow>
                                        <TableCell><strong>Base ref:</strong></TableCell>
                                        <TableCell>{job.gitBaseRef}</TableCell>
                                    </TableRow>
                                    <TableRow>
                                        <TableCell><strong>Pulls:</strong></TableCell>
                                        <TableCell>{this.getPulls()}</TableCell>
                                    </TableRow>
                                    <TableRow>
                                        <TableCell><strong>Links:</strong></TableCell>
                                        <TableCell>
                                            <a href={job.url}  target="_blank" rel={"noreferrer noopener"}>Prow <LaunchIcon style={{fontSize:"1em"}} /></a>
                                        </TableCell>
                                    </TableRow>
                                </TableBody>
                            </Table>
                        </Grid>
                    </Grid>
                </CardContent>
            </Card>
        </Box>
    }

    getPulls = () => {
        const job = this.state.job
        if (job?.pulls == null) {
            return null
        }
        return job.pulls.map(pull => {
            return <span><a href={pull.pullLink} target="_blank" rel={"noreferrer noopener"}>{pull.number} <LaunchIcon style={{fontSize:"1em"}} /></a> by <a href={pull.authorLink} target="_blank" rel={"noreferrer noopener"}>{pull.author} <LaunchIcon style={{fontSize:"1em"}} /></a></span>
        })
    }
}
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
                            <Table>
                                <TableBody>
                                    <TableRow>
                                        <TableCell>Status:</TableCell>
                                        <TableCell>{job.status}</TableCell>
                                    </TableRow>
                                    <TableRow>
                                        <TableCell>Start time:</TableCell>
                                        <TableCell>{job.startTime}</TableCell>
                                    </TableRow>
                                    <TableRow>
                                        <TableCell>Pending time:</TableCell>
                                        <TableCell>{job.pendingTime}</TableCell>
                                    </TableRow>
                                    <TableRow>
                                        <TableCell>Completion time:</TableCell>
                                        <TableCell>{job.completionTime}</TableCell>
                                    </TableRow>
                                </TableBody>
                            </Table>
                        </Grid>
                        <Grid item xs={6}>
                            <Table>
                                <TableBody>
                                    <TableRow>
                                        <TableCell>Repository:</TableCell>
                                        <TableCell>{job.gitOrg != null && job.gitRepo != null?<a href={job.gitRepoLink} target="_blank" rel={"noreferrer noopener"}>{job.gitOrg}/{job.gitRepo} <LaunchIcon style={{fontSize:"1em"}} /></a>:null}</TableCell>
                                    </TableRow>
                                    <TableRow>
                                        <TableCell>Base ref:</TableCell>
                                        <TableCell>{job.gitBaseRef}</TableCell>
                                    </TableRow>
                                    <TableRow>
                                        <TableCell>Pulls:</TableCell>
                                        <TableCell>{this.getPulls()}</TableCell>
                                    </TableRow>
                                    <TableRow>
                                        <TableCell>Prow:</TableCell>
                                        <TableCell><a href={job.url}  target="_blank" rel={"noreferrer noopener"}>Open <LaunchIcon style={{fontSize:"1em"}} /></a></TableCell>
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
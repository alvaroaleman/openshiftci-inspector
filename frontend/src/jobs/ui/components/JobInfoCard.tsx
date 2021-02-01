import React from "react";
import JobsGetService from "../../get";
import {JobWithAssetURL} from "../../../api-client";
import {
    Box, Button,
    Card, CardActions,
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
import JobDuration from "./JobDuration";

interface IJobInfoCardState {
    isLoaded: boolean,
    isRefreshing: boolean,
    job?: JobWithAssetURL,
}

interface IJobInfoCardProps {
    id: string
    jobsGetService: JobsGetService
}

export default class JobInfoCard extends React.Component<IJobInfoCardProps, IJobInfoCardState> {
    constructor(props: IJobInfoCardProps) {
        super(props);
        this.state = {
            isLoaded: false,
            isRefreshing: true,
        }
    }

    componentDidMount = () => {
        this.reload()
    }

    componentDidUpdate = (prevProps: Readonly<IJobInfoCardProps>, prevState: Readonly<IJobInfoCardState>, snapshot?: any) => {
        if (prevProps.id !== this.props.id) {
            this.reload()
        }
    }

    reload = async () => {
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
        return <React.Fragment>
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
                                </TableBody>
                            </Table>
                        </Grid>
                        <Grid item xs={6}>
                            <Table size="small">
                                <TableBody>
                                    <TableRow>
                                        <TableCell><strong>Started:</strong></TableCell>
                                        <TableCell><JobTime time={job.startTime} /></TableCell>
                                    </TableRow>
                                    <TableRow>
                                        <TableCell><strong>Pending:</strong></TableCell>
                                        <TableCell><JobTime time={job.pendingTime} /></TableCell>
                                    </TableRow>
                                    <TableRow>
                                        <TableCell><strong>Completed:</strong></TableCell>
                                        <TableCell><JobTime time={job.completionTime} /></TableCell>
                                    </TableRow>
                                    <TableRow>
                                        <TableCell><strong>Total time:</strong></TableCell>
                                        <TableCell><JobDuration startTime={job.pendingTime} completionTime={job.completionTime} /></TableCell>
                                    </TableRow>

                                </TableBody>
                            </Table>
                        </Grid>
                    </Grid>
                </CardContent>
                <CardActions>
                    <a href={job.url}  target="_blank" rel={"noreferrer noopener"}>
                        <Button variant="contained" color="primary" size={"small"}>Open in Prow&nbsp;<LaunchIcon style={{fontSize:"1em"}} /></Button>
                    </a>
                    {!job.assetURL?null:
                        <a href={job.assetURL}  target="_blank" rel={"noreferrer noopener"}>
                            <Button variant="contained" color="primary" size={"small"}>Artifacts&nbsp;<LaunchIcon style={{fontSize:"1em"}} /></Button>
                        </a>
                    }
                </CardActions>
            </Card>
        </React.Fragment>
    }

    getPulls = () => {
        const job = this.state.job
        if (job?.pulls == null) {
            return null
        }
        return job.pulls.map(pull => {
            return <span key={pull.number}><a href={pull.pullLink} target="_blank" rel={"noreferrer noopener"}>{pull.number} <LaunchIcon style={{fontSize:"1em"}} /></a> by <a href={pull.authorLink} target="_blank" rel={"noreferrer noopener"}>{pull.author} <LaunchIcon style={{fontSize:"1em"}} /></a></span>
        })
    }
}
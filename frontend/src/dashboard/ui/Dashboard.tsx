import {
    Box,
    LinearProgress,
    Paper,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow, TextField,
} from "@material-ui/core";
import ErrorIcon from '@material-ui/icons/Error';
import HourglassEmptyIcon from '@material-ui/icons/HourglassEmpty';
import React from "react";
import JobsListService from "../../jobs/list";
import {Job} from "../../api-client";

interface IDashboardState {
    jobFilter: string,
    repoFilter: string,
    isLoaded: boolean,
    isRefreshing: boolean,
    jobs: Array<Job>
}

interface IDashboardProps {
    jobsListService: JobsListService,
}

class Dashboard extends React.Component<IDashboardProps, IDashboardState> {

    constructor(props: IDashboardProps) {
        super(props);
        this.state = {
            jobFilter: "",
            repoFilter: "",
            isLoaded: false,
            isRefreshing: false,
            jobs: new Array<Job>()
        }
    }

    componentDidMount = () => {
        this.props.jobsListService.register(this)
        // noinspection JSIgnoredPromiseFromCall
        this.props.jobsListService.refresh()
    }

    componentWillUnmount = () => {
        this.props.jobsListService.deregister(this)
    }

    update = () => {
        this.setState({
            isLoaded: this.props.jobsListService.isLoaded(),
            isRefreshing: this.props.jobsListService.isRefreshing(),
            jobs: this.props.jobsListService.getJobs()
        })
    }

    changeJobFilter = (e: React.ChangeEvent<HTMLInputElement>) => {
        this.setState({
            jobFilter: e.target.value,
        })
    }

    changeRepoFilter = (e: React.ChangeEvent<HTMLInputElement>) => {
        this.setState({
            repoFilter: e.target.value,
        })
    }

    render = () => {
        return <Box m={2}>
            <TableContainer component={Paper}>
                {this.state.isRefreshing?<LinearProgress />:null}
                <Table aria-label="simple table" size="small">
                    <TableHead>
                        <TableRow>
                            <TableCell>Job</TableCell>
                            <TableCell>Git repository</TableCell>
                            <TableCell>Base</TableCell>
                            <TableCell>Pulls</TableCell>
                        </TableRow>
                        <TableRow>
                            <TableCell><TextField id="filter-job" label="Filter by job" size={"small"} value={this.state.jobFilter} onChange={this.changeJobFilter} /></TableCell>
                            <TableCell><TextField id="filter-repo" label="Filter by repository" size={"small"} value={this.state.repoFilter} /></TableCell>
                            <TableCell />
                            <TableCell />
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {this.state.jobs.filter(job => {
                            if (this.state.jobFilter !== "" && !job.job.includes(this.state.jobFilter)) {
                                return false
                            }
                            // noinspection RedundantIfStatementJS
                            if (job.gitOrg !== null &&
                                job.gitRepo !== null &&
                                this.state.repoFilter !== "" &&
                                !(job.gitOrg + "/" + job.gitRepo).includes(this.state.repoFilter)) {
                                return false
                            }
                            return true
                        }).map(job => {
                                return <TableRow id={job.id}>
                                    <TableCell style={{
                                        display: 'flex',
                                        alignItems: 'center',
                                        flexWrap: 'wrap',
                                    }}>{job.job} {this.getJobStatus(job.status)}</TableCell>
                                    <TableCell>{job.gitOrg != null && job.gitRepo != null?<a href={job.gitRepoLink} target="_blank" rel={"noreferrer noopener"}>{job.gitOrg}/{job.gitRepo}</a>:null}</TableCell>
                                    <TableCell>{job.gitBaseRef}</TableCell>
                                    <TableCell>{this.getPulls(job)}</TableCell>
                                </TableRow>
                            })
                        }
                    </TableBody>
                </Table>
            </TableContainer>
        </Box>;
    }

    getJobStatus = (status: string) => {
        switch (status) {
            case "success":
                return null
            case "failure":
                return <Box color="error.main" component={"span"} fontSize="small" style={{ marginLeft: "0.5rem"}}><ErrorIcon /></Box>
            case "pending":
                return <Box color="warning.main" component={"span"} fontSize="small" style={{ marginLeft: "0.5rem"}}><HourglassEmptyIcon /></Box>
            default:
                return null
        }
    }

    getPulls = (job: Job) => {
        if (job.pulls == null) {
            return null
        }
        return job.pulls.map(pull => {
            return <a href={pull.pullLink} target="_blank" rel={"noreferrer noopener"}>{pull.number}</a>
        })
    }
}

export default Dashboard;
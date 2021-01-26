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
import React from "react";
import JobsListService from "../list";
import {Job} from "../../api-client";
import JobStatus from "./components/JobStatus";
import Link from "../../common/Link"

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

export default class JobDashboardPage extends React.Component<IDashboardProps, IDashboardState> {

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
                <Table size="small">
                    <TableHead>
                        <TableRow>
                            <TableCell>Job</TableCell>
                            <TableCell>Git repository</TableCell>
                            <TableCell>Base</TableCell>
                            <TableCell>Pulls</TableCell>
                        </TableRow>
                        <TableRow>
                            <TableCell colSpan={4}>
                                <Box mr={2} component={"span"}><TextField id="filter-job" label="Filter by job" size={"small"} value={this.state.jobFilter} onChange={this.changeJobFilter} /></Box>
                                <Box mr={2} component={"span"}><TextField id="filter-repo" label="Filter by repository" size={"small"} value={this.state.repoFilter} onChange={this.changeRepoFilter} /></Box>
                            </TableCell>
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
                                return <TableRow key={job.id}>
                                    <TableCell style={{
                                        display: 'flex',
                                        alignItems: 'center',
                                        flexWrap: 'wrap',
                                    }}>
                                        <Box component={"span"} mr={1}><JobStatus status={job.status} fontSize={"inherit"} /></Box>
                                        <Link to={"/" + job.id}>
                                            {job.job}
                                        </Link>
                                    </TableCell>
                                    <TableCell>
                                        {job.gitOrg != null && job.gitRepo != null?<a href={job.gitRepoLink} target="_blank" rel={"noreferrer noopener"}>{job.gitOrg}/{job.gitRepo}</a>:null}
                                    </TableCell>
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

    getPulls = (job: Job) => {
        if (job.pulls == null) {
            return null
        }
        return job.pulls.map(pull => {
            return <span key={pull.number}><a href={pull.pullLink} target="_blank" rel={"noreferrer noopener"}>{pull.number}</a> by <a href={pull.authorLink} target="_blank" rel={"noreferrer noopener"}>{pull.author}</a></span>
        })
    }
}
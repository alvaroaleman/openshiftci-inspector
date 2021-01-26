import {Job} from "../../../api-client";
import React from "react";
import JobsGetRelatedService from "../../related";
import {
    Box,
    Card,
    CardHeader,
    LinearProgress, Table, TableBody,
    TableCell,
    TableHead,
    TableRow,
} from "@material-ui/core";
import JobStatus from "./JobStatus";
import JobTime from "./JobTime";
import Link from "../../../common/Link";

interface IJobRelatedCardState {
    isLoaded: boolean,
    isRefreshing: boolean,
    jobs: Array<Job>,
}

interface IJobRelatedCardProps {
    id: string
    jobsGetRelatedService: JobsGetRelatedService
}

export default class JobRelatedCard extends React.Component<IJobRelatedCardProps, IJobRelatedCardState> {
    constructor(props: IJobRelatedCardProps) {
        super(props);
        this.state = {
            isLoaded: false,
            isRefreshing: true,
            jobs: new Array<Job>()
        }
    }

    componentDidMount = () => {
        this.reload()
    }

    componentDidUpdate = (prevProps: Readonly<IJobRelatedCardProps>, prevState: Readonly<IJobRelatedCardState>, snapshot?: any) => {
        if (prevProps.id !== this.props.id) {
            this.reload()
        }
    }

    reload = async () => {
        this.setState({
            isRefreshing: true
        })
        const jobs = await this.props.jobsGetRelatedService.getJob(this.props.id)
        this.setState({
            isRefreshing: false,
            isLoaded: true,
            jobs: jobs,
        })
    }

    render = () => {
        return <React.Fragment>
            <Card>
                {this.state.isRefreshing?<LinearProgress />:null}
                <CardHeader title="Related jobs" subheader={"Lists other jobs for the same repository, base, and pull request ID."} />
                <Table size="small">
                    <TableHead>
                        <TableRow>
                            <TableCell>Job</TableCell>
                            <TableCell>Start time</TableCell>
                            <TableCell>Git repository</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {this.state.jobs.map(job => {
                            return <TableRow key={job.id}>
                                <TableCell>
                                    <Box component={"span"} mr={1}><JobStatus status={job.status} fontSize={"inherit"} /></Box>
                                    <Link to={"/" + job.id}>
                                        {job.job}
                                    </Link>
                                </TableCell>
                                <TableCell>
                                    <JobTime time={job.startTime} />
                                </TableCell>
                                <TableCell>
                                    {job.gitOrg != null && job.gitRepo != null?<a href={job.gitRepoLink} target="_blank" rel={"noreferrer noopener"}>{job.gitOrg}/{job.gitRepo}</a>:null}
                                </TableCell>
                            </TableRow>
                        })}
                    </TableBody>
                </Table>
            </Card>
        </React.Fragment>
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
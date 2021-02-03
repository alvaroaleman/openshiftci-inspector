import {Job} from "../../../api-client";
import React from "react";
import JobsGetPreviousService from "../../previous";
import {
    Box,
    Card,
    CardHeader,
    LinearProgress, Table, TableBody,
    TableCell, TableContainer,
    TableHead,
    TableRow,
} from "@material-ui/core";
import JobStatus from "./JobStatus";
import JobTime from "./JobTime";
import Link from "../../../common/Link";

interface IJobPreviousCardState {
    isLoaded: boolean,
    isRefreshing: boolean,
    jobs: Array<Job>,
}

interface IJobPreviousCardProps {
    id: string
    jobsGetPreviousService: JobsGetPreviousService
}

export default class JobPreviousCard extends React.Component<IJobPreviousCardProps, IJobPreviousCardState> {
    constructor(props: IJobPreviousCardProps) {
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

    componentDidUpdate = (prevProps: Readonly<IJobPreviousCardProps>, prevState: Readonly<IJobPreviousCardState>, snapshot?: any) => {
        if (prevProps.id !== this.props.id) {
            this.reload()
        }
    }

    reload = async () => {
        this.setState({
            isRefreshing: true
        })
        const jobs = await this.props.jobsGetPreviousService.getJob(this.props.id)
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
                <CardHeader title="Previous jobs" subheader={"Lists previous jobs for the same job type and repository."} />
                <TableContainer style={{height:"500px"}}>
                    <Table size="small" stickyHeader>
                        <TableHead>
                            <TableRow>
                                <TableCell>Time</TableCell>
                                <TableCell>Pulls</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {this.state.jobs.map(job => {
                                return <TableRow key={job.id}>
                                    <TableCell>
                                        <Box display={"flex"} flexDirection={"row"}>
                                            <Box component={"span"} mr={1} style={{paddingTop:"2px"}}>
                                                <JobStatus status={job.status} fontSize={"inherit"} />
                                            </Box>
                                            <Box flex={1}>
                                                <Link to={"/" + job.id}>
                                                    <JobTime time={job.startTime} />
                                                </Link>
                                            </Box>
                                        </Box>
                                    </TableCell>
                                    <TableCell>{this.getPulls(job)}</TableCell>
                                </TableRow>
                            })}
                        </TableBody>
                    </Table>
                </TableContainer>
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
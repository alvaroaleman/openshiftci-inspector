import {
    Box, Button,
    LinearProgress, ListItemIcon, ListItemText, Menu, MenuItem,
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
import JobTime from "./components/JobTime";
import SearchIcon from '@material-ui/icons/Search';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import ImportContactsIcon from '@material-ui/icons/ImportContacts';
import { useHistory } from "react-router-dom";

interface IPullsProps {
    job: Job
}

function Pulls(props: IPullsProps) {
    if (props.job.pulls == null) {
        return null
    }
    return <React.Fragment>{props.job.pulls.map(pull => {
        return <span key={pull.number}><a href={pull.pullLink} target="_blank" rel={"noreferrer noopener"}>{pull.number}</a> by <a href={pull.authorLink} target="_blank" rel={"noreferrer noopener"}>{pull.author}</a></span>
    })}</React.Fragment>
}

interface IRowProps {
    job: Job,
    onSearch: (query: string, repository: string) => void
}

function Row(props: IRowProps) {
    const history = useHistory();

    const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);

    const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
        setAnchorEl(event.currentTarget);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };

    const job = props.job

    return (
        <TableRow key={job.id}>
            <TableCell>
                <Box display={"flex"} flexDirection={"row"}>
                    <Box component={"span"} mr={1} style={{paddingTop:"4px"}}>
                        <JobStatus status={job.status} fontSize={"inherit"} />
                    </Box>
                    <Box flex={1}>
                        <Link to={"/" + job.id} title={"Click to show details page..."}>
                            {job.job}
                        </Link>
                    </Box>
                </Box>
            </TableCell>
            <TableCell>
                <JobTime time={job.startTime} />
            </TableCell>
            <TableCell>
                {job.gitOrg != null && job.gitRepo != null?<a href={job.gitRepoLink} target="_blank" rel={"noreferrer noopener"}>{job.gitOrg}/{job.gitRepo}</a>:null}
            </TableCell>
            <TableCell>{job.gitBaseRef}</TableCell>
            <TableCell><Pulls job={job} /></TableCell>
            <TableCell width={"2rem"}>
                <Box display={"inline-block"} ml={1}>
                    <Button size={"small"} aria-controls="simple-menu" aria-haspopup="true" onClick={handleClick}>
                        <MoreVertIcon fontSize={"small"} />
                    </Button>
                </Box>
                <Menu
                    anchorEl={anchorEl}
                    keepMounted
                    open={Boolean(anchorEl)}
                    onClose={handleClose}
                >
                    <MenuItem onClick={() => {
                        handleClose()
                        history.push("/" + job.id)
                    }}>
                        <ListItemIcon><ImportContactsIcon /></ListItemIcon>
                        <ListItemText>Open Details Page...</ListItemText>
                    </MenuItem>
                    <MenuItem onClick={() => {
                        handleClose()
                        props.onSearch(job.job, "")
                    }}>
                        <ListItemIcon><SearchIcon /></ListItemIcon>
                        <ListItemText>Show only &ldquo;<code>{job.job}</code>&rdquo; jobs</ListItemText>
                    </MenuItem>
                    {job.gitOrg != null && job.gitRepo != null?<MenuItem onClick={() => {
                        handleClose()
                        props.onSearch("", job.gitOrg as string + "/" + job.gitRepo)
                    }}>
                        <ListItemIcon><SearchIcon /></ListItemIcon>
                        <ListItemText>Show only jobs for the &ldquo;<code>{job.gitOrg}/{job.gitRepo}</code>&rdquo; repository</ListItemText>
                    </MenuItem>:null}
                    {job.gitOrg != null && job.gitRepo != null?<MenuItem onClick={() => {
                        handleClose()
                        props.onSearch(job.job, job.gitOrg as string + "/" + job.gitRepo)
                    }}>
                        <ListItemIcon><SearchIcon /></ListItemIcon>
                        <ListItemText>Show only &ldquo;<code>{job.job}</code>&rdquo; jobs for the &ldquo;<code>{job.gitOrg}/{job.gitRepo}</code>&rdquo; repository</ListItemText>
                    </MenuItem>:null}
                </Menu>
            </TableCell>
        </TableRow>
    );
}

interface IDashboardState {
    jobFilter: string,
    repoFilter: string,
    typingTimer?: number,
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

    search = () => {
        if (this.state.typingTimer) {
            window.clearTimeout(this.state.typingTimer)
        }
        this.setState({typingTimer:undefined})
        this.props.jobsListService.setFilters(this.state.jobFilter, this.state.repoFilter)
    }

    searchFor = (keyword: string, repository: string) => {
        if (this.state.typingTimer) {
            window.clearTimeout(this.state.typingTimer)
        }

        this.setState({
            jobFilter: keyword,
            repoFilter: repository,
            typingTimer:undefined,
        })
        this.props.jobsListService.setFilters(keyword, repository)
    }

    changeJobFilter = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (this.state.typingTimer) {
            window.clearTimeout(this.state.typingTimer)
        }
        this.setState({
            jobFilter: e.target.value,
            typingTimer: window.setTimeout(this.search, 300)
        })
    }

    changeRepoFilter = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (this.state.typingTimer) {
            window.clearTimeout(this.state.typingTimer)
        }
        this.setState({
            repoFilter: e.target.value,
            typingTimer: window.setTimeout(this.search, 300)
        })
    }

    render = () => {
        return <Box m={2}>
            <h1>Last job runs</h1>
            <TableContainer component={Paper}>
                {this.state.isRefreshing?<LinearProgress />:null}
                <Table size="small">
                    <TableHead>
                        <TableRow>
                            <TableCell>Job</TableCell>
                            <TableCell>Started</TableCell>
                            <TableCell>Git repository</TableCell>
                            <TableCell>Base</TableCell>
                            <TableCell>Pulls</TableCell>
                            <TableCell width={"2rem"} />
                        </TableRow>
                        <TableRow>
                            <TableCell colSpan={4}>
                                <Box mr={2} component={"span"}><TextField id="filter-job" label="Filter by job" size={"small"} value={this.state.jobFilter} onChange={this.changeJobFilter} /></Box>
                                <Box mr={2} component={"span"}><TextField id="filter-repo" label="Filter by repository" size={"small"} value={this.state.repoFilter} onChange={this.changeRepoFilter} /></Box>
                            </TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {this.state.jobs.map(job => <Row job={job} onSearch={this.searchFor} />)}
                    </TableBody>
                </Table>
            </TableContainer>
        </Box>;
    }
}
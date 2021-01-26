import React from "react";
import JobsGetService from "../get";
import JobInfoCard from "./components/JobInfoCard";
import {Box, Grid} from "@material-ui/core";
import JobPreviousCard from "./components/JobPreviousCard";
import JobRelatedCard from "./components/JobRelatedCard";
import JobsGetRelatedService from "../related";
import JobsGetPreviousService from "../previous";

interface IJobDetailsState {
}

interface IJobDetailsProps {
    id: string
    jobsGetService: JobsGetService
    jobsGetRelatedService: JobsGetRelatedService
    jobsGetPreviousService: JobsGetPreviousService
}

export default class JobDetailsPage extends React.Component<IJobDetailsProps, IJobDetailsState> {
    constructor(props: IJobDetailsProps) {
        super(props);
        this.state = {
        }
    }

    render = () => {
        return <React.Fragment>
            <Box m={2}>
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
            </Box>
        </React.Fragment>
    }
}
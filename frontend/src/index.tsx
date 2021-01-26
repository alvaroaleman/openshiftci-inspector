import CssBaseline from '@material-ui/core/CssBaseline';
import { ThemeProvider } from '@material-ui/core/styles';
import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import JobDashboardPage from "./jobs/ui/JobDashboardPage";
import './index.css';
import NotificationServiceFactory from "./notification/service/NotificationServiceFactory";
import ToastHandlerFactory from "./notification/ui/ToastHandlerFactory";
import theme from './theme';
import {Configuration, JobsApi} from "./api-client";
import JobsListService from "./jobs/list";
import JobsGetService from "./jobs/get";
import JobDetailsPage from "./jobs/ui/JobDetailsPage";
import {
    BrowserRouter as Router,
    Switch,
    Route,
} from "react-router-dom";
import {RouteComponentProps} from "react-router";
import JobsGetPreviousService from "./jobs/previous";
import JobsGetRelatedService from "./jobs/related";

const notificationServiceFactory = new NotificationServiceFactory();
const toastHandlerFactory = new ToastHandlerFactory(
    notificationServiceFactory
);

const baseURL = window.location.protocol + "//" + window.location.hostname + (
    (window.location.protocol === "http:" && window.location.port !== "80") || (window.location.protocol === "https:" && window.location.port !== "443") ? ":" + window.location.port:""
)

const apiConfiguration = new Configuration({
    basePath: baseURL
})
const jobsAPI = new JobsApi(apiConfiguration)
const jobsListService = new JobsListService(jobsAPI, notificationServiceFactory.create())
const jobsGetService = new JobsGetService(jobsAPI, notificationServiceFactory.create())
const jobsGetPreviousService = new JobsGetPreviousService(jobsAPI, notificationServiceFactory.create())
const jobsGetRelatedService = new JobsGetRelatedService(jobsAPI, notificationServiceFactory.create())

ReactDOM.render(
    <ThemeProvider theme={theme}>
        <CssBaseline />
        <Router>
            <App
                toastHandler={toastHandlerFactory.create()}
            >
                <Switch>
                    <Route exact path="/">
                        <JobDashboardPage jobsListService={jobsListService} />
                    </Route>
                    <Route exact path="/:id" component={jobDetailsRoute} />
                </Switch>
            </App>
        </Router>
    </ThemeProvider>,
    document.getElementById('root') as HTMLElement
);

function jobDetailsRoute(props: RouteComponentProps<any>) {
    return <JobDetailsPage
        jobsGetService={jobsGetService}
        id={props.match.params.id}
        jobsGetPreviousService={jobsGetPreviousService}
        jobsGetRelatedService={jobsGetRelatedService}
    />
}
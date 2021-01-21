import CssBaseline from '@material-ui/core/CssBaseline';
import { ThemeProvider } from '@material-ui/core/styles';
import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import Dashboard from "./dashboard/ui/Dashboard";
import './index.css';
import NotificationServiceFactory from "./notification/service/NotificationServiceFactory";
import ToastHandlerFactory from "./notification/ui/ToastHandlerFactory";
import BrowserHistoryServiceFactory from "./router/service/BrowserHistoryServiceFactory";
import PageTitleServiceFactory from "./router/service/PageTitleServiceFactory";
import RoutingServiceFactory from "./router/service/RoutingServiceFactory";
import LinkFactory from "./router/ui/LinkFactory";
import RouterFactory from "./router/ui/RouterFactory";
import theme from './theme';
import SidebarFactory from "./ui/SidebarFactory";

const notificationServiceFactory = new NotificationServiceFactory();
const toastHandlerFactory = new ToastHandlerFactory(
    notificationServiceFactory
);

const routingServiceFactory = new RoutingServiceFactory(window.location.pathname);

const titleServiceFactory = new PageTitleServiceFactory(
    window
);

const routerFactory = new RouterFactory(
    routingServiceFactory.create(),
    titleServiceFactory.create()
);

const browserHistoryServiceFactory = new BrowserHistoryServiceFactory(
    window.location.protocol + "//" + window.location.hostname + (
        (window.location.protocol === "http:" && window.location.port !== "80") || (window.location.protocol === "https:" && window.location.port !== "443") ? ":" + window.location.port:""
    ),
    window,
    routingServiceFactory.create()
);
browserHistoryServiceFactory.create().register();

const linkFactory = new LinkFactory(routingServiceFactory.create());

const sidebarFactory = new SidebarFactory(linkFactory);

ReactDOM.render(
    <ThemeProvider theme={theme}>
        <CssBaseline />
        <App
            toastHandler={toastHandlerFactory.create()}
            sidebar={sidebarFactory.create()}
        >
            {routerFactory.create("/", "Dashboard", <Dashboard/>)}
        </App>
    </ThemeProvider>,
    document.getElementById('root') as HTMLElement
);

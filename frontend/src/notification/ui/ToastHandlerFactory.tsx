import NotificationServiceFactory from "../service/NotificationServiceFactory";
import ToastHandler from "./ToastHandler";
// @ts-ignore
import React from "react";

class ToastHandlerFactory {
    public constructor(
        readonly notificationServiceFactory: NotificationServiceFactory
    ) {

    }

    public create():JSX.Element {
        return <ToastHandler notificationService={this.notificationServiceFactory.create()}/>
    }
}

export default ToastHandlerFactory;
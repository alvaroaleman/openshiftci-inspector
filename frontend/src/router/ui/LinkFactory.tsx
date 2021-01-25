import RoutingService from "../service/RoutingService";
import Link from "./Link";
// @ts-ignore
import * as React from "react";

class LinkFactory {
    constructor(
        private readonly routingService : RoutingService
    ) {
    }

    public create(
        children: undefined|string|JSX.Element|JSX.Element[],
        className: string|undefined,
        path: string
    ): JSX.Element {
        return <Link
            className={className}
            path={path}
            routingService={this.routingService}
        >{children}</Link>
    }
}

export default LinkFactory;
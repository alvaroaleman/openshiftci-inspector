import * as React from "react";
import {Link as MUILink} from "@material-ui/core"
import {Link as RouterLink} from "react-router-dom";

interface ILinkProps {
    to: string
    children: string|JSX.Element|JSX.Element[]
    title?: string
}

export default class Link extends React.Component<ILinkProps> {
    render = () => {
        return <MUILink component={RouterLink} to={this.props.to} title={this.props.title}>
            {this.props.children}
        </MUILink>
    }
}
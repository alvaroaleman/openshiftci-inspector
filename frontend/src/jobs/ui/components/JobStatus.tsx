import {Tooltip} from "@material-ui/core";
import CheckCircleIcon from "@material-ui/icons/CheckCircle";
import ErrorIcon from "@material-ui/icons/Error";
import WatchLaterIcon from "@material-ui/icons/WatchLater";
import CancelIcon from "@material-ui/icons/Cancel";
import HelpIcon from "@material-ui/icons/Help";
import React from "react";
import theme from "../../../theme"

export interface IJobStatusProps{
    status: string
    fontSize: 'inherit' | 'default' | 'small' | 'large'
}

export default function JobStatus(props: IJobStatusProps) {
    const status = props.status
    switch (status) {
        case "success":
            return <Tooltip title="Success" aria-label="success">
                <CheckCircleIcon htmlColor={theme.palette.success.main} fontSize={props.fontSize} />
            </Tooltip>
        case "failure":
            return <Tooltip title="Failure" aria-label="failure"><ErrorIcon htmlColor={theme.palette.error.main} fontSize={props.fontSize} /></Tooltip>
        case "pending":
            return <Tooltip title="Pending" aria-label="pending"><WatchLaterIcon htmlColor={theme.palette.warning.main} fontSize={props.fontSize} /></Tooltip>
        case "aborted":
            return <Tooltip title="Aborted" aria-label="aborted"><CancelIcon htmlColor={theme.palette.secondary.main} fontSize={props.fontSize}  /></Tooltip>
        default:
            return <Tooltip title={ status }><HelpIcon htmlColor={theme.palette.secondary.main} fontSize={props.fontSize}  /></Tooltip>
    }
}
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
            return <CheckCircleIcon htmlColor={theme.palette.success.main} fontSize={props.fontSize} />
        case "failure":
            return <ErrorIcon htmlColor={theme.palette.error.main} fontSize={props.fontSize} />
        case "pending":
            return <WatchLaterIcon htmlColor={theme.palette.warning.main} fontSize={props.fontSize} />
        case "aborted":
            return <CancelIcon htmlColor={theme.palette.secondary.main} fontSize={props.fontSize}  />
        default:
            return <HelpIcon htmlColor={theme.palette.secondary.main} fontSize={props.fontSize}  />
    }
}
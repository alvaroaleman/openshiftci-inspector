import {Tooltip} from "@material-ui/core";
import * as React from "react";

export interface IJobTimeProps{
    startTime?: string
    completionTime?: string
}

function timeDuration(seconds: number) {
    const years = Math.floor(seconds / 31536000);
    if (years >= 1) {
        return years + " years";
    }
    const months = Math.floor(seconds / 2592000);
    if (months >= 1) {
        return months + " months";
    }
    const days = Math.floor(seconds / 86400);
    if (days >= 1) {
        return days + " days";
    }
    const hours = Math.floor(seconds / 3600);
    if (hours >= 1) {
        return hours + " hours";
    }
    const minutes = Math.floor(seconds / 60);
    if (minutes >= 1) {
        return minutes + " minutes";
    }
    return Math.floor(seconds) + " seconds";
}

export default function JobDuration(props: IJobTimeProps) {
    if (!props.startTime || props.startTime === "") {
        return <span>&mdash;</span>
    }
    if (!props.completionTime || props.completionTime === "") {
        return <span>&mdash;</span>
    }
    const parsedStartTime = Date.parse(props.startTime)
    const parsedCompletionTime = Date.parse(props.completionTime)
    const seconds = Math.floor((parsedCompletionTime - parsedStartTime) / 1000);
    return <Tooltip title={seconds + " seconds"}><span><span style={{whiteSpace:"nowrap"}}>{timeDuration(seconds)}</span></span></Tooltip>
}

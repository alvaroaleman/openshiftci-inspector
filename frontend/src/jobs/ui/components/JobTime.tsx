import {Tooltip} from "@material-ui/core";
import * as React from "react";

export interface IJobTimeProps{
    time?: string
}

function timeSince(currentDate: number, lookupTime: number) {

    const seconds = Math.floor((currentDate - lookupTime) / 1000);

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
    const minutes = seconds / 60;
    if (minutes >= 1) {
        return minutes + " minutes";
    }
    return Math.floor(seconds) + " seconds";
}

export default function JobTime(props: IJobTimeProps) {
    if (!props.time || props.time === "") {
        return <span>&mdash;</span>
    }
    const parsedTime = Date.parse(props.time)
    return <Tooltip title={props.time}><span>{timeSince(new Date().getTime(), parsedTime)} ago</span></Tooltip>
}

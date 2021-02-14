import * as React from 'react';
import {LinearProgress} from "@material-ui/core";
import LineGraphWidget from "./graph/line/linegraph";

export interface WidgetProps {
    jobID: string
    type: string
}

export interface WidgetState {
    loading: boolean
    data: any
}

export default class Widget extends React.Component<WidgetProps, WidgetState> {
    constructor(props: WidgetProps) {
        super(props);
        this.state = {
            loading: true,
            data: null,
        }
    }

    componentDidMount = () => {

    }

    componentDidUpdate = (prevProps: WidgetProps, prevState: WidgetState) => {

    }

    render = () => {
        if (this.state.loading) {
            return <LinearProgress />
        }
        switch (this.props.type) {
            case "linegraph":
                return <LineGraphWidget data={this.state.data} />
            default:
                throw "unsupported widget type: " + this.props.type
        }
    }
}
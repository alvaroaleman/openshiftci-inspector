import * as React from "react";
import {Card, CardContent, CardHeader} from "@material-ui/core";
import {ResponsiveLine, Serie} from "@nivo/line"

interface IJobMetricsCardProps {
    id: string
}

interface IJobMetricsCardState {
    data: [Serie]
}

export default class JobPreviousCard extends React.Component<IJobMetricsCardProps, IJobMetricsCardState> {
    render = () => {
        return <React.Fragment>
            <Card>
                <CardHeader title={"Metrics"} />
                <CardContent>
                    <ResponsiveLine
                        data={this.state.data}
                        />
                </CardContent>
            </Card>
        </React.Fragment>
    }
}
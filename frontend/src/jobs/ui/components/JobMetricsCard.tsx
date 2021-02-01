import * as React from "react";
import {Box, Button, Card, CardContent, CardHeader, LinearProgress, TextField} from "@material-ui/core";
import {ResponsiveLine, Serie} from "@nivo/line"
import PlayArrowIcon from '@material-ui/icons/PlayArrow';
import HourglassEmptyIcon from '@material-ui/icons/HourglassEmpty';
import JobsMetricsService from "../../metrics";
import {QueryPoint, QuerySample, QuerySeries} from "../../../api-client";
import moment from "moment-timezone";

interface IJobMetricsCardProps {
    id: string
    metricsService: JobsMetricsService
}

interface IJobMetricsCardState {
    query: string
    loading: boolean
    loaded: boolean
    line: Serie[]
    timeRangeLowest: number,
    timeRangeHighest: number,
}

export default class JobMetricsCard extends React.Component<IJobMetricsCardProps, IJobMetricsCardState> {
    constructor(props: IJobMetricsCardProps) {
        super(props);
        this.state = {
            query: "histogram_quantile(0.99, rate(etcd_disk_wal_fsync_duration_seconds_bucket[5m]))",
            loading: false,
            loaded: false,
            line: [],
            timeRangeLowest: 0,
            timeRangeHighest: 0
        }
    }
    onQueryChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        this.setState({
            query: e.target.value
        })
    }

    onRun = async () => {
        this.setState({
            loading: true,
        })
        try {
            const response = await this.props.metricsService.getMetrics(this.props.id, this.state.query)
            if (response.matrix) {
                this.updateMatrix(response.matrix)
            } else if (response.vector) {
                this.updateVector(response.vector)
            } else if (response.scalar) {
                this.updateScalar(response.scalar)
            }
        } catch (e) {
            this.setState({
                loading: false,
            })
            throw e
        }
    }

    updateMatrix = (matrix: Array<QuerySeries>) => {
        let lowestTimestamp = 0
        let highestTimestamp = 0
        const series = matrix.map(function (value):Serie {
            return {
                id: value.labels.map(function (label) {
                    return label.name + "=" + label.value
                }).reduce(function (prev, current) {
                    return prev + "," + current
                }),
                data: value.points.map(function (v, index) {
                    if (lowestTimestamp === 0 || v.timestamp < lowestTimestamp) {
                        lowestTimestamp = v.timestamp
                    }
                    if (v.timestamp > highestTimestamp) {
                        highestTimestamp = v.timestamp
                    }
                    let d = moment.unix(Math.floor(v.timestamp/1000)).tz("UTC")
                    return {
                        x: d.format("HH:mm:ss"),
                        y: v.value
                    }
                })
            }
        })
        this.setState({
            loaded: true,
            loading: false,
            line: series,
            timeRangeLowest: lowestTimestamp,
            timeRangeHighest: highestTimestamp,
        })
    }

    updateVector = (vector: Array<QuerySample>) => {
        this.setState({
            loaded: true,
            loading: false,
        })
    }

    updateScalar = (scalar: QueryPoint) => {
        this.setState({
            loaded: true,
            loading: false,
        })
    }

    getTickValues = () => {
        const diff = (this.state.timeRangeHighest - this.state.timeRangeLowest)
        let result = []
        if (diff !== 0) {
            for (let i = this.state.timeRangeLowest; i < this.state.timeRangeHighest; i = i + diff / 10) {
                result.push(new Date(i).getUTCDate())
            }
        }
        return result
    }

    render = () => {
        return <React.Fragment>
            <Card>
                <CardHeader title={"Metrics"} />
                <CardContent>
                    <Box display={"flex"} flexDirection={"row"}>
                        <Box flex={1}>
                            <TextField
                                label="Query"
                                disabled={this.state.loading}
                                onChange={this.onQueryChange}
                                value={this.state.query}
                                fullWidth={true}
                                autoComplete={"on"}
                            />
                        </Box>
                        <Box>
                            <Button
                                variant="contained"
                                color="primary"
                                disabled={this.state.loading || !this.state.query}
                                onClick={this.onRun}>
                                Run {this.state.loading?<HourglassEmptyIcon />:<PlayArrowIcon />}
                            </Button>
                        </Box>
                    </Box>
                    <Box height={"400px"}>
                        {this.state.loading?<LinearProgress />:
                            this.state.loaded?
                            <ResponsiveLine
                                xScale={{
                                    type: 'time',
                                    format: "%H:%M:%S",
                                    precision: "second",
                                }}
                                xFormat="time:%H:%M:%S"
                                yScale={{ type: 'linear', min: 'auto', max: 'auto' }}
                                margin={{ top: 50, right: 110, bottom: 50, left: 60 }}
                                data={this.state.line}
                                pointSize={1}
                                pointColor={{ from: 'color', modifiers: [] }}
                                pointBorderWidth={1}
                                pointBorderColor={{ from: 'serieColor', modifiers: [] }}
                                axisTop={null}
                                axisRight={null}
                                axisBottom={{
                                    orient: 'bottom',
                                    tickSize: 5,
                                    tickPadding: 5,
                                    tickRotation: -60,
                                    tickValues: "every 1 minutes",
                                    format: '%H:%M:%S',
                                    legend: 'Time',
                                    legendOffset: 36,
                                    legendPosition: 'middle'
                                }}
                                axisLeft={{
                                    orient: 'left',
                                    tickSize: 5,
                                    tickPadding: 5,
                                    tickRotation: 0,
                                }}
                                useMesh={true}
                            />
                            :null
                        }
                    </Box>
                </CardContent>
            </Card>
        </React.Fragment>
    }
}

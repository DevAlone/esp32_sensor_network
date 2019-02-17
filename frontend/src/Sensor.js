import React, {Component} from 'react';
import Block from "./Block";
import Graph from "./Graph";
import withStyles from "@material-ui/core/es/styles/withStyles";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableRow from "@material-ui/core/TableRow";
import TableCell from "@material-ui/core/TableCell";
import WS from "./websocket";
import "./Sensor.css";
import moment from "moment";

const CustomTableCell = withStyles(theme => ({
    head: {
        backgroundColor: theme.palette.common.black,
        color: theme.palette.common.white,
    },
    body: {
        fontSize: 14,
    },
}))(TableCell);

class Sensor extends Component {
    constructor(params) {
        super(params);

        this.state = {
            "sensor": this.props.sensor,
            "sensorData": null,
        };
    }

    componentDidMount() {
        WS.SubscribeOnModelAdded("sensor_data", model => {
            if (model.sensor_id !== this.state.sensor.id) {
                return;
            }
            this.setState({
                sensorData: model,
            });
        });
    }

    render() {
        const classes = this.props;

        const sensor = this.state.sensor;
        const tableRows = [
            // ["sensor's ID", sensor.id],
            ["sensor's pin", sensor.pin],
            // ["sensor's type", sensor.type],
            // ["sensor's node mac address", this.state.sensor.sensor_node_mac_address],
        ];
        if (this.state.sensorData != null) {
            tableRows.push([
                sensor.type,
                this.state.sensorData.value
            ]);
            tableRows.push([
              "время последнего обновления",
                moment(new Date(this.state.sensorData.timestamp)).format("H:mm:ss"),
            ]);
        }
        console.log(sensor);

        return (<Block>
            <Table className={classes.table}>
                <TableBody>
                    {
                        tableRows.map(row => (
                            <TableRow className={classes.row} key={row[0]}>
                                <CustomTableCell component={"th"} scope={"row"}>{row[0]}</CustomTableCell>
                                <CustomTableCell align={"right"}>{row[1]}</CustomTableCell>
                            </TableRow>
                        ))
                    }
                </TableBody>
            </Table>
            <Graph
                modelName={"sensor_data"}
                itemId={this.state.sensor.id}
                itemIdFieldName={"sensor_id"}
                timestampMultiplier={1}
            />
        </Block>);
    }
}

export default withStyles()(Sensor);

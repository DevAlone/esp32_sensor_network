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
import {useTranslation, withTranslation, Trans} from 'react-i18next';

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
        const {t, i18n} = this.props;

        const sensor = this.state.sensor;
        const tableRows = [
            [t("common:sensor.sensor_pin"), sensor.pin],
            // ["sensor's type", sensor.type],
        ];
        if (this.state.sensorData != null) {
            let translatedSensorType = sensor.type;
            tableRows.push([
                t("common:sensor.sensor_type." + sensor.type),
                this.state.sensorData.value
            ]);
            tableRows.push([
              t("common:sensor.last_update_time"),
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

export default withTranslation()(Sensor);

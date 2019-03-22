import React, {Component} from 'react';
import withStyles from "@material-ui/core/styles/withStyles";
import 'amcharts3';
import 'amcharts3/amcharts/serial';
import AmCharts from '@amcharts/amcharts3-react';
import ToggleButton from '@material-ui/lab/ToggleButton';
import ToggleButtonGroup from '@material-ui/lab/ToggleButtonGroup';
import DoRequest from "./api";
import WS from "./websocket";
import Button from "@material-ui/core/Button";
import {useTranslation, withTranslation, Trans} from 'react-i18next';

const styles = theme => ({
    toggleButtonGroup: {
        background: "none",
        border: "none",
        boxShadow: "none",
        width: "100%",
        display: "flex",
        flexDirection: "row",
        flexWrap: "wrap",
        justifyContent: "space-between",
        alignItems: "center",
        alignContent: "center",
    },
    toggleButton: {
        fontSize: "11px",
        padding: "5px 10px",
        width: "auto",
        minWidth: "0",
    }
});

class Graph extends Component {
    constructor(props) {
        super(props);
        this.modelName = this.props.modelName;
        this.itemId = this.props.itemId;
        this.state = {};
        this.isLogarithmic = false;
        this.graphType = "line";  // "column";  // for bar charts
        this.xIsTimestamp = typeof this.props.xIsTimestamp !== "undefined" ? this.props.xIsTimestamp : true;
        this.minimumNumberOfResults = 0;
        // in ms
        this.timestampFilterButtons = [
            {
                "type": "last10Minutes",
                "pretty": "common:graph.last10Minutes",
                "filterFrom": 10 * 60 * 1000,
            },
            {
                "type": "lastHour",
                "pretty": "common:graph.lastHour",
                "filterFrom": 3600 * 1000,
            },
            {
                "type": "lastDay",
                "pretty": "common:graph.lastDay",
                "filterFrom": 24 * 3600 * 1000,
            },
            {
                "type": "lastWeek",
                "pretty": "common:graph.lastWeek",
                "filterFrom": 7 * 24 * 3600 * 1000,
            },
            {
                "type": "lastMonth",
                "pretty": "common:graph.lastMonth",
                "filterFrom": 30 * 24 * 3600 * 1000,
            },
            {
                "type": "lastYear",
                "pretty": "common:graph.lastYear",
                "filterFrom": 12 * 30 * 24 * 3600 * 1000,
            },
        ];
        this.itemIdFieldName = typeof this.props.itemIdFieldName === "undefined" ?
            "item_id" :
            this.props.itemIdFieldName;
        this.timestampMultiplier = typeof this.props.timestampMultiplier === "undefined" ?
            1000 :
            this.props.timestampMultiplier;
    }

    onWebsocketAdded = (message) => {
        message = message.data;
        if (message.model_name === this.modelName) {
            const data = message.data;
            if (data.sensor_id !== this.itemId) {
                return;
            }

            this.setState(prevState => {
                let newData = prevState.data.slice();
                newData.push({
                    timestamp: new Date(data.timestamp),
                    value: data.value,
                });

                return {
                    data: newData,
                };
            });
        }
    };

    toggleTimestampFilterButton = (buttonType) => {
        this.setState({
            data: null,
        });
        this.timestampFilterCurrentButton = buttonType;
        const filterButton = this.timestampFilterButtons.find(value => {
            return value.type === buttonType;
        });
        if (typeof filterButton === "undefined") {
            console.error("unknown type \"" + buttonType + "\"");
            return;
        }

        this.setTimestampFromFilter(filterButton.filterFrom);
    };

    setTimestampFromFilter = (from) => {
        if (from === 0) {
            this.filterTimestampFrom = 0;
        } else {
            this.filterTimestampFrom = Math.round((+new Date()) / this.timestampMultiplier) - from;
        }
        this.updateData();
    };

    loadData = (offset, limit, dataAccumulator, callback) => {
        let filter = this.itemId != null ?
            this.itemIdFieldName + " == " + this.itemId + "u"
            : "";

        if (this.filterTimestampFrom > 0) {
            if (filter.length > 0) {
                filter += " && ";
            }
            filter += "timestamp > " + (this.filterTimestampFrom - 1);
        }

        DoRequest("list_model", {
            "name": this.modelName,
            "order_by_fields": "timestamp",
            "filter": filter,
            "offset": offset,
            "limit": limit,
        }).then(response => {
            let data = [];
            for (let i in response.data.results) {
                var item = response.data.results[i];
                data.push({
                    timestamp: new Date(item.timestamp * this.timestampMultiplier),
                    value: item.value,
                });
            }
            if (data.length === 0) {
                callback(dataAccumulator);
            } else {
                dataAccumulator = dataAccumulator.concat(data);
                this.loadData(offset + limit, limit, dataAccumulator, callback);
            }
        });
    };

    updateData = () => {
        let data = [];
        this.loadData(0, 512, data, (resultData) => {
            if (resultData.length < this.minimumNumberOfResults) {
                const filterButtonIndex = this.timestampFilterButtons.findIndex(value => {
                    return value.type === this.timestampFilterCurrentButton;
                });
                if (typeof filterButtonIndex !== "undefined" && filterButtonIndex >= 0 && filterButtonIndex < this.timestampFilterButtons.length - 1) {
                    this.timestampFilterButtons.splice(filterButtonIndex, 1);
                    this.toggleTimestampFilterButton(this.timestampFilterButtons[filterButtonIndex].type);
                    return;
                }
            }

            this.setState({
                data: resultData,
            });
        });
    };

    componentDidMount() {
        this.timestampFilterCurrentButton = this.props.defaultTimestampFilter || this.timestampFilterButtons[0].type;
        this.toggleTimestampFilterButton(this.timestampFilterCurrentButton)
    }

    render() {
        const config = {
            "type": "serial",
            "theme": "light",
            "marginRight": 20,
            "marginLeft": 20,
            "autoMarginOffset": 20,
            "mouseWheelZoomEnabled": false,
            "valueAxes": [{
                "logarithmic": this.isLogarithmic,
                "id": "v1",
                "axisAlpha": 0.2,
                "position": "left",
                "ignoreAxisWidth": true
            }],
            "balloon": {
                "borderThickness": 1,
                "shadowAlpha": 0
            },
            "graphs": [{
                "id": "g1",
                "lineColor": "#77c0e2",
                "fillAlphas": 0.2,
                "balloon": {
                    "drop": true,
                    "adjustBorderColor": false,
                    "color": "#ffffff"
                },
                "bullet": "round",
                "bulletBorderAlpha": 1,
                "bulletColor": "#FFFFFF",
                "bulletSize": 5,
                "hideBulletsCount": 50,
                "lineThickness": 2,
                "title": "", // TODO: add title
                "type": this.graphType,
                "useLineColorForBulletBorder": true,
                "valueField": "value",
                "balloonText": "<span style='font-size:18px;'>[[value]]</span>"
            }],
            /*"chartScrollbar": {
                "graph": "g1",
                "scrollbarHeight": 80,
                "autoGridCount": true,
            },*/
            "chartCursor": {
                "limitToGraph": "g1",
            },
            "categoryField": "timestamp",
            "categoryAxis": {
                "minPeriod": "ss",
                "parseDates": this.xIsTimestamp,
                "axisColor": "#DADADA",
                "dashLength": 1,
                "minorGridEnabled": true
            },
            "dataProvider": this.state.data,
            "export": {"enabled": false},
            "legend": {"enabled": false},
            "zoomControl": {"zoomControlEnabled": false},
        };

        const {classes, t} = this.props;

        return (
            <div style={{width: "100%"}}>
                {
                    <ToggleButtonGroup
                        className={classes.toggleButtonGroup}
                        exclusive
                        value={this.timestampFilterCurrentButton}>
                        {
                            this.timestampFilterButtons.map(button => {
                                return (
                                    <ToggleButton
                                        key={button.type}
                                        component={Button}
                                        value={button.type}
                                        onClick={() => this.toggleTimestampFilterButton(button.type)}
                                        disabled={typeof this.state.data === "undefined" || this.state.data === null}
                                        className={classes.toggleButton}
                                    >
                                        {t(button.pretty)}
                                    </ToggleButton>
                                );
                            })
                        }
                    </ToggleButtonGroup>
                }
                {
                    typeof (this.state.data) === "undefined" || this.state.data === null ?
                        <h4>{t("common:loading") + "..."}</h4>
                        : this.state.data.length === 0 ?
                        <h4>{t("common:graph:nothing_is_here")}</h4>
                        : null
                }
                {
                    <div style={{
                        width: "100%",
                        height: "500px",
                    }}>
                        <AmCharts.React style={{
                            width: "100%",
                            height: typeof (this.state.data) !== "undefined" && this.state.data !== null && this.state.data.length > 0 ? "500px" : "0",
                            // overflow: "hidden",
                        }} options={config}/>
                    </div>
                }
            </div>
        );
    }
}

export default withTranslation()(withStyles(styles)(Graph));

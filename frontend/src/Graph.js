import React, {Component} from 'react';
import withStyles from "@material-ui/core/styles/withStyles";
import 'amcharts3';
import 'amcharts3/amcharts/serial';
import AmCharts from '@amcharts/amcharts3-react';
import DoRequest from "./api";
import WS from "./websocket";

const styles = theme => ({});

class Graph extends Component {
    constructor(props) {
        super(props);
        this.modelName = this.props.modelName;
        this.itemId = this.props.itemId;
        this.itemIdFieldName = typeof this.props.itemIdFieldName === "undefined" ?
            "item_id" :
            this.props.itemIdFieldName;
        this.state = {};
        this.isLogarithmic = false;
        this.graphType = "line";  // "column";  // for bar charts
        this.xIsTimestamp = true;
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

    loadData = (offset, limit, dataAccumulator, callback) => {
        let filter = this.itemId != null ?
            this.itemIdFieldName + " == " + this.itemId + "u"
            : "";

        if (filter.length > 0) {
            filter += " && ";
        }
        filter += "timestamp > " + (Math.round(new Date()) - 1000 * 3600 * 2);

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

    componentDidMount() {
        let data = [];
        this.loadData(0, 512, data, (resultData) => {
            this.setState({
                data: resultData,
            });
            WS.SubscribeOnMessage("model_added", this.onWebsocketAdded);
        });
    }

    componentWillUnmount() {
        // clearInterval(this.state.timer);
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

        return (
            <div style={{width: "100%"}}>
                {
                    typeof (this.state.data) === "undefined" ?
                        <h4>Загрузка...</h4>
                        : this.state.data.length === 0 ?
                        <h4>Ничего нет :(</h4>
                        : null
                }
                {
                    <AmCharts.React style={{
                        width: "100%",
                        height: typeof (this.state.data) !== "undefined" && this.state.data.length > 0 ? "250px" : "0px",
                        overflow: "hidden",
                    }} options={config}/>
                }
            </div>
        );
    }
}

export default withStyles(styles)(Graph);

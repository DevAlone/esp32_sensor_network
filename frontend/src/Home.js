import React, {Component} from 'react';
import Block from "./Block";
import NiceLink from "./NiceLink";
import Graph from "./Graph";
import DoRequest from "./api";
import ExpansionPanel from "@material-ui/core/ExpansionPanel";
import ExpansionPanelSummary from "@material-ui/core/ExpansionPanelSummary";
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import Typography from "@material-ui/core/Typography";
import ExpansionPanelDetails from "@material-ui/core/ExpansionPanelDetails";
import {useTranslation, withTranslation, Trans} from 'react-i18next';
import "./Home.css";
import Paper from "@material-ui/core/Paper";
import Sensor from "./Sensor";
import Logs from "./Logs";

class Home extends Component {
    constructor(params) {
        super(params);

        this.state = {
            nodes: null,
        }
    }

    componentDidMount() {
        this.loadData();
    }

    loadData = () => {
        DoRequest("list_model", {
            "name": "sensor_node",
            // "order_by_fields": "timestamp",
            // "filter": filter,
            "offset": 0,
            // "limit": 0,
        }).then(response => {
            if (typeof response === 'undefined') {
                return;
            }
            this.setState(prevState => {
                prevState.nodes = {};
                for (let i in response.data.results) {
                    const sensorNode = response.data.results[i];
                    const macAddress = sensorNode.mac_address;
                    prevState.nodes[macAddress] = sensorNode;

                    this.loadSensors(macAddress);
                }

                return prevState;
            });
        });
    };

    loadSensors = (macAddress) => {
        // fetch sensors
        DoRequest("list_model", {
            "name": "sensor",
            // "order_by_fields": "timestamp_ms",
            "filter": "sensor_node_mac_address == " + macAddress + "u",
        }).then(response => {
            if (typeof response === 'undefined') {
                return;
            }
            this.setState(prevState => {
                let sensors = {};
                for (let i in response.data.results) {
                    let sensor = response.data.results[i];
                    sensors[sensor.id] = sensor
                    // TODO: fetch data
                }
                prevState.nodes[macAddress].sensors = sensors;

                return prevState;
            });
        })
    };

    changeLanguage = languageCode => {
        this.props.i18n.changeLanguage(languageCode);
    };

    render() {
        const classes = this.props;
        const {t, i18n} = this.props;

        return (
            <div>
                <div class="languageButtons">
                    <button onClick={() => this.changeLanguage('en')}>English</button>
                    <button onClick={() => this.changeLanguage('ru')}>Русский</button>
                    <button onClick={() => this.changeLanguage('ua')}>Українська</button>
                </div>
                <Block><Typography>esp32_sensor_network {t("common:home.demo")}</Typography></Block>
                <Typography>{t("common:home.available_nodes")}:</Typography>
                {
                    this.state.nodes === null ?
                        t("common:loading") + "..." :
                        Object.keys(this.state.nodes).map((key) => {
                            const node = this.state.nodes[key];
                            return <div key={key}>
                                <Block>
                                    <Typography>
                                        {t("common:home.node_mac_is") + " " + node.mac_address}
                                        {
                                            node.name.length > 0 ? <Typography>
                                                {t("common:home.node_name_is") + " " + node.name}

                                            </Typography> : null
                                        }
                                    </Typography>
                                </Block>

                                <div>
                                    <Typography>{t("common:home.available_sensors")}:</Typography>
                                </div>
                                <div className={"sensorsContainer"}>
                                    {
                                        typeof node.sensors === "undefined" ?
                                            t("common:loading") + "..." :
                                            Object.keys(node.sensors).map((key) => {
                                                const sensor = node.sensors[key];
                                                return <div className={"sensorContainer"} key={key}>
                                                    <Sensor sensor={sensor}/>
                                                </div>;
                                            })
                                    }
                                </div>
                            </div>;
                        })
                }
                <Logs />
            </div>
        );
    }
}

export default withTranslation()(Home);

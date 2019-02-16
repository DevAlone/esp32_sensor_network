import React, {Component} from 'react';
import Block from "./Block";
import NiceLink from "./NiceLink";
import Graph from "./Graph";
import DoRequest from "./api";

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

    render() {
        return (
            <div>
                {
                    this.state.nodes === null ?
                        "Loading..." :
                        Object.keys(this.state.nodes).map((key) => {
                            const node = this.state.nodes[key];
                            return <Block key={key}>
                                <p>Node's mac is {node.mac_address}</p>
                                <p>Node's name is {node.name.length > 0 ? '"' + node.name + '"' : "empty"}</p>
                                <p>Available sensors:</p>
                                {
                                    typeof node.sensors === "undefined" ?
                                        "loading..." :
                                        Object.keys(node.sensors).map((key) => {
                                            const sensor = node.sensors[key];
                                            return <div key={key}>
                                                <p>sensor's id is {sensor.id}</p>
                                                <p>sensor's pin is {sensor.pin}</p>
                                                <p>sensor's type is {sensor.type}</p>
                                                <p>sensor's node mac address is {sensor.sensor_node_mac_address}</p>
                                                <p>graph:</p>
                                                <Graph
                                                    modelName={"sensor_data"}
                                                    itemId={sensor.id}
                                                    itemIdFieldName={"sensor_id"}
                                                    timestampMultiplier={1}
                                                />
                                            </div>;
                                        })
                                }
                            </Block>;
                        })
                }
                <Block>Home</Block>
            </div>
        );
    }
}

export default Home;

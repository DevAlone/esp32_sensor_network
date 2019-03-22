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
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemText from "@material-ui/core/ListItemText";
import ExpansionPanel from "@material-ui/core/ExpansionPanel";
import ExpansionPanelSummary from "@material-ui/core/ExpansionPanelSummary";
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import ExpansionPanelDetails from "@material-ui/core/ExpansionPanelDetails";
import Typography from "@material-ui/core/Typography";
import "./Logs.css";

class Logs extends Component {
    constructor(params) {
        super(params);
        this.logsSize = 32;

        this.state = {
            "logs": [],
        };
    }

    componentDidMount() {
        WS.SubscribeOnModelAdded("log_request", model => {
            this.setState(prevState => {
                while (prevState.logs.length > this.logsSize) {
                    prevState.logs.shift();
                }
                prevState.logs.push(model);
                return prevState;
            });
        });
    }

    render() {
        const classes = this.props;
        const {t, i18n} = this.props;

        console.log(this.state.logs);

        return (
            <ExpansionPanel className={"expansionPanel"}>
                <ExpansionPanelSummary  expandIcon={<ExpandMoreIcon />}>
                    <Typography>{t("common:logs.logs")}:</Typography>
                </ExpansionPanelSummary>
                <ExpansionPanelDetails>
                    <List>
                        {
                            this.state.logs.map((item, i) => {
                                return <ListItem key={i}>
                                    <ListItemText primary={JSON.stringify(item)}/>
                                </ListItem>
                            })
                        }
                    </List>
                </ExpansionPanelDetails>
            </ExpansionPanel>
        );
    }
}

export default withTranslation()(withStyles()(Logs));

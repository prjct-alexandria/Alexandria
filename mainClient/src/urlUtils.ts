import configData from "./config.json";

export default function backEndUrl() {
    return configData.protocol + "://" + window.location.hostname + ":" + configData.back_end_port
}
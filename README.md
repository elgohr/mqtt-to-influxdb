# mqtt-to-influxdb

Bridge to write MQTT to InfluxDB.

## Usage

| Environment Variable | Description |
| --- | --- |
| MQTT_NAME | Name that the bridge will use on mqtt |
| MQTT_URL | MQTT Broker URL |
| MQTT_USERNAME | Username for authentication |
| MQTT_PASSWORD | Password for authentication |

| Environment Variable | InfluxDb < 2.0 | InfluxDb 2.0 |
| --- | --- | --- |
| INFLUX_URL | Url to the InfluxDB | Url to the InfluxDB |
| INFLUX_TOKEN | use "username:password" | Token for authentication |
| INFLUX_ORGANIZATION | leave empty | Organization to use |
| INFLUX_BUCKET |  use database name | Bucket to use |

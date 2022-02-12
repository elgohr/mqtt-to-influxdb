# mqtt-to-influxdb

Bridge to write MQTT to InfluxDB.

## Install

Please have a look at the binaries in https://github.com/elgohr/mqtt-to-influxdb/releases

## Usage

### Mandatory Configuration
| Environment Variable | Description                           |
|----------------------|---------------------------------------|
| MQTT_NAME            | Name that the bridge will use on mqtt |
| MQTT_URL             | MQTT Broker URL                       |
| MQTT_USERNAME        | Username for authentication           |
| MQTT_PASSWORD        | Password for authentication           |

| Environment Variable | InfluxDb < 2.0          | InfluxDb 2.0             |
|----------------------|-------------------------|--------------------------|
| INFLUX_URL           | Url to the InfluxDB     | Url to the InfluxDB      |
| INFLUX_TOKEN         | use "username:password" | Token for authentication |
| INFLUX_ORGANIZATION  | leave empty             | Organization to use      |
| INFLUX_BUCKET        | use database name       | Bucket to use            |

### Optional

| Environment Variable  | Description                                                             |
|-----------------------|-------------------------------------------------------------------------|
| INFLUX_RETRY_INTERVAL | Retry interval for writing to InfluxDB in Milliseconds (max 4294967295) |
| INFLUX_MAX_RETRIES    | Numbers of retries for writing to InfluxDB                              |
| INFLUX_BATCH_SIZE     | Number of points send in a single request to InfluxDB                   |

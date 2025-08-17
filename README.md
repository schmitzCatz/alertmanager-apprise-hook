# ![Logo](./logo.png "Logo") Apprise Alertmanager Webhook

Simple Go application to forward Prometheus Alertmanager Alerts to Apprise

## Installation

``docker pull ghcr.io/schmitzcatz/alertmanager-apprise-hook:latest ``

## Usage

``docker run -d -p 8080:8080 -e APPRISE_URL=http://apprise.example.com ghcr.io/schmitzcatz/alertmanager-apprise-hook:latest``

## Configuration

| Variable       | Default | Description                 | Mandatory |
|----------------|---------|-----------------------------|-----------|
| TAG            | all     | Apprise Tag/Group to notify | ❌         |
| APPRISE_URL    |         | URL to your Apprise server  | ✅         |
| LISTEN_ADDRESS | :8080   | Hook listen Address         | ❌         |

## Contribution

Consider contributing? Check [Contributing](CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md) beforehand.

## License

[GNU GPLv3](https://spdx.org/licenses/GPL-3.0-or-later.html)


# Windows Installation

## Installing

To install the collector on Windows run the Powershell command below to install the MSI with no UI.
```pwsh
msiexec /i "https://github.com/observIQ/observiq-otel-collector/releases/latest/download/observiq-otel-collector.msi" /quiet
```

Alternately, for an interactive installation [download the latest MSI](https://github.com/observIQ/observiq-otel-collector/releases/latest).

After downloading the MSI, simply double click it to open the installation wizard. Follow the instructions to configure and install the collector.

## Configuring the Collector

After installing, the `observiq-otel-collector` service will be running and ready for configuration! 

The collector logs to `C:\Program Files\observIQ OpenTelemetry Collector\log\collector.log` by default.

By default, the config file for the collector can be found at `C:\Program Files\observIQ OpenTelemetry Collector\config.yaml`. When changing the configuration,the collector service must be restarted in order for config changes to take effect.

For more information on configuring the collector, see the [OpenTelemetry docs](https://opentelemetry.io/docs/collector/configuration/).

## Restarting the Collector
Restarting the collector may be done through the services dialog.
To access the services dialog, press Win + R, enter `services.msc` into the Run dialog, and press enter.

![The run dialog](./screenshots/windows/launch-services.png)

Locate the "observIQ Distro for OpenTelemetry Collector" service, right click the entry, and click "Restart" to restart the collector.

![The services dialog](./screenshots/windows/stop-restart-service.png)

Alternatively, the Powershell command below may be run to restart the collector service.
```pwsh
Restart-Service -Name "observiq-otel-collector"
```

## Stopping the Collector

Stopping the collector may be done through the services dialog.
To access the services dialog, press Win + R, enter `services.msc` into the Run dialog, and press enter.

![The run dialog](./screenshots/windows/launch-services.png)

Locate the "observIQ Distro for OpenTelemetry Collector" service, right click the entry, and click "Stop" to stop the collector.

![The services dialog](./screenshots/windows/stop-restart-service.png)

Alternatively, the Powershell command below may be run to stop the collector service.
```pwsh
Stop-Service -Name "observiq-otel-collector"
```

## Starting the Collector

Starting the collector may be done through the services dialog.
To access the services dialog, press Win + R, enter `services.msc` into the Run dialog, and press enter.

![The run dialog](./screenshots/windows/launch-services.png)

Locate the "observIQ Distro for OpenTelemetry Collector" service, right click the entry, and click "Start" to start the collector.

![The services dialog](./screenshots/windows/start-service.png)

Alternatively, the Powershell command below may be run to start the collector service.
```pwsh
Start-Service -Name "observiq-otel-collector"
```

## Uninstalling

To uninstall the collector on Windows, navigate to the control panel, then to the "Uninstall a program" dialog.

![The control panel](./screenshots/windows/control-panel-uninstall.png)

Locate the `"observIQ Distro for OpenTelemetry Collector"` entry, and select uninstall. 

![The uninstall or change a program dialog](./screenshots/windows/uninstall-collector.png)

Follow the wizard to complete removal of the collector.

Alternatively, Powershell command below may be run to uninstall the collector.
```pwsh
(Get-WmiObject -Class Win32_Product -Filter "Name = 'observIQ Distro for OpenTelemetry Collector'").Uninstall()
```

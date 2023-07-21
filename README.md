# GoJira 🦖

![Silly dinosaur with sunglasses](https://preview.redd.it/5j0lfvzx8c861.jpg?width=640&crop=smart&auto=webp&s=6a0c5951ed2da6166d7a8a1d53269d7d3d8ad039)

GoJira is a fully open source data distribution tool that connects JIRA webhooks to external services like GitHub, Discord, and Slack.

## Quick Start ⚡

You can either download the binary relative to your operating system and fire that baby up, clone the repository and run it natively, or use docker to compose the app and it's extra goodies ( Prom, Jaeger, Redis, etc... )

### Install binary

Download, install, and run the binary for the latest release that matches your Operating System.

### Run the project

Clone/fork this project and run it manually using the makefile

```bash
make run
```

### Run with docker

Make sure you have your local config file set up with your environment variables first. Then either use the makefile or a custom docker command to run the project.

Makefile
```bash
make develop
```

Docker
```bash
docker-compose -f docker-compose.local.yml up --build
```

### Enable Discord Bot

The discord bot is run with [discordgo](https://github.com/bwmarrin/discordgo). To enable the bot to send messages in discord you will need a few things to get started. The first is setting up a bot in discord developer console which can be found [here](https://discord.com/developers). Create your bot, give it the proper permissions to send messages/embeds, copy the Bot Token, then invite it to your server(s). The bot token should be placed in the config under `discord: BotToken`. Enable the discord bot in the config then start/restart this application. The bot can be set up to filter all incoming webhook messages to one channel, or to split incoming messages into unique channels. 

## How it works 🔥

JIRA is the staple of this project, and this was built for personal use for digesting my own JIRA data and effectively distributing it to systems that I currently use.

See JIRA WebHooks documentation for reference to payloads [here](https://developer.atlassian.com/server/jira/platform/webhooks/)

### Receiving

Initially JIRA will need to be set up to send data to this service. The endpoint that accepts the JIRA data is at `/v1/receiver/accept`. The incoming payload will be converted to a matching struct and the event types will be interpreted.

### Interpreting

JIRA has a handful of event action types that it can send. It's the interpreters job to parse these actions and create a simple yet concrete decision on what to do with the parsed actions. Once the event is interpreted and an `EventMap` is created its on to distribution.

### Distribution

The distributor has arguably the biggest lift of them all. With payload body and event map in hand it's job is to send the data to where it needs to go. Based on configurations and mappings it will decide what, where, and how to send the data. As a quick example of a simple JIRA->GitHub communication the distributor will determine: Repository->Issue|Project->IssueType|ProjectTaskType->Assignee.

## Telemetry 🌎

OpenTelemtry is included in this project for distributed tracing. It's not absolutely necessary but in a world of integrations this is lightweight and can help seamlessly connect and trace this project when connected to others.

Example usage:
```go
ctx, span := otel.Tracer("Receiver").Start(utils.GetRequestCtx(c), "receiverHandlers.Accept")

p := &models.JiraWebhookMessage{}
if err := c.Bind(p); err != nil {
    utils.LogResponseError(c, h.logger, err)
    span.End()
    return c.JSON(httpErrors.ErrorResponse(err))
}

if err := h.receiverService.Accept(ctx, p); err != nil {
    utils.LogResponseError(c, h.logger, err)
    span.End()
    return c.JSON(httpErrors.ErrorResponse(err))
}
```

The telemetry library used is `otel` for Go. You can read more about OpenTelemetry [here](https://opentelemetry.io/docs/what-is-opentelemetry/).

## Observability 🔎

### Pprof
By default this project comes with the ability to spin up a separate pprof endpoint at port `5555` for debugging and taking snapshots for application insights.

### Prometheus
By default this project comes packed with middleware that exposes a `/metrics` endpoint and returns prometheus typed data. Additionally the docker compose file will also spin up a Prometheus container bound to ports `9101:9100`.

### Grafana
Not absolutely necessary but the docker compose will also spin up a Grafana container accessible at port `3000:3000` for use with exporting the prometheus/app metrics to create visible charts/graphs.

### Jaeger
Originally I had included Jaeger in this build, and it still does have some remnants but with all the other observability baked in I didn't see it as absolutely necessary. If you plan on configuring this as a microservice or want to create a distributed system with this project build in Jaeger can quickly be added to the core service spin-up.



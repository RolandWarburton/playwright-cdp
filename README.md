# Playwright CDP Experiment

Test exercise to see how to set up a client to control playwright via chrome CDP (chrome devtools protocol) and a web server to accept commands from.

## Situation

In restrictive environments it may not be possible to run your own code,
in situations where you have access to a browser you can control a remote playwright instance (agent) via an API (server).

![Diagram showing different parts of the test](./diagrams/dia.svg)

![Program flow of the test](./diagrams/flow.svg)

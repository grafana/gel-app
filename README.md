# GEL POC

As of Sept 9th 2019 - the backend GEL code in `pkg` is no longer runnable. This is because the code was moved from the grafana-enterprise repo. We aim to have it runnable again as a plugin by Sept 13th.

The `pkg/gel` folder in particular was meant for enterprise so needs to be reworked. The `pkg/data`
and `pkg/mathexp/...` packages can run tests.

## Instructions

Last Updated: 2019-09-11

Currently Grafana views GEL as a backend plugin. When you have built it, Grafana should run the service.

In this repo:

### Build Frontend

1. `yarn install`
2. `yarn dev`. If you run into an issue, try deleting the `yarn.lock` file and running `yarn install` again. If that does not work, ask Ryan for help in the #gel channel.

### Build Go Backend

1. `make vendor`
2. `make build` (on Linux), (or `make build-darwin` for mac) - which should output an executable to `dist/datasource`.

**Warning**: The `go.mod` file has `replace` overrides to reference specific *commits* of the bidirectional_poc branches for Grafana. So I believe if those branches get new commits, the `replace` directives in `go.mod` will need to be updated to get the changes.

### Running

You can then run grafana with the [bidirectional_poc](https://github.com/grafana/grafana/tree/bidirectional_poc) branch. If you enable debug logging in Grafana's ini you should see the plugin executed (or errors/problems). You should also see a gel process running after grafana starts `ps aux | grep gel | grep -v grep` or task manager or something :-).

## Reference Branches

- [Grafana bidirectional_poc](https://github.com/grafana/grafana/tree/bidirectional_poc)
- [Grafana Plugin Model bidirectional_poc](https://github.com/grafana/grafana-plugin-model/tree/bidirectional_poc)

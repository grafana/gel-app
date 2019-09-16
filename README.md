# GEL POC

As of Sept 16, 2019 this will run as a datasource plugin.

## Instructions

Last Updated: 2019-09-16

Currently Grafana will run GEL as a backend plugin. When you have built it, Grafana should run the service with the [bidirectional_poc](https://github.com/grafana/grafana/tree/bidirectional_poc) of Grafana.

In this repo:

### Build Frontend

If you run into an issue with the following two steps, try deleting the `yarn.lock` file and running `yarn install` again. If that does not work, ask Ryan for help in the `#gel` slack channel.

1. `yarn install`
2. `yarn dev`.

### Build Go Backend

1. `make vendor`
2. `make build` (on Linux), (or `make build-darwin` for mac) - which should output an executable to `dist/datasource`.

**Warning**: The `go.mod` file has `replace` overrides to reference specific *commits* of the bidirectional_poc branches for Grafana. So I believe if those branches get new commits, the `replace` directives in `go.mod` will need to be updated to get the changes.

### Running

You will need to run grafana with the [bidirectional_poc](https://github.com/grafana/grafana/tree/bidirectional_poc) branch. Be sure to run `GO111MODULE=on go mod vendor` on the branch to get the updated plugin-model.

The plugin will need to be added to grafana's `data/plugins` through your preferred method. One way is to create a symlink, in the `data/plugins` directory something like `ln -s /home/kbrandt/src/github.com/grafana/gel-app/ gel-app`.

When grafana is running, in Grafana's UI:

1. Enable the GEL plugin in Grafana's UI.
2. Add the GEL datasource, for datasource settings set the URL to `/api/tsdb/query`.

If you enable debug logging in Grafana's ini you should see the plugin executed (or errors/problems). You should also see a gel process running after grafana starts `ps aux | grep gel | grep -v grep` or task manager or something :-).

### Using GEL

Select "GEL Data", the GEL datasource as the Query.

Currently it is expected you *only do one top level (upper right) Add "Query"* which would create refId "A".

From there, you can click "Add Query" and/or "Add GEL" within the GEL query that was added. Each of these will be refs like "GA", "GB", etc. You can reference queries with `$GA` for example in the math text text field, or in the Fields of reduction.

#### Caveats

The first query within the GEL datasource must be a GEL Node ("Add GEL" in the UI). That is to say it has to be a math expr or a reduction function. (This is because the grafana `/api/tsdb/query` endpoint looks at the first query to pick the datasource).

## Reference Branches

- [Grafana bidirectional_poc](https://github.com/grafana/grafana/tree/bidirectional_poc)
- [Grafana Plugin Model bidirectional_poc](https://github.com/grafana/grafana-plugin-model/tree/bidirectional_poc)

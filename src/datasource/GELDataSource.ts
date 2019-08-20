// Types
import { DataQueryRequest, DataQueryResponse, DataSourceApi, DataSourceInstanceSettings } from '@grafana/ui';
import { getBackendSrv, config } from '@grafana/runtime';

import { GELQuery, GELDataSourceOptions } from './types';
import { responseToDataFrame } from './util';

export class GELDataSource extends DataSourceApi<GELQuery, GELDataSourceOptions> {
  constructor(private instanceSettings: DataSourceInstanceSettings<GELDataSourceOptions>) {
    super(instanceSettings);
  }

  /**
   * Convert a query to a simple text string
   */
  getQueryDisplayText(query: GELQuery): string {
    return 'GEL: ' + query;
  }

  async query(options: DataQueryRequest<GELQuery>): Promise<DataQueryResponse> {
    const { url } = this.instanceSettings;
    const { targets, startTime, ...opts } = options;
    if (targets.length > 1) {
      return Promise.reject('Only query supported right now');
    }
    if (targets.length < 1) {
      return Promise.resolve({ data: [] });
    }
    const first: GELQuery = targets[0];
    const target = {
      ...first,
      queries: first.queries.map(q => {
        const ds = config.datasources[q.datasource || config.defaultDatasource];
        console.log('DS', ds);
        return {
          ...q,
          datasourceId: ds.id,
        };
      }),
    };

    return getBackendSrv()
      .post(url!, {
        options: opts,
        gel: target,
      })
      .then(res => {
        console.log('RESPONSE', res);
        return { data: [responseToDataFrame(res)] };
      })
      .catch(err => {
        err.isHandled = true;
        console.error('Error', err);
        return { data: [] };
      });
  }

  async testDatasource() {
    return Promise.resolve({
      status: 'success',
      message: 'TODO, actually check?',
    });
  }
}

export default GELDataSource;

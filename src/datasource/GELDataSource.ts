// Types
import { DataQueryRequest, DataQueryResponse, DataSourceApi, DataSourceInstanceSettings } from '@grafana/ui';
import { getBackendSrv, config } from '@grafana/runtime';

import { GELDataSourceOptions, TempGELQueryWrapper, GEL_DS_KEY } from './types';
import { responseToDataFrame } from './util';

export class GELDataSource extends DataSourceApi<TempGELQueryWrapper, GELDataSourceOptions> {
  constructor(private instanceSettings: DataSourceInstanceSettings<GELDataSourceOptions>) {
    super(instanceSettings);
  }

  /**
   * Convert a query to a simple text string
   */
  getQueryDisplayText(query: TempGELQueryWrapper): string {
    return 'GEL: ' + query;
  }

  async query(options: DataQueryRequest<TempGELQueryWrapper>): Promise<DataQueryResponse> {
    const { url } = this.instanceSettings;
    const { targets, startTime, ...opts } = options;
    if (targets.length > 1) {
      return Promise.reject('Only query supported right now');
    }
    if (targets.length < 1) {
      return Promise.resolve({ data: [] });
    }
    const first: TempGELQueryWrapper = targets[0];
    (opts as any).targets = first.queries.map(q => {
      if (q.datasource === GEL_DS_KEY) {
        return q;
      }
      const ds = config.datasources[q.datasource || config.defaultDatasource];
      return {
        ...q,
        datasourceId: ds.id,
      };
    });

    return getBackendSrv()
      .post(url!, {
        options: opts,
      })
      .then(res => {
        return { data: responseToDataFrame(res) };
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

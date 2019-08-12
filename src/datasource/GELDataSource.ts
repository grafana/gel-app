// Types
import { DataQueryRequest, DataQueryResponse, DataSourceApi, DataSourceInstanceSettings, MetricFindValue } from '@grafana/ui';
import { DataFrame } from '@grafana/data';

import { GELQuery, GELDataSourceOptions } from './types';

export class GELDataSource extends DataSourceApi<GELQuery, GELDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<GELDataSourceOptions>) {
    super(instanceSettings);
  }

  /**
   * Convert a query to a simple text string
   */
  getQueryDisplayText(query: GELQuery): string {
    return 'Plugin: ' + query;
  }

  metricFindQuery(query: string, options?: any): Promise<MetricFindValue[]> {
    return new Promise((resolve, reject) => {
      const names: MetricFindValue[] = [];
      // for (const series of this.data) {
      //   for (const field of series.fields) {
      //     // TODO, match query/options?
      //     names.push({
      //       text: field.name,
      //     });
      //   }
      // }
      resolve(names);
    });
  }

  async query(options: DataQueryRequest<GELQuery>): Promise<DataQueryResponse> {
    const results: DataFrame[] = [];
    for (const query of options.targets) {
      console.log('QUERY: ', query);
    }
    return { data: results };
  }

  async testDatasource() {
    return Promise.resolve({
      status: 'success',
      message: 'TODO, actually check?',
    });
  }
}

export default GELDataSource;

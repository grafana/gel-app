// Types
import { DataQueryRequest, DataQueryResponse, DataSourceApi, DataSourceInstanceSettings } from '@grafana/ui';
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

  async query(options: DataQueryRequest<GELQuery>): Promise<DataQueryResponse> {
    const results: DataFrame[] = [];
    for (const query of options.targets) {
      console.log('QUERY: ', query);
    }
    console.log('RETURN empty results: ', results);
    return Promise.resolve({ data: results });
  }

  async testDatasource() {
    return Promise.resolve({
      status: 'success',
      message: 'TODO, actually check?',
    });
  }
}

export default GELDataSource;

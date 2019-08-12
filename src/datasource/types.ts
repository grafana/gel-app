import { DataQuery, DataSourceJsonData } from '@grafana/ui';

export interface GELQuery extends DataQuery {
  queries: DataQuery[];
  expression: string;
}

export interface GELDataSourceOptions extends DataSourceJsonData {
  // maybe a beta flag?
}

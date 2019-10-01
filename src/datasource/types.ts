import { DataQuery, DataSourceJsonData } from '@grafana/ui';

export const GEL_DS_KEY = '-- GEL --';

/**
 * This is just a temporary wrapper around the query interface
 */
export interface TempGELQueryWrapper extends DataQuery {
  queries: DataQuery[];
}

export enum GELQueryType {
  math = 'math',
  reduce = 'reduce',
  resample = 'resample',
}

/**
 * For now this is a single object to cover all the types.... would likely
 * want to split this up by type as the complexity increases
 */
export interface GELQuery extends DataQuery {
  type: GELQueryType;
  reducer?: string;
  expression?: string;
  rule?: string;
  downsampler?: string;
}

export interface GELDataSourceOptions extends DataSourceJsonData {
  // maybe a beta flag?
}

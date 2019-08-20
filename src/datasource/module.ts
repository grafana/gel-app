import { DataSourcePlugin } from '@grafana/ui';

import { GELDataSource } from './GELDataSource';

import { GELConfigEditor } from './GELConfigEditor';
import { GELQueryEditor } from './GELQueryEditor';
import { TempGELQueryWrapper, GELDataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<GELDataSource, TempGELQueryWrapper, GELDataSourceOptions>(GELDataSource)
  .setConfigEditor(GELConfigEditor)
  .setQueryEditor(GELQueryEditor);

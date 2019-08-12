// Libraries
import React, { PureComponent } from 'react';

// Types
import { GELDataSourceOptions } from './types';

import { DataSourcePluginOptionsEditorProps, DataSourceSettings } from '@grafana/ui';

type PluginSettings = DataSourceSettings<GELDataSourceOptions>;

interface Props extends DataSourcePluginOptionsEditorProps<PluginSettings> {}

interface State {
  text: string;
}

export class GELConfigEditor extends PureComponent<Props, State> {
  state = {
    text: '',
  };

  componentDidMount() {}

  render() {
    return <div>GEL Datasource Config Editor</div>;
  }
}

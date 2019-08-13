// Libraries
import React, { PureComponent, ChangeEvent } from 'react';

// Types
import { GELDataSourceOptions } from './types';

import { DataSourcePluginOptionsEditorProps, DataSourceSettings, FormField } from '@grafana/ui';

type PluginSettings = DataSourceSettings<GELDataSourceOptions>;

interface Props extends DataSourcePluginOptionsEditorProps<PluginSettings> {}

interface State {}

export class GELConfigEditor extends PureComponent<Props, State> {
  state = {};

  componentDidMount() {}

  onURLChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      url: event.target.value,
      access: 'direct', // HARDCODE For now!
    });
  };

  render() {
    const { options } = this.props;

    return (
      <div className="gf-form-group">
        <div className="gf-form">
          <FormField label="URL" labelWidth={6} onChange={this.onURLChange} value={options.url} placeholder="GEL Endpoint URL" />
        </div>
      </div>
    );
  }
}

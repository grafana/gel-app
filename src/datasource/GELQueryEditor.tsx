// Libraries
import React, { PureComponent } from 'react';

// Types
import { GELDataSource } from './GELDataSource';
import { GELQuery, GELDataSourceOptions } from './types';

import { QueryEditorProps } from '@grafana/ui';
type Props = QueryEditorProps<GELDataSource, GELQuery, GELDataSourceOptions>;

interface State {}

export class GELQueryEditor extends PureComponent<Props, State> {
  state: State = {
    // nothing
  };

  render() {
    const { query } = this.props;

    return (
      <div>
        TODO... query....
        <br>
          <pre>{JSON.stringify(query)}</pre>
        </br>
      </div>
    );
  }
}

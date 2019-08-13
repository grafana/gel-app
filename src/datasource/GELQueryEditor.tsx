// Libraries
import React, { PureComponent } from 'react';

// Types
import { GELDataSource } from './GELDataSource';
import { GELQuery, GELDataSourceOptions } from './types';

import { getBackendSrv } from '@grafana/runtime';

import { QueryEditorProps, Select, FormLabel, DataQuery } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';
type Props = QueryEditorProps<GELDataSource, GELQuery, GELDataSourceOptions>;

interface State {
  datasources: Array<SelectableValue<number>>;
}

export class GELQueryEditor extends PureComponent<Props, State> {
  state: State = {
    datasources: [],
  };

  async componentDidMount() {
    const dslist: any[] = await getBackendSrv().get('/api/datasources');
    this.setState({
      datasources: dslist.map(ds => {
        return {
          label: ds.name,
          value: ds.id, // number
          imgUrl: ds.typeLogoUrl,
        };
      }),
    });
  }

  onSelectDataDource = (item: SelectableValue<number>) => {
    const { query, onChange } = this.props;
    if(!query.queries) {
      query.queries = [];
    }

    query.queries.push({
      refId: 'A',
      datasource:item.label,
    });

    onChange(query);
    console.log('SELECT', item);
  };

  renderQuery(query:DataQuery) {
    return <div>
      <pre>{JSON.stringify(query)}</pre>
    </div>
  }

  render() {
    const { query } = this.props;
    const { datasources } = this.state;
    const selected = {
      label: '   ',
      value: undefined,
    };

    return (
      <div>
        TODO... query....
        <pre>{JSON.stringify(query)}</pre>
        {query.queries.map( q => {
          return this.renderQuery(q);
        })}

        <div className="form-field">
          <FormLabel width={6}>Add Query</FormLabel>
          <Select options={datasources} value={selected} onChange={this.onSelectDataDource} />
        </div>
      </div>
    );
  }
}

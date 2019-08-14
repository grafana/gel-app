// Libraries
import React, { PureComponent } from 'react';

// Types
import { GELDataSource } from './GELDataSource';
import { GELQuery, GELDataSourceOptions } from './types';

import { getBackendSrv } from '@grafana/runtime';

import { QueryEditorProps, Select, FormLabel, DataQuery } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';
import { QueryEditorRow } from './QueryEditorRow';
type Props = QueryEditorProps<GELDataSource, GELQuery, GELDataSourceOptions>;

interface State {
  datasources: Array<SelectableValue<number>>;
}

let x = 1;

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
    if (!query.queries) {
      query.queries = [];
    }

    query.queries.push({
      refId: (x++).toString(),
      datasource: item.label,
    });

    onChange(query);
    console.log('SELECT', item);
  };

  onRemoveQuery = (remove: DataQuery) => {
    const { query, onChange } = this.props;
    const queries = query.queries.filter(q => {
      return q.refId !== remove.refId;
    });
    onChange({ ...query, queries });
  };

  render() {
    const { query, panelData, onChange } = this.props;
    const { datasources } = this.state;
    const selected = {
      label: '   ',
      value: undefined,
    };
    if (!query.queries) {
      query.queries = [];
    }

    return (
      <div>
        {query.queries.map((q, index) => {
          return (
            <QueryEditorRow
              key={index}
              query={query}
              data={panelData}
              onRemoveQuery={this.onRemoveQuery}
              onChange={onChange as (query: DataQuery) => void}
            />
          );
        })}

        <div className="form-field">
          <FormLabel width={6}>Add Query</FormLabel>
          <Select options={datasources} value={selected} onChange={this.onSelectDataDource} />
        </div>
      </div>
    );
  }
}

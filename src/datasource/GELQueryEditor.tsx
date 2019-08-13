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

  renderQuery(query: DataQuery, index: number) {
    return (
      <div key={index}>
        <div className="query-editor-row__header">
          <div className="query-editor-row__ref-id">
            <span>{query.refId}</span>
          </div>
          <div className="query-editor-row__collapsed-text"></div>
          <div className="query-editor-row__actions">
            <button
              className="query-editor-row__action"
              title="Remove query"
              onClick={() => {
                this.onRemoveQuery(query);
              }}
            >
              <i className="fa fa-fw fa-trash"></i>
            </button>
          </div>
        </div>
        <div>
          <pre>{JSON.stringify(query)}</pre>
        </div>
      </div>
    );
  }

  render() {
    const { query } = this.props;
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
          return this.renderQuery(q, index);
        })}
        <div className="form-field">
          <FormLabel width={6}>Add Query</FormLabel>
          <Select options={datasources} value={selected} onChange={this.onSelectDataDource} />
        </div>
      </div>
    );
  }
}

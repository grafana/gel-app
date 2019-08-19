// Libraries
import React, { PureComponent, ChangeEvent } from 'react';

// Types
import { GELDataSource } from './GELDataSource';
import { GELQuery, GELDataSourceOptions } from './types';

import { getBackendSrv } from '@grafana/runtime';

import { QueryEditorProps, Select, FormLabel, DataQuery } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';
import { QueryEditorRow } from './QueryEditorRow';
import { getNextQueryID } from './util';

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
    if (!query.queries) {
      query.queries = [];
    }

    query.queries.push({
      refId: getNextQueryID(query),
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

  onExpressionChange = (evt: ChangeEvent<HTMLTextAreaElement>) => {
    const { query, onChange } = this.props;
    onChange({
      ...query,
      expression: evt.target.value,
    });
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
              query={q}
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

        <br />
        <br />

        <div>
          <div className="query-editor-row__header">
            <div className="query-editor-row__ref-id">
              <span>GEL Expression:</span>
            </div>
          </div>
          <div>
            <textarea value={query.expression} onChange={this.onExpressionChange} className="gf-form-input" rows={3} />
          </div>
        </div>
      </div>
    );
  }
}

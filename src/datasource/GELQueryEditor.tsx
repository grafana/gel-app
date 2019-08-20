// Libraries
import React, { PureComponent } from 'react';

// Types
import { GELDataSource } from './GELDataSource';
import { TempGELQueryWrapper, GELDataSourceOptions, GELQuery, GELQueryType, GEL_DS_KEY } from './types';

import { getBackendSrv } from '@grafana/runtime';

import { QueryEditorProps, Select, FormLabel, DataQuery } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';
import { QueryEditorRow } from './QueryEditorRow';
import { getNextQueryID } from './util';
import { GELQueryNode } from './GELQueryNode';

type Props = QueryEditorProps<GELDataSource, TempGELQueryWrapper, GELDataSourceOptions>;

interface State {
  datasources: Array<SelectableValue<number>>;
}

const gelTypes: Array<SelectableValue<GELQueryType>> = [
  { value: GELQueryType.math, label: 'Math Expression' },
  { value: GELQueryType.reduce, label: 'Reduce Results' },
];

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

  onSelectDataSource = (item: SelectableValue<number>) => {
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

  onSelectGELType = (item: SelectableValue<GELQueryType>) => {
    const { query, onChange } = this.props;
    if (!query.queries) {
      query.queries = [];
    }

    query.queries.push({
      refId: getNextQueryID(query),
      datasource: GEL_DS_KEY, // GEL Type!!!
      type: item.value,
    } as GELQuery);

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

  onChangeGELQuery = (gel: GELQuery) => {
    const { query, onChange } = this.props;
    const queries = query.queries.map(q => {
      if (q.refId === gel.refId) {
        return gel;
      }
      return q;
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
          if (q.datasource === GEL_DS_KEY) {
            const gelQuery = q as GELQuery;
            return (
              <div key={index}>
                <div className="query-editor-row__header">
                  <div className="query-editor-row__ref-id">
                    <span>{q.refId}</span>
                  </div>
                  <div className="query-editor-row__collapsed-text">
                    <span>GEL: {gelQuery.type}</span>
                  </div>
                  <div className="query-editor-row__actions">
                    <button className="query-editor-row__action" title="Remove query" onClick={() => this.onRemoveQuery(q)}>
                      <i className="fa fa-fw fa-trash"></i>
                    </button>
                  </div>
                </div>
                <div>
                  <div>
                    <GELQueryNode query={gelQuery} onChange={this.onChangeGELQuery} />
                  </div>
                </div>
                <br />
              </div>
            );
          }

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
          <Select options={datasources} value={selected} onChange={this.onSelectDataSource} />

          <FormLabel width={6}>Add GEL</FormLabel>
          <Select options={gelTypes} value={selected} onChange={this.onSelectGELType} />
        </div>
      </div>
    );
  }
}

// Libraries
import React, { PureComponent } from 'react';

// Types
import { GELDataSource } from './GELDataSource';
import { TempGELQueryWrapper, GELDataSourceOptions, GELQuery, GELQueryType, GEL_DS_KEY } from './types';

import { getBackendSrv } from '@grafana/runtime';

import { QueryEditorProps, Select, FormLabel, DataQuery, Button } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';
import { QueryEditorRow } from './QueryEditorRow';
import { getNextQueryID } from './util';
import { GELQueryNode } from './GELQueryNode';

type Props = QueryEditorProps<GELDataSource, TempGELQueryWrapper, GELDataSourceOptions>;

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
      datasources: dslist
        .filter(ds => {
          return ds.type !== 'gel-datasource';
        })
        .map(ds => {
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

  addGEL = () => {
    const { query, onChange } = this.props;
    if (!query.queries) {
      query.queries = [];
    }

    query.queries.push({
      refId: getNextQueryID(query),
      datasource: GEL_DS_KEY, // GEL Type!!!
      type: GELQueryType.math,
    } as GELQuery);

    onChange(query);
  };

  onRemoveQuery = (remove: DataQuery) => {
    const { query, onChange } = this.props;
    const queries = query.queries.filter(q => {
      return q.refId !== remove.refId;
    });
    onChange({ ...query, queries });
  };

  onChangeGELQuery = (update: DataQuery) => {
    const { query, onChange } = this.props;
    const queries = query.queries.map(q => {
      if (q.refId === update.refId) {
        return update;
      }
      return q;
    });
    onChange({ ...query, queries });
  };

  onToggleQueryHide = (update: GELQuery) => {
    this.onChangeGELQuery({
      ...update,
      hide: !update.hide,
    });
  };

  render() {
    const { query, data } = this.props;
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
            const isDisabled = q.hide;
            const gelQuery = q as GELQuery;
            return (
              <div key={index}>
                <div className="query-editor-row__header">
                  <div className="query-editor-row__ref-id">
                    <span>{q.refId}</span>
                  </div>
                  <div className="query-editor-row__collapsed-text">
                    <span>GEL:</span>
                  </div>
                  <div className="query-editor-row__actions">
                    <button className="query-editor-row__action" onClick={() => this.onToggleQueryHide(gelQuery)} title="Disable/enable query">
                      {isDisabled && <i className="fa fa-fw fa-eye-slash" />}
                      {!isDisabled && <i className="fa fa-fw fa-eye" />}
                    </button>
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

          return <QueryEditorRow key={index} query={q} data={data!} onRemoveQuery={this.onRemoveQuery} onChange={this.onChangeGELQuery} />;
        })}

        <div className="form-field">
          <FormLabel width={6}>Add Query</FormLabel>
          <Select options={datasources} value={selected} onChange={this.onSelectDataSource} />
          <Button variant={'inverse'} onClick={this.addGEL}>
            Add GEL
          </Button>
        </div>
      </div>
    );
  }
}

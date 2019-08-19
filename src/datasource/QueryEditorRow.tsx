// Libraries
import React, { PureComponent } from 'react';

import { AngularComponent, getAngularLoader, getDataSourceSrv } from '@grafana/runtime';

import { DataQuery, DataSourceApi, PanelData, DataQueryRequest } from '@grafana/ui';
import { LoadingState } from '@grafana/data';

interface Props {
  query: DataQuery;
  data: PanelData;
  onRemoveQuery: (query: DataQuery) => void;
  onChange: (query: DataQuery) => void;
}

interface State {
  datasource?: DataSourceApi;
  loadedDataSourceValue?: string;
  queryResponse?: PanelData;
}

interface PretendPanelModel {
  targets: DataQuery[];
}

export class QueryEditorRow extends PureComponent<Props, State> {
  element?: HTMLElement | null;
  angularScope?: AngularQueryComponentScope;
  angularQueryEditor?: AngularComponent;

  state: State = {
    datasource: undefined,
    queryResponse: undefined,
  };

  componentDidMount() {
    this.loadDatasource();
  }

  componentWillUnmount() {
    if (this.angularQueryEditor) {
      this.angularQueryEditor.destroy();
    }
  }

  getAngularQueryComponentScope(): AngularQueryComponentScope {
    const { query } = this.props;
    const { datasource } = this.state;

    return {
      datasource: datasource!,
      target: query,
      panel: {
        targets: [query],
      },
      refresh: () => {},
      render: () => {},
    };
  }

  async loadDatasource() {
    const { query } = this.props;
    const ds = query.datasource ? query.datasource : undefined;
    const datasource = await getDataSourceSrv().get(ds);

    this.setState({
      datasource,
      loadedDataSourceValue: datasource.name,
    });
  }

  componentDidUpdate(prevProps: Props) {
    const { loadedDataSourceValue, datasource } = this.state;
    const { data, query } = this.props;

    if (data !== prevProps.data) {
      this.setState({ queryResponse: filterPanelDataToQuery(data, query.refId) });

      if (this.angularScope) {
        // this.angularScope.range = getTimeSrv().timeRange();
      }

      if (this.angularQueryEditor) {
        // Some query controllers listen to data error events and need a digest
        // for some reason this needs to be done in next tick
        setTimeout(this.angularQueryEditor.digest);
      }
    }

    // check if we need to load another datasource
    if (!datasource || loadedDataSourceValue !== datasource.name) {
      if (this.angularQueryEditor) {
        this.angularQueryEditor.destroy();
        this.angularQueryEditor = undefined;
      }
      this.loadDatasource();
      return;
    }

    if (!this.element || this.angularQueryEditor) {
      return;
    }

    const loader = getAngularLoader();
    const template = '<plugin-component type="query-ctrl" />';
    const scopeProps = { ctrl: this.getAngularQueryComponentScope() };
    this.angularQueryEditor = loader.load(this.element, scopeProps, template);
    this.angularScope = scopeProps.ctrl;
  }

  renderPluginEditor() {
    const { query, data, onChange } = this.props;
    const { datasource, queryResponse } = this.state;

    if (!datasource || !datasource.components) {
      return <div>no datasource</div>;
    }

    if (datasource.components.QueryCtrl) {
      return <div ref={element => (this.element = element)} />;
    }

    if (datasource.components.QueryEditor) {
      const QueryEditor = datasource.components.QueryEditor;

      return (
        <QueryEditor
          query={query}
          datasource={datasource}
          onChange={onChange}
          onRunQuery={this.onRunQuery}
          queryResponse={queryResponse}
          panelData={data}
        />
      );
    }

    return <div>Data source plugin does not export any Query Editor component</div>;
  }

  onRunQuery = () => {
    console.log('Run Query');
  };

  onRemoveQuery = () => {
    this.props.onRemoveQuery(this.props.query);
  };

  render() {
    const { query } = this.props;
    const { datasource } = this.state;

    if (!datasource) {
      return null;
    }

    return (
      <div>
        <div className="query-editor-row__header">
          <div className="query-editor-row__ref-id">
            <span>{query.refId}</span>
          </div>
          <div className="query-editor-row__collapsed-text"></div>
          <div className="query-editor-row__actions">
            <button className="query-editor-row__action" title="Remove query" onClick={this.onRemoveQuery}>
              <i className="fa fa-fw fa-trash"></i>
            </button>
          </div>
        </div>
        <div>
          <div>{this.renderPluginEditor()}</div>
        </div>
        <br/>
      </div>
    );
  }
}

export interface AngularQueryComponentScope {
  target: DataQuery;
  panel: PretendPanelModel;
  // dashboard: DashboardModel;
  // events: Emitter;
  refresh: () => void;
  render: () => void;
  datasource: DataSourceApi;
  toggleEditorMode?: () => void;
  getCollapsedText?: () => string;
  // range: TimeRange;
}

/**
 * Get a version of the PanelData limited to the query we are looking at
 */
export function filterPanelDataToQuery(data: PanelData, refId: string): PanelData | undefined {
  const series = data.series.filter(series => series.refId === refId);

  // No matching series
  if (!series.length) {
    return undefined;
  }

  // Don't pass the request if all requests are the same
  const request: DataQueryRequest | undefined = undefined;
  // TODO: look in sub-requets to match the info

  // Only say this is an error if the error links to the query
  let state = LoadingState.Done;
  const error = data.error && data.error.refId === refId ? data.error : undefined;
  if (error) {
    state = LoadingState.Error;
  }

  return {
    state,
    series,
    request,
    error,
  };
}

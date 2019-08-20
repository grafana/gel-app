// Libraries
import React, { PureComponent, ChangeEvent } from 'react';

// Types
import { GELQuery, GELQueryType } from './types';

interface Props {
  query: GELQuery;
  onChange: (value: GELQuery) => void;
}

interface State {}

export class GELQueryNode extends PureComponent<Props, State> {
  state: State = {};

  async componentDidMount() {}

  onExpressionChange = (evt: ChangeEvent<HTMLTextAreaElement>) => {
    const { query, onChange } = this.props;
    onChange({
      ...query,
      expression: evt.target.value,
    });
  };

  render() {
    const { query } = this.props;

    if (query.type === GELQueryType.math) {
      return (
        <div>
          <textarea value={query.expression} onChange={this.onExpressionChange} className="gf-form-input" rows={2} />
        </div>
      );
    }

    if (query.type === GELQueryType.reduce) {
      return <div>REDUCE: {JSON.stringify(query)}</div>;
    }

    return <div>Unknown Query Type: {JSON.stringify(query)}</div>;
  }
}

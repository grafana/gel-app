// Libraries
import React, { PureComponent, ChangeEvent } from 'react';

// Types
import { GELQuery, GELQueryType } from './types';
import { SelectableValue, ReducerID } from '@grafana/data';
import { FormLabel, Select, FormField } from '@grafana/ui';

interface Props {
  query: GELQuery;
  onChange: (value: GELQuery) => void;
}

interface State {}

const gelTypes: Array<SelectableValue<GELQueryType>> = [{ value: GELQueryType.math, label: 'Math' }, { value: GELQueryType.reduce, label: 'Reduce' }];

const reducerTypes: Array<SelectableValue<string>> = [
  { value: ReducerID.min, label: 'Min', description: 'Get the minimum value' },
  { value: ReducerID.max, label: 'Max', description: 'Get the maximum value' },
  { value: ReducerID.mean, label: 'Mean', description: 'Get the average value' },
];

export class GELQueryNode extends PureComponent<Props, State> {
  state: State = {};

  onSelectGELType = (item: SelectableValue<GELQueryType>) => {
    const { query, onChange } = this.props;
    const q = {
      ...query,
      type: item.value!,
    };

    if (q.type === GELQueryType.reduce) {
      if (!q.reducer) {
        q.reducer = ReducerID.mean;
      }
      q.expression = undefined;
    } else {
      q.reducer = undefined;
    }

    onChange(q);
  };

  onSelectReducer = (item: SelectableValue<string>) => {
    const { query, onChange } = this.props;
    onChange({
      ...query,
      reducer: item.value!,
    });
  };

  onExpressionChange = (evt: ChangeEvent<any>) => {
    const { query, onChange } = this.props;
    onChange({
      ...query,
      expression: evt.target.value,
    });
  };

  render() {
    const { query } = this.props;
    const selected = gelTypes.find(o => o.value === query.type);
    const reducer = reducerTypes.find(o => o.value === query.reducer);

    return (
      <div>
        <div className="form-field">
          <Select options={gelTypes} value={selected} onChange={this.onSelectGELType} />

          {query.type === GELQueryType.reduce && (
            <>
              <FormLabel width={5}>Function:</FormLabel>
              <Select options={reducerTypes} value={reducer} onChange={this.onSelectReducer} />

              <FormField label="Fields:" labelWidth={5} onChange={this.onExpressionChange} value={query.expression} />
            </>
          )}
        </div>
        {query.type === GELQueryType.math && (
          <textarea value={query.expression} onChange={this.onExpressionChange} className="gf-form-input" rows={2} />
        )}
      </div>
    );

    //   if (query.type === GELQueryType.reduce) {
    //     return <div>REDUCE: {JSON.stringify(query)}</div>;
    //   }

    //   return <div>
    //   <FormLabel width={6}>Add GEL</FormLabel>
    //   <Select options={gelTypes} value={selected} onChange={this.onSelectGELType} />

    //   <textarea value={query.expression} onChange={this.onExpressionChange} className="gf-form-input" rows={2} />
    // </div>
    //   return <div>Unknown Query Type: {JSON.stringify(query)}</div>;
  }
}

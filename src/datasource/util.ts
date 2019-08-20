import { GELQuery } from './types';
import { DataFrame, DataFrameHelper } from '@grafana/data';

export function getNextQueryID(query: GELQuery) {
  if (!query || !query.queries) {
    return 'GA';
  }
  const A = 'A'.charCodeAt(0);
  const ids = query.queries.map(q => q.refId);
  for (let i = query.queries.length; i < 27; i++) {
    const id = 'G' + String.fromCharCode(A + i);
    if (!ids.includes(id)) {
      return id;
    }
  }
  return 'G' + Date.now(); //
}

export function responseToDataFrame(rsp: any): DataFrame {
  const frame = new DataFrameHelper();
  frame.name = rsp.Name;
  for (const f of rsp.Fields) {
    frame.addField({
      name: f.Name,
      values: f.Vector,
    });
  }
  return frame;
}

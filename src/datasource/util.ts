import { TempGELQueryWrapper } from './types';
import { DataFrame, DataFrameHelper, dateTime, FieldType } from '@grafana/data';

export function getNextQueryID(query: TempGELQueryWrapper) {
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

export function gelResponseToDataFrames(rsp: any[]): DataFrame[] {
  return rsp.map((v: any) => {
    const frame = new DataFrameHelper();
    frame.name = v.name;
    frame.refId = v.refId;
    if (v.labels) {
      frame.labels = v.labels;
    }
    for (const f of v.fields) {
      let v = f.values;
      const type: FieldType = f.type;
      // HACK: this should be supported out-of-the-box
      if (type === FieldType.time) {
        v = v.map((str: string) => {
          return dateTime(str).valueOf();
        });
      }
      frame.addField({
        ...f,
        values: v,
      });
    }
    return frame;
  });
}

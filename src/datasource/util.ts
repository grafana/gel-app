import { TempGELQueryWrapper } from './types';
import { DataFrame, MutableDataFrame, FieldType } from '@grafana/data';
import { Table } from 'apache-arrow';
import { decode } from 'base64-arraybuffer';

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

export function gelResponseToDataFrames(rsp: any): DataFrame[] {
  if (rsp.results) {
    const frames: DataFrame[] = [];
    for (const v of Object.values(rsp.results)) {
      const raw = decode((v as any).meta.GC);
      const uints = new Uint8Array(raw);
      const arrowTable = Table.from([uints]);
      console.log('ENC', arrowTable.toString());
      console.log('FIELDS', arrowTable.schema.fields.length);
    }
    return frames;
  }
  return rsp.Frames.map((v: any) => {
    const frame = new MutableDataFrame();
    frame.name = v.name;
    frame.refId = v.refId;
    if (v.labels) {
      frame.labels = v.labels;
    }
    for (const f of v.fields) {
      let v = f.values;
      const type: FieldType = f.type;
      // HACK: this should be supported out-of-the-box
      // String as ms date
      if (type === FieldType.time) {
        v = v.map((str: string) => {
          return parseInt(str.slice(0, -6), 10);
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

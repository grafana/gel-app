import { TempGELQueryWrapper } from './types';
import { DataFrame, MutableDataFrame, FieldType, Field, Vector, ArrayVector } from '@grafana/data';
import { Table, ArrowType } from 'apache-arrow';

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

export function base64StringToArrowTable(text: string) {
  const b64 = atob(text);
  const arr = Uint8Array.from(b64, c => {
    return c.charCodeAt(0);
  });
  return Table.from(arr);
}

export function arrowTableToDataFrame(table: Table): DataFrame {
  const fields: Field[] = [];
  for (let i = 0; i < table.numCols; i++) {
    const col = table.getColumnAt(i);
    if (col) {
      const schema = table.schema.fields[i];
      let type = FieldType.other;
      let values: Vector<any> = col; //
      if ('Time' === col.name) {
        type = FieldType.time;

        // Silly conversion until we agree on date formats
        const ms: number[] = new Array(col.length);
        for (let j = 0; j < col.length; j++) {
          ms[j] = col.get(j) / 1000000; // nanoseconds to milliseconds
        }
        values = new ArrayVector(ms);
      } else {
        switch ((schema.typeId as unknown) as ArrowType) {
          case ArrowType.FloatingPoint: {
            break;
          }
          default:
            console.log('UNKOWN Type:', schema);
        }
      }

      fields.push({
        name: col.name,
        type,
        config: {}, // TODO, pull from metadata
        values,
      });
    }
  }
  return {
    fields,
    length: table.length,
  };
}

export function gelResponseToDataFrames(rsp: any): DataFrame[] {
  if (rsp.results) {
    const frames: DataFrame[] = [];
    for (const res of Object.values(rsp.results)) {
      for (const b of (res as any).meta) {
        const t = base64StringToArrowTable(b as string);
        frames.push(arrowTableToDataFrame(t));
      }
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

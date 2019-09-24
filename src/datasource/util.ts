import { TempGELQueryWrapper } from './types';
import { DataFrame, MutableDataFrame, FieldType, Field, Vector } from '@grafana/data';
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
      const values: Vector<any> = col;
      switch ((schema.typeId as unknown) as ArrowType) {
        case ArrowType.Decimal:
        case ArrowType.Int:
        case ArrowType.FloatingPoint: {
          type = FieldType.number;
          break;
        }
        case ArrowType.Bool: {
          type = FieldType.boolean;
          break;
        }
        case ArrowType.Timestamp: {
          type = FieldType.time;
          break;
        }
        default:
          console.log('UNKNOWN Type:', schema);
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
      frame.addField({
        ...f,
        values: f.values,
      });
    }
    return frame;
  });
}

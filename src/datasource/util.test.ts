import { gelResponseToDataFrames } from './util';
import { toDataFrameDTO } from '@grafana/data';

/* tslint:disable */
const resp = {
  results: {
    '': {
      refId: '',
      meta: [
        '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',
      ],
      series: [],
      tables: null,
      frames: null,
    },
  },
};
/* tslint:enable */

describe('GEL Utils', () => {
  // test('should parse sample GEL output', () => {
  //   const frames = gelResponseToDataFrames(resp);
  //   const frame = frames[0];
  //   expect(frame.name).toEqual('BBB');
  //   expect(frame.fields.length).toEqual(2);
  //   expect(frame.length).toEqual(resp.Frames[0].fields[0].values.length);

  //   const timeField = frame.fields[0];
  //   expect(timeField.name).toEqual('Time');

  //   // The whole response
  //   expect(frames).toMatchSnapshot();
  // });

  test('should parse output with dataframe', () => {
    const frames = gelResponseToDataFrames(resp);
    for (const frame of frames) {
      console.log('Frame', frame.refId + ' // ' + frame.labels);
      for (const field of frame.fields) {
        console.log(' > ', field.name, field.values.toArray());
      }
    }

    const norm = frames.map( f => toDataFrameDTO(f) );
    expect(norm).toMatchSnapshot();
  });
});

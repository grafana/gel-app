import { gelResponseToDataFrames } from './util';

const resp = {
  Frames: [
    {
      name: 'BBB',
      fields: [
        {
          name: 'Time',
          type: 'time',
          values: [
            '1568390765000000000',
            '1568390770000000000',
            '1568390775000000000',
            '1568390780000000000',
            '1568390785000000000',
            '1568390790000000000',
            '1568390795000000000',
          ],
        },
        { name: 'GB-series', type: 'number', values: [2, 2, 0, 0, 0, 2, 2] },
      ],
      refId: 'GB',
    },
    {
      name: 'CCC',
      fields: [
        {
          name: 'Time',
          type: 'time',
          values: [
            '1568390765000000000',
            '1568390770000000000',
            '1568390775000000000',
            '1568390780000000000',
            '1568390785000000000',
            '1568390790000000000',
            '1568390795000000000',
          ],
        },
        { type: 'number', values: [5, 5, 3, 3, 3, 5, 5] },
      ],
      refId: 'GC',
    },
  ],
};

describe('GEL Utils', () => {
  test('should parse sample GEL output', () => {
    const frames = gelResponseToDataFrames(resp);
    const frame = frames[0];
    expect(frame.name).toEqual('BBB');
    expect(frame.fields.length).toEqual(2);
    expect(frame.length).toEqual(resp.Frames[0].fields[0].values.length);

    const timeField = frame.fields[0];
    expect(timeField.name).toEqual('Time');

    // The whole response
    expect(frames).toMatchSnapshot();
  });
});

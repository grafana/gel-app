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

/* tslint:disable */
const respDF = {
  results: {
    '': {
      refId: '',
      meta: {
        GX:
          'QVJST1cxAADsAQAAEAAAAAAACgAOAAwACwAEAAoAAAAUAAAAAAAAAQMACgAMAAAACAAEAAoAAAAIAAAAfAAAAAMAAABMAAAAKAAAAAQAAACg/v//CAAAAAwAAAACAAAAR0IAAAUAAAByZWZJZAAAAMD+//8IAAAADAAAAAAAAAAAAAAABAAAAG5hbWUAAAAA4P7//wgAAAAUAAAACQAAAHRlc3Q9dGVzdAAAAAYAAABsYWJlbHMAAAIAAACsAAAABAAAAG7///8UAAAAaAAAAHAAAAAAAAMBcAAAAAIAAAAsAAAABAAAADj///8IAAAAEAAAAAYAAABudW1iZXIAAAQAAAB0eXBlAAAAAFz///8IAAAAFAAAAAkAAABHQi1zZXJpZXMAAAAEAAAAbmFtZQAAAAAAAAAAAAAGAAgABgAGAAAAAAACAAkAAABHQi1zZXJpZXMAEgAYABQAEwASAAwAAAAIAAQAEgAAABQAAABsAAAAcAAAAAAABQFsAAAAAgAAADQAAAAEAAAA3P///wgAAAAQAAAABAAAAHRpbWUAAAAABAAAAHR5cGUAAAAACAAMAAgABAAIAAAACAAAABAAAAAEAAAAVGltZQAAAAAEAAAAbmFtZQAAAAAAAAAABAAEAAQAAAAEAAAAVGltZQAAAAAAAAAAzAAAABQAAAAAAAAADAAWABQAEwAMAAQADAAAAMgAAAAAAAAAFAAAAAAAAAMDAAoAGAAMAAgABAAKAAAAFAAAAGgAAAAGAAAAAAAAAAAAAAAFAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIAAAAAAAAAAgAAAAAAAAAHgAAAAAAAAAmAAAAAAAAAAAAAAAAAAAAJgAAAAAAAAAMAAAAAAAAAAAAAAAAgAAAAYAAAAAAAAAAAAAAAAAAAAGAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAATAAAAJgAAADkAAABMAAAAXwAAAHIAAAAAAAAAMTU2ODk4NjgwMDAwMDAwMDAwMDE1Njg5ODY4NjAwMDAwMDAwMDAxNTY4OTg2OTIwMDAwMDAwMDAwMTU2ODk4Njk4MDAwMDAwMDAwMDE1Njg5ODcwNDAwMDAwMDAwMDAxNTY4OTg3MTAwMDAwMDAwMDAwAAAAAAAAAAAAAAAA8D8AAAAAAADwPwAAAAAAAABAAAAAAAAAAEAAAAAAAAAAQAAAAAAAAPA/EAAAAAwAFAASAAwACAAEAAwAAAAQAAAALAAAADwAAAAAAAMAAQAAAPgBAAAAAAAA0AAAAAAAAADIAAAAAAAAAAAAAAAAAAAAAAAAAAAACgAMAAAACAAEAAoAAAAIAAAAfAAAAAMAAABMAAAAKAAAAAQAAACg/v//CAAAAAwAAAACAAAAR0IAAAUAAAByZWZJZAAAAMD+//8IAAAADAAAAAAAAAAAAAAABAAAAG5hbWUAAAAA4P7//wgAAAAUAAAACQAAAHRlc3Q9dGVzdAAAAAYAAABsYWJlbHMAAAIAAACsAAAABAAAAG7///8UAAAAaAAAAHAAAAAAAAMBcAAAAAIAAAAsAAAABAAAADj///8IAAAAEAAAAAYAAABudW1iZXIAAAQAAAB0eXBlAAAAAFz///8IAAAAFAAAAAkAAABHQi1zZXJpZXMAAAAEAAAAbmFtZQAAAAAAAAAAAAAGAAgABgAGAAAAAAACAAkAAABHQi1zZXJpZXMAEgAYABQAEwASAAwAAAAIAAQAEgAAABQAAABsAAAAcAAAAAAABQFsAAAAAgAAADQAAAAEAAAA3P///wgAAAAQAAAABAAAAHRpbWUAAAAABAAAAHR5cGUAAAAACAAMAAgABAAIAAAACAAAABAAAAAEAAAAVGltZQAAAAAEAAAAbmFtZQAAAAAAAAAABAAEAAQAAAAEAAAAVGltZQAAAAAYAgAAQVJST1cx',
      },
      series: [],
      tables: null,
      frames: null,
    },
  },
};
/* tslint:enable */

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

  test('should parse output with dataframe', () => {
    const frames = gelResponseToDataFrames(respDF);
    for (const frame of frames) {
      console.log('Frame', frame.refId + ' // ' + frame.labels);
      for (const field of frame.fields) {
        console.log(' > ', field.name, field.values.toArray());
      }
    }
    expect(frames).toBeDefined();
  });
});

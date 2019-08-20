import { responseToDataFrame } from './util';

const resp = {
  Values: [
    {
      Name: 'AAA',
      Fields: [
        {
          Name: 'Time',
          Type: 'time',
          Unit: '',
          Vector: [
            '2019-08-20T12:51:00-04:00',
            '2019-08-20T12:52:00-04:00',
            '2019-08-20T12:53:00-04:00',
            '2019-08-20T12:54:00-04:00',
            '2019-08-20T12:55:00-04:00',
            '2019-08-20T12:56:00-04:00',
          ],
        },
        {
          Name: '',
          Type: 'number',
          Unit: '',
          Vector: [0.06666666666666667, 0.06666666666666667, 0.06666666666666667, 0.16666666666666666, 0.16666666666666666, 0.16666666666666666],
        },
      ],
      Labels: null,
    },
  ],
};

describe('GEL Utils', () => {
  test('should parse sample GEL output', () => {
    const frame = responseToDataFrame(resp)[0];
    expect(frame.name).toEqual('AAA');
    expect(frame.fields.length).toEqual(2);
    expect(frame.length).toEqual(resp.Values[0].Fields[0].Vector.length);

    const timeField = frame.fields[0];

    expect(timeField.name).toEqual('Time');
    expect(timeField).toMatchSnapshot();
  });
});

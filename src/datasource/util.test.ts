import { gelResponseToDataFrames } from './util';

const resp = [
  {
    name: 'AAA',
    fields: [
      {
        name: 'Time',
        type: 'time',
        values: ['2019-08-27T16:18:00-04:00'],
      },
      {
        name: 'GE-series',
        type: 'number',
        values: [2],
      },
    ],
    labels: null,
    refId: 'GE',
  },
  {
    name: 'BBB',
    fields: [
      {
        name: 'Time',
        type: 'time',
        values: [
          '2019-08-27T16:18:34.433-04:00',
          '2019-08-27T16:18:35.433-04:00',
          '2019-08-27T16:18:36.433-04:00',
          '2019-08-27T16:18:37.433-04:00',
          '2019-08-27T16:18:38.433-04:00',
        ],
      },
      {
        name: 'GA-series',
        type: 'number',
        values: [60.466028797961954, 60.906537886006966, 61.07109793922545, 61.00881212641243, 60.933449623483696],
      },
    ],
    labels: null,
    refId: 'GA',
  },
  {
    name: 'CCC',
    fields: [
      {
        name: '',
        type: 'number',
        values: [60.8771852746181],
      },
    ],
    labels: null,
    refId: 'GC',
  },
  {
    name: 'DDD',
    fields: [
      {
        name: 'Time',
        type: 'time',
        values: [
          '2019-08-27T16:18:34.433-04:00',
          '2019-08-27T16:18:35.433-04:00',
          '2019-08-27T16:18:36.433-04:00',
          '2019-08-27T16:18:37.433-04:00',
          '2019-08-27T16:18:38.433-04:00',
        ],
      },
      {
        name: '',
        type: 'number',
        values: [0.9932461319490806, 1.0004821611126804, 1.0031853093031915, 1.0021621704617347, 1.0009242271733128],
      },
    ],
    labels: null,
    refId: 'GD',
  },
];

describe('GEL Utils', () => {
  test('should parse sample GEL output', () => {
    const frames = gelResponseToDataFrames(resp);
    const frame = frames[0];
    expect(frame.name).toEqual('AAA');
    expect(frame.fields.length).toEqual(2);
    expect(frame.length).toEqual(resp[0].fields[0].values.length);

    const timeField = frame.fields[0];
    expect(timeField.name).toEqual('Time');

    // The whole response
    expect(frames).toMatchSnapshot();
  });
});

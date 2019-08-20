import { responseToDataFrame } from './util';

const resp = {
  Values: [
    {
      Name: 'AAA',
      Fields: [
        {
          Name: 'Time',
          Type: 4,
          Unit: '',
          Vector: [
            '2019-08-20T10:29:00-04:00',
            '2019-08-20T10:30:00-04:00',
            '2019-08-20T10:31:00-04:00',
            '2019-08-20T10:32:00-04:00',
            '2019-08-20T10:33:00-04:00',
            '2019-08-20T10:34:00-04:00',
          ],
        },
        { Name: '', Type: 1, Unit: '', Vector: [2, 4, 4, 4, 2, 2] },
      ],
      Labels: null,
    },
  ],
};

describe('PluginDatasource', () => {
  describe('when querying', () => {
    test('should return the saved data with a query', () => {
      const frame = responseToDataFrame(resp)[0];
      expect(frame.name).toEqual('AAA');
      expect(frame.fields.length).toEqual(2);
      expect(frame.fields[0].name).toEqual('Time');
      expect(frame.length).toEqual(resp.Values[0].Fields[0].Vector.length);
    });
  });
});

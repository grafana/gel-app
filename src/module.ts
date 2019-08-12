import { ExampleConfigCtrl } from './legacy/config';
import { AppPlugin } from '@grafana/ui';
import { ExampleRootPage } from './ExampleRootPage';
import { GELAppSettings } from './types';

// Legacy exports just for testing
export { ExampleConfigCtrl as ConfigCtrl };

export const plugin = new AppPlugin<GELAppSettings>().setRootPage(ExampleRootPage);

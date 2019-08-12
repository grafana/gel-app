import { PluginMeta } from '@grafana/ui';

export class ExampleConfigCtrl {
  static templateUrl = 'legacy/config.html';

  appEditCtrl: any;
  appModel: PluginMeta;

  /** @ngInject */
  constructor($scope: any, $injector: any) {
    this.appEditCtrl.setPostUpdateHook(this.postUpdate.bind(this));

    const app = (this as any).appModel || {};
    if (!app.jsonData) {
      app.jsonData = {};
    }
    this.appModel = app as PluginMeta;

    console.log('ExampleConfigCtrl', this);
  }

  postUpdate() {
    if (!this.appModel.enabled) {
      console.log('Not enabled...');
      return;
    }

    // TODO, can do stuff after update
    console.log('Post Update:', this);
  }
}

library view_test_types;

import 'dart:convert';
import 'dart:html';

import 'package:angular/angular.dart';

@Component(
    selector: 'view-test-types',
    templateUrl: 'packages/stein/component/view_test_types.html',
    cssUrl: 'packages/stein/component/view_test_types.css',
    publishAs: 'cmp')
class ViewTestTypesCtrl {
  @NgOneWay('testTypes')
  List<String> testTypes;

  String _project = 'none';

  String get project => this._project;

  ViewTestTypesCtrl(RouteProvider routeProvider) {
    this._project = routeProvider.parameters['project'];
    HttpRequest.getString('/projects/$project/types').then((String data) => this.testTypes = JSON.decode(data));
    print("ViewTestTypesCtrl: $_project");
  }
}

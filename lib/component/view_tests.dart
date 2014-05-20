library view_tests;

import 'dart:convert';
import 'dart:html';

import 'package:angular/angular.dart';

@Component(
    selector: 'view-tests',
    templateUrl: 'packages/stein/component/view_tests.html',
    cssUrl: 'packages/stein/component/view_tests.css',
    publishAs: 'cmp')
class ViewTestsCtrl {
  @NgOneWay('tests')
  List<String> tests;
  @NgOneWay('testStatus')
  Map<String, bool> testStatus = {};

  String _project;
  String _testType;

  String get project => _project;
  String get testType => _testType;

  String status(test) => testStatus[test] ? "pass" : "fail";

  ViewTestsCtrl(RouteProvider routeProvider) {
    this._project = routeProvider.parameters['project'];
    this._testType = routeProvider.parameters['testType'];
    HttpRequest.getString('/projects/$project/types/$testType')
      .then((String data) => this.tests = JSON.decode(data))
      .then((List<String> tests) {
          tests.forEach((t) => HttpRequest.getString('/projects/$project/tests/$t').then((data) {
              var status = JSON.decode(data);
              this.testStatus[t] = (status["Fail"] == 0 && status["Error"] == 0);
            }));
        });
  }
}

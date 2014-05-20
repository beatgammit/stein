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
  List<String> tests = ['first test', 'second test'];

  String _project;
  String _testType;

  String get project => _project;
  String get testType => _testType;

  ViewTestsCtrl(RouteProvider routeProvider) {
    this._project = routeProvider.parameters['project'];
    this._testType = routeProvider.parameters['testType'];
    HttpRequest.getString('/projects/$project/types/$testType').then((String data) => this.tests = JSON.decode(data));
  }
}

library view_test;

import 'dart:convert';
import 'dart:html';

import 'package:angular/angular.dart';

@Component(
    selector: 'view-test',
    templateUrl: 'packages/stein/component/view_test.html',
    cssUrl: 'packages/stein/component/view_test.css',
    publishAs: 'cmp')
class ViewTestCtrl {

  @NgOneWay('suite')
  Map suite;

  Map<String, int> get counts => suite == null ? null : suite["Final"]["Counts"] as Map<String, int>;
  int get total => counts == null ? null : counts["Total"];
  int get pass => counts == null ? null : counts["Pass"];
  int get error => counts == null ? null : counts["Error"];
  int get fail => counts == null ? null : counts["Fail"];
  int get skip => counts == null ? null : counts["Omit"];
  int get todo => counts == null ? null : counts["Todo"];

  Map get cases => suite["Cases"];

  ViewTestCtrl(RouteProvider routeProvider) {
    String project = routeProvider.parameters['project'];
    String id = routeProvider.parameters['id'];

    HttpRequest.getString('/projects/$project/tests/$id').then((String data) => this.suite = JSON.decode(data));
  }
}

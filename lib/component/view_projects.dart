library view_projects_component;

import 'package:angular/angular.dart';
import 'dart:html';
import 'dart:convert';

@Component(
    selector: 'view-projects',
    templateUrl: 'packages/stein/component/view_projects.html',
    cssUrl: 'packages/stein/component/view_projects.css',
    publishAs: 'cmp')
class ViewProjectsCtrl {
  @NgOneWay('projects')
  List<String> projects = ['hello', 'there'];

  ViewProjectsCtrl(RouteProvider routeProvider) {
    HttpRequest.getString('/projects').then((String data) => this.projects = JSON.decode(data));
  }
}

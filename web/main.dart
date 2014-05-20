library stein;

import 'dart:async';
import 'dart:html';
import 'dart:convert';
import 'package:angular/angular.dart';
import 'package:angular/application_factory.dart';
import 'package:logging/logging.dart';

import 'package:stein/component/view_projects.dart';
import 'package:stein/component/view_test_types.dart';
import 'package:stein/component/view_tests.dart';
import 'package:stein/component/view_test.dart';

// Temporary, please follow https://github.com/angular/angular.dart/issues/476
@MirrorsUsed(targets: const['stein'], override: '*')
import 'dart:mirrors';

class SteinModule extends Module {
  SteinModule() {
    //type(SteinController);
    type(ViewProjectsCtrl);
    type(ViewTestTypesCtrl);
    type(ViewTestsCtrl);
    type(ViewTestCtrl);
    value(RouteInitializerFn, steinRouteInitializer);
    value(NgRoutingUsePushState, new NgRoutingUsePushState.value(false));
  }
}

void steinRouteInitializer(Router router, RouteViewFactory view) {
  // specific -> general
  router.root
    ..addRoute(
        name: 'test',
        path: '/projects/:project/tests/:id',
        enter: view('./view/view_test.html'))
    ..addRoute(
        name: 'tests',
        path: '/projects/:project/types/:testType',
        enter: view('./view/view_tests.html'))
    ..addRoute(
        name: 'testTypes',
        path: '/projects/:project/types',
        enter: view('./view/view_test_types.html'))
    ..addRoute(
        defaultRoute: true,
        name: 'projects',
        path: '/projects',
        enter: view('./view/view_projects.html'))
    ;
}

void main() {
  Logger.root..level = Level.FINEST
             ..onRecord.listen((LogRecord r) { print(r.message); });
  applicationFactory()
      .addModule(new SteinModule())
      .run();
}

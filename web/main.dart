library stein;

import 'dart:async';
import 'dart:html';
import 'dart:convert';
import 'package:angular/angular.dart';

// Temporary, please follow https://github.com/angular/angular.dart/issues/476
@MirrorsUsed(targets: const['stein'], override: '*')
import 'dart:mirrors';

@NgController(selector: '[stein-app]', publishAs: 'ctrl')
class SteinController {
  List<String> tests;
  List<String> projects;
  String selectedProject;
  String selectedTest;
  TestSuite curTest;

  SteinController() {
    _loadProjects().then((proj) => this.projects = proj);
  }

  Future<List<String>> _loadProjects() {
    return HttpRequest.getString('/projects').then((data) {
      var projects = JSON.decode(data);

      if (projects == null || projects.isEmpty) {
        return [];
      }

      projects.sort();
      return projects;
    });
  }

  Future<List<String>> _loadTests(String project) {
    return HttpRequest.getString('/projects/$project/tests').then((data) {
      var tests = JSON.decode(data);

      if (tests == null || tests.isEmpty) {
        return [];
      }

      tests.sort();
      return tests;
    });
  }

  Future<List<String>> _loadTest(String project, String test) {
    return HttpRequest.getString('/projects/$project/tests/$test').then((data) {
      var test = JSON.decode(data);

      if (test == null) {
        return [];
      }

      return test;
    });
  }

  void selectProject() {
    // TODO: remove this once bug 399 is fixed:
    // https://github.com/angular/angular.dart/issues/399
    new Future(() {
      _loadTests(selectedProject).then((tests) {
        this.tests = tests;
      });
    });
  }

  void selectTest() {
    // TODO: remove this once bug 399 is fixed:
    // https://github.com/angular/angular.dart/issues/399
    new Future(() {
      _loadTest(selectedProject, selectedTest).then((test) {
        this.curTest = test;
      });
    });
  }
}

class SteinModule extends Module {
  SteinModule() {
    type(SteinController);
  }
}

void main() {
  ngBootstrap(module: new SteinModule());
}

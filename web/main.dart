import 'dart:html';
import 'dart:convert';

void loadTest(project, test) {
  HttpRequest.getString('/projects/$project/tests/$test').then((data) {
    var suite = JSON.decode(data);
    if (suite == null) {
      return;
    }

    querySelectorAll('.suite:not(.template)').forEach((el) => el.remove());

    var suiteEl = querySelector('.suite.template').clone(true);
    suiteEl.classes.remove('template');

    print(suite);

    suiteEl.querySelector('.tally > .time').innerHtml = "${suite['Final']['Time']}";
    suiteEl.querySelector('.tally > .total').innerHtml = "${suite['Final']['Counts']['Total']}";
    suiteEl.querySelector('.tally > .pass').innerHtml = "${suite['Final']['Counts']['Pass']}";
    suiteEl.querySelector('.tally > .fail').innerHtml = "${suite['Final']['Counts']['Error']}";
    suiteEl.querySelector('.tally > .error').innerHtml = "${suite['Final']['Counts']['Fail']}";
    suiteEl.querySelector('.tally > .omit').innerHtml = "${suite['Final']['Counts']['Omit']}";
    suiteEl.querySelector('.tally > .todo').innerHtml = "${suite['Final']['Counts']['Todo']}";

    querySelector('#content').append(suiteEl);
  });
}

void loadTests(project) {
  var testList = querySelector('#tests');

  HttpRequest.getString('/projects/$project/tests').then((data) {
    var tests = JSON.decode(data);

    testList.children.clear();

    if (tests == null || tests.isEmpty) {
      return;
    }

    tests.sort();
    tests.forEach((test) {
      var el = document.createElement('option');
      el.value = test;
      el.innerHtml = test;
      testList.append(el);
    });

    loadTest(project, tests.first);
  });
}

void loadProjects() {
  HttpRequest.getString('/projects').then((data) {
    var projectList = querySelector('#projects');
    var projects = JSON.decode(data);

    if (projects == null || projects.isEmpty) {
      return;
    }

    projects.sort();
    projects.forEach((project) {
      var el = document.createElement('option');
      el.value = project;
      el.innerHtml = project;
      projectList.append(el);
    });

    loadTests(projects.first);
  });
}

void main() {
  loadProjects();

  querySelector('#projects').onChange.listen((ev) => loadTests(ev.target.value));
  querySelector('#tests').onChange.listen((ev) => loadTest(querySelector('#projects').value, ev.target.value));
}

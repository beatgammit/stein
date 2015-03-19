define([
	'angular',
	'angularRoute',

	'./views/projects/projects',
	'./views/project-tests/project-tests',
	'./views/project-test-types/project-test-types',
	'./views/project-type-tests/project-type-tests',
	'./views/project-test-run/project-test-run',
], function(angular, angularRoute, view1, view2) {
	'use strict';

	// Declare app level module which depends on views, and components
	return angular.module('stein', [
		'ngRoute',

		'stein.projects',
		'stein.project-tests',
		'stein.project-test-types',
		'stein.project-type-tests',
		'stein.project-test-run',
	]).
	config(['$routeProvider', function($routeProvider) {
		$routeProvider
			// http://tio-test.local/#/projects/NIO/tests/2014-05-30T13:13:42-06:00
			.when('/projects', {
				templateUrl: 'views/projects/projects.html',
				controller: 'ProjectsCtrl',
			})
			.when('/projects/:project/tests', {
				templateUrl: 'views/project-tests/project-tests.html',
				controller: 'ProjectTestsCtrl',
			})
			.when('/projects/:project/types', {
				templateUrl: 'views/project-test-types/project-test-types.html',
				controller: 'ProjectTestTypesCtrl',
			})
			.when('/projects/:project/types/:type', {
				templateUrl: 'views/project-type-tests/project-type-tests.html',
				controller: 'ProjectTypeTestsCtrl',
			})
			.when('/projects/:project/tests/:test', {
				templateUrl: 'views/project-test-run/project-test-run.html',
				controller: 'ProjectTestRunCtrl',
			})
			.otherwise({
				redirectTo: '/projects',
			});
	}]);
});


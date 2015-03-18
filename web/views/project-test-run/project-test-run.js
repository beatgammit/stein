define([
	'angular',
], function (angular) {
	'use strict';

	angular.module('stein.project-test-run', [])
		.controller('ProjectTestRunCtrl', ['$scope', '$routeParams', '$http', function ($scope, $routeParams, $http) {
			$scope.project = $routeParams.project;
			$scope.test = $routeParams.test;
			$http.get('/projects/' + $routeParams.project + '/tests/' + $routeParams.test).success(function (data) {
				$scope.testResult = data;
			});
		}]);
});
